package webhook

import (
	"context"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/pkg/async"
)

type Sink struct {
	Logger *telemetry.Logger
	DB     *database.Database

	worker *async.Worker
}

// NewSink creates a Sink backed by async.Worker
func NewSink(tel *telemetry.Telemetry, db *database.Database) *Sink {
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

	sink := &Sink{
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
func (s *Sink) Start(ctx context.Context) {
	s.worker.Start(ctx)
}

// Enqueue sends an event to the worker
func (s *Sink) Enqueue(eventType EventType, payload map[string]any) {
	s.worker.Enqueue(async.Event(Event{
		Type: eventType,
		Data: payload,
	}))
}

// flush persists a batch of events (called by async.Worker)
func (s *Sink) flush(ctx context.Context, batch []async.Event) error {
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
