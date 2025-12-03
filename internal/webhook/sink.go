package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/pkg/async"
	"github.com/google/uuid"

	"github.com/segmentio/kafka-go"
)

type Sink interface {
	Start(ctx context.Context)
	Enqueue(eventType EventType, payload map[string]any)
}

type InMemorySink struct {
	Logger *telemetry.Logger
	DB     *database.Database

	worker *async.Worker
}

// NewInMemorySink creates a Sink backed by async.Worker
func NewInMemorySink(tel *telemetry.Telemetry, db *database.Database) Sink {
	// Create worker metrics
	queueDepth, _ := tel.Metrics.Meter.Int64UpDownCounter("webhook_queue_depth")
	droppedEvents, _ := tel.Metrics.Meter.Int64Counter("webhook_dropped_total")
	flushFailures, _ := tel.Metrics.Meter.Int64Counter("webhook_flush_failures")
	batchSize, _ := tel.Metrics.Meter.Int64Histogram("webhook_flush_batch_size")

	metrics := &async.WorkerMetrics{
		QueueDepth:    queueDepth,
		DroppedEvents: droppedEvents,
		FlushFailures: flushFailures,
		BatchSize:     batchSize,
	}

	cfg := async.WorkerConfig{
		Logger:     tel.Logger,
		QueueSize:  200,
		BatchSize:  100,
		FlushEvery: 1 * time.Second,
		Metrics:    metrics,
		Tracer:     tel.Tracing.Tracer,
	}

	sink := &InMemorySink{
		Logger: tel.Logger,
		DB:     db,
	}

	// Worker flush function calls Sink.flush
	sink.worker = async.NewWorker(cfg, func(ctx context.Context, batch []async.Event) error {
		return sink.flush(ctx, batch)
	})

	return sink
}

// Start launches the worker
func (s *InMemorySink) Start(ctx context.Context) {
	s.worker.Start(ctx)
}

// Enqueue sends an event to the worker
func (s *InMemorySink) Enqueue(eventType EventType, payload map[string]any) {
	s.worker.Enqueue(async.Event(Event{
		Type: eventType,
		Data: payload,
	}))
}

// flush persists a batch of events (called by async.Worker)
func (s *InMemorySink) flush(ctx context.Context, batch []async.Event) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }() // safe if already committed

	var eventID string

	for _, event := range batch {
		ev := event.(Event)
		if err := tx.QueryRow(ctx, `
			INSERT INTO webhook.tbl_event (type, payload)
			VALUES ($1, $2)
			RETURNING id
		`, ev.Type, ev.Data).Scan(&eventID); err != nil {
			s.Logger.Error("failed inserting webhook event", "error", err)
			continue
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO webhook.tbl_delivery 
				(event_id, subscription_id, attempt_count, status, created_at)
			SELECT e.id, s.id, 0, 'pending', NOW()
			FROM webhook.tbl_event e
			CROSS JOIN webhook.tbl_subscription s
			WHERE e.id = $1
			  AND s.is_active = TRUE
			  AND $2 = ANY(s.event_types)
		`, eventID, ev.Type); err != nil {
			s.Logger.Error("failed inserting delivery rows", "event_id", eventID, "error", err)
		}
	}

	return tx.Commit(ctx)
}

type KafkaSink struct {
	Logger *telemetry.Logger
	DB     *database.Database

	Writer    *kafka.Writer
	DLQWriter *kafka.Writer
	Reader    *kafka.Reader

	// Producer buffering & workers
	enqueueC     chan kafka.Message
	producerWg   chan struct{} // signal worker finished (len = numWorkers)
	numProducers int

	// Consumer batching config
	batchSize  int
	flushEvery time.Duration
	buffer     []KafkaEvent
	lastFlush  time.Time

	// Limits / policies
	maxMessageBytes int // max allowed incoming payload size (bytes)

	// // Observability (placeholder)
	// metrics struct {
	// 	queueLen metric.Int64Observable // replace with your metrics
	// }
}

type KafkaEvent struct {
	Message kafka.Message
	Payload map[string]any
}

// NewKafkaSink returns a configured sink.
// brokers: list of Kafka brokers
// topic: main topic, dlqTopic: dead letter topic (e.g. "webhook-events-dlq")
func NewKafkaSink(logger *telemetry.Logger, db *database.Database, brokers []string) *KafkaSink {
	// configure writers/readers. Add TLS/SASL into Transport here if needed.
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          brokers,
		Topic:            "webhook-events",
		Balancer:         &kafka.LeastBytes{},
		CompressionCodec: kafka.Snappy.Codec(),
		RequiredAcks:     int(kafka.RequireAll),
		BatchSize:        100,
		BatchTimeout:     5 * time.Second,
		MaxAttempts:      5,
		Async:            false, // keep synchronous writes for retries
		// Transport: &kafka.Transport{TLS: &tls.Config{MinVersion: tls.VersionTLS12}}, // optional TLS
	})

	dlqWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        "webhook-events-dlq",
		Balancer:     &kafka.Hash{},
		BatchSize:    1,
		BatchTimeout: 500 * time.Millisecond,
		MaxAttempts:  3,
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        "gopenehr-webhook-group",
		Topic:          "webhook-events",
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0, // manual commits
		MaxBytes:       10e6,
		MaxWait:        500 * time.Millisecond,
	})

	return &KafkaSink{
		Logger:          logger,
		DB:              db,
		Writer:          writer,
		DLQWriter:       dlqWriter,
		Reader:          reader,
		enqueueC:        make(chan kafka.Message, 500), // tuned buffer
		producerWg:      make(chan struct{}, 4),        // numWorkers capacity
		numProducers:    4,
		batchSize:       100,
		flushEvery:      5 * time.Second,
		buffer:          make([]KafkaEvent, 0, 256),
		lastFlush:       time.Now(),
		maxMessageBytes: 1 << 20, // 1 MB
	}
}

// Start producers and consumer. ctx cancellation triggers graceful shutdown.
func (s *KafkaSink) Start(ctx context.Context) {
	// start producer workers
	for i := 0; i < s.numProducers; i++ {
		s.producerWg <- struct{}{} // reserve slot, used to wait
		go s.produceWorker(ctx)
	}

	// start consumer
	go s.consumeAndFlush(ctx)

	// handle shutdown: wait for producers to finish then close writers/readers
	go func() {
		<-ctx.Done()

		// close enqueue channel so producers drain
		close(s.enqueueC)

		// wait for producers to return
		for i := 0; i < s.numProducers; i++ {
			<-s.producerWg
		}

		// close writers/readers (log errors but continue)
		if err := s.Writer.Close(); err != nil {
			s.Logger.Error("failed closing kafka writer", "error", err)
		}
		if err := s.DLQWriter.Close(); err != nil {
			s.Logger.Error("failed closing kafka dlq writer", "error", err)
		}
		if err := s.Reader.Close(); err != nil {
			s.Logger.Error("failed closing kafka reader", "error", err)
		}
	}()
}

// Enqueue queues an event for production to Kafka. Non-blocking with drop-oldest policy.
func (s *KafkaSink) Enqueue(eventType EventType, payload map[string]any) {
	// serialize payload; check size
	data, err := json.Marshal(payload)
	if err != nil {
		s.Logger.Error("failed marshaling event data", "error", err)
		return
	}
	if len(data) > s.maxMessageBytes {
		s.Logger.Warn("payload too large, dropping", "size", len(data))
		return
	}

	msg := kafka.Message{
		Key:   []byte(eventType),
		Value: data,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(eventType)},
			{Key: "timestamp", Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
		},
		Time: time.Now(),
	}

	// non-blocking send with drop-oldest policy
	select {
	case s.enqueueC <- msg:
		// queued
	default:
		// channel full -> drop oldest to make room (policy: drop_oldest)
		select {
		case <-s.enqueueC:
			// removed oldest
		default:
			// if still blocked, log and drop this message
			s.Logger.Warn("enqueue channel full and cannot drop oldest, dropping new message", "type", eventType)
			return
		}
		// try again (should succeed)
		select {
		case s.enqueueC <- msg:
		default:
			s.Logger.Warn("enqueue channel still full after dropping oldest, dropping message", "type", eventType)
		}
	}
}

// produceWorker reads from enqueueC, batches a few messages and writes them to Kafka with retry/backoff.
func (s *KafkaSink) produceWorker(ctx context.Context) {
	defer func() { s.producerWg <- struct{}{} }() // signal finished by writing back to channel

	batch := make([]kafka.Message, 0, 16)
	flushTicker := time.NewTicker(200 * time.Millisecond)
	defer flushTicker.Stop()

	flushBatch := func(ctx context.Context, b []kafka.Message) {
		if len(b) == 0 {
			return
		}
		// Exponential backoff retry
		var attempt int
		for {
			attempt++
			writeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := s.Writer.WriteMessages(writeCtx, b...)
			cancel()
			if err == nil {
				return
			}
			s.Logger.Error("failed writing kafka batch", "err", err, "attempt", attempt, "batch_size", len(b))

			if attempt >= 5 {
				// push batch to DLQ as fallback (individual messages)
				for _, m := range b {
					_ = s.writeToDLQ(context.Background(), m, err)
				}
				return
			}
			// backoff
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	for {
		select {
		case <-ctx.Done():
			// flush remaining
			if len(batch) > 0 {
				flushBatch(ctx, batch)
			}
			return
		case msg, ok := <-s.enqueueC:
			if !ok {
				// channel closed, flush remaining and exit
				if len(batch) > 0 {
					flushBatch(ctx, batch)
				}
				return
			}
			batch = append(batch, msg)
			if len(batch) >= 16 {
				flushBatch(ctx, batch)
				batch = batch[:0]
			}
		case <-flushTicker.C:
			if len(batch) > 0 {
				flushBatch(ctx, batch)
				batch = batch[:0]
			}
		}
	}
}

func (s *KafkaSink) writeToDLQ(ctx context.Context, m kafka.Message, cause error) error {
	// clear topic because writer already has one
	m.Topic = ""

	// add DLQ error reason
	m.Headers = append(m.Headers, kafka.Header{
		Key:   "dlq_reason",
		Value: []byte(cause.Error()),
	})

	writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.DLQWriter.WriteMessages(writeCtx, m); err != nil {
		s.Logger.Error("failed writing to dlq", "error", err)
		return err
	}

	s.Logger.Info("wrote message to dlq", "key", string(m.Key))
	return nil
}

// consumeAndFlush consumes messages from Kafka, buffers them, writes to DB in a transaction,
// then commits Kafka offsets only after DB commit succeeds.
func (s *KafkaSink) consumeAndFlush(ctx context.Context) {
	ticker := time.NewTicker(s.flushEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// flush remaining then return
			if len(s.buffer) > 0 {
				if err := s.flushBatchWithRetry(context.Background()); err != nil {
					s.Logger.Error("failed flushing final batch", "error", err)
				}
			}
			return

		case <-ticker.C:
			if len(s.buffer) > 0 && time.Since(s.lastFlush) >= s.flushEvery {
				if err := s.flushBatchWithRetry(ctx); err != nil {
					s.Logger.Error("failed flushing batch on ticker", "batch_size", len(s.buffer), "error", err)
				}
			}

		default:
			// read with a short timeout so we can periodically flush on time
			readCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			msg, err := s.Reader.FetchMessage(readCtx)
			cancel()

			if err != nil {
				// normal timeout
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					continue
				}
				// network/other error -> log and continue (reader may retry internally)
				s.Logger.Error("failed reading kafka message", "error", err)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			var payload map[string]any
			if err := json.Unmarshal(msg.Value, &payload); err != nil {
				s.Logger.Error("failed unmarshaling kafka message", "error", err, "key", string(msg.Key))
				// commit offset to avoid reprocessing malformed payloads
				if err := s.Reader.CommitMessages(ctx, msg); err != nil {
					s.Logger.Error("failed committing offset after unmarshal error", "error", err)
				}
				continue
			}

			s.buffer = append(s.buffer, KafkaEvent{Message: msg, Payload: payload})

			// Flush when batch size reached
			if len(s.buffer) >= s.batchSize {
				if err := s.flushBatchWithRetry(ctx); err != nil {
					s.Logger.Error("failed flushing batch on size threshold", "batch_size", len(s.buffer), "error", err)
					// keep buffer for retry
				}
			}
		}
	}
}

// flushBatchWithRetry wraps flushBatch with retry/backoff and DLQ behavior for permanent failure.
func (s *KafkaSink) flushBatchWithRetry(ctx context.Context) error {
	if len(s.buffer) == 0 {
		return nil
	}

	var attempt int
	for {
		attempt++
		err := s.flushBatch(ctx)
		if err == nil {
			return nil
		}

		s.Logger.Error("flushBatch failed", "err", err, "attempt", attempt)
		if attempt >= 5 {
			// permanent failure -> push messages to DLQ
			for _, ke := range s.buffer {
				_ = s.writeToDLQ(context.Background(), ke.Message, err)
				// commit the offset to avoid reprocessing
				_ = s.Reader.CommitMessages(context.Background(), ke.Message)
			}
			// clear buffer
			s.buffer = s.buffer[:0]
			s.lastFlush = time.Now()
			return err
		}
		// backoff
		time.Sleep(time.Duration(attempt) * time.Second)
	}
}

// flushBatch writes buffered events into the DB inside a single transaction,
// then commits Kafka offsets only after successful commit.
// It returns an error if DB commit or offset commit fails.
func (s *KafkaSink) flushBatch(ctx context.Context) error {
	if len(s.buffer) == 0 {
		return nil
	}
	batchSize := len(s.buffer)
	s.Logger.Info("flushing webhook batch", "size", batchSize)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var messagesToCommit []kafka.Message

	// Prepared sql: ensure your webhook.tbl_event has UNIQUE(event_id) and returns id
	insertEventSQL := `
		INSERT INTO webhook.tbl_event (id, type, payload, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (id) DO NOTHING
		RETURNING id
	`

	insertDeliverySQL := `
		INSERT INTO webhook.tbl_delivery (event_id, subscription_id, attempt_count, status, created_at)
		SELECT $1, s.id, 0, 'pending', NOW()
		FROM webhook.tbl_subscription s
		WHERE s.is_active = TRUE
		  AND $2::text = ANY(s.event_types)
		ON CONFLICT DO NOTHING
	`

	for _, ke := range s.buffer {
		// idempotency: use event_id from payload if present, otherwise generate one
		var eventID string
		if v, ok := ke.Payload["event_id"]; ok {
			if s, ok2 := v.(string); ok2 && s != "" {
				eventID = s
			}
		}
		if eventID == "" {
			eventID = uuid.NewString()
		}

		// insert event - payload as jsonb may be driver-specific; tx.QueryRow should Scan into returned id
		var returnedID string
		if err := tx.QueryRow(ctx, insertEventSQL, eventID, string(ke.Message.Key), ke.Payload).Scan(&returnedID); err != nil {
			// If the error indicates nothing was returned due to ON CONFLICT DO NOTHING, that's ok:
			// we still want to create deliveries using existing event_id.
			// For simplicity, swallow "no rows" and continue using eventID.
			// Note: depending on DB client, Scan may return sql.ErrNoRows
			// Log the error but continue: we still should create deliveries.
			s.Logger.Warn("insert event may have been duplicate or failed to return id", "err", err, "event_id", eventID)
			returnedID = eventID
		}

		// create delivery rows
		if _, err := tx.Exec(ctx, insertDeliverySQL, returnedID, string(ke.Message.Key)); err != nil {
			s.Logger.Error("failed inserting delivery rows", "event_id", returnedID, "err", err)
			// continue: don't abort entire batch because of one subscription insert failure
			continue
		}

		messagesToCommit = append(messagesToCommit, ke.Message)
	}

	// commit DB transaction
	if err := tx.Commit(ctx); err != nil {
		s.Logger.Error("failed committing webhook batch transaction", "error", err)
		return err
	}

	// commit kafka offsets after DB commit
	if len(messagesToCommit) > 0 {
		if err := s.Reader.CommitMessages(ctx, messagesToCommit...); err != nil {
			s.Logger.Error("failed committing kafka offsets", "count", len(messagesToCommit), "error", err)
			// On offset commit failure, we return error so caller can retry; DB already has been committed.
			return err
		}
	}

	// clear buffer and update lastFlush
	s.buffer = s.buffer[:0]
	s.lastFlush = time.Now()
	s.Logger.Info("successfully flushed webhook batch", "size", batchSize)
	return nil
}
