package audit

import (
	"context"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/pkg/async"
	"github.com/freekieb7/gopenehr/pkg/audit"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Sink struct {
	Logger *telemetry.Logger
	DB     *database.Database

	worker *async.Worker[audit.Event]
}

// NewSink creates a Sink backed by async.Worker
func NewSink(tel *telemetry.Telemetry, db *database.Database) *Sink {
	// Create worker metrics
	queueDepth, _ := tel.Metrics.Meter.Int64UpDownCounter("audit_queue_depth")
	droppedEvents, _ := tel.Metrics.Meter.Int64Counter("audit_dropped_total")
	flushFailures, _ := tel.Metrics.Meter.Int64Counter("audit_flush_failures")
	batchSize, _ := tel.Metrics.Meter.Int64Histogram("audit_flush_batch_size")

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
	sink.worker = async.NewWorker(cfg, func(ctx context.Context, batch []audit.Event) error {
		return sink.flush(ctx, batch)
	})

	return sink
}

// Start launches the worker
func (s *Sink) Start(ctx context.Context) {
	s.worker.Start(ctx)
}

// Enqueue sends an event to the worker
func (s *Sink) Enqueue(event audit.Event) {
	s.worker.Enqueue(event)
}

// flush persists a batch of events (called by async.Worker)
func (s *Sink) flush(ctx context.Context, batch []audit.Event) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }() // safe if already committed

	// var eventID string

	rows := make([][]interface{}, len(batch))
	for i, event := range batch {
		id := uuid.New() // replace with v7 generator if needed
		rows[i] = []interface{}{
			id,
			event.ActorID,
			event.ActorType,
			event.Resource,
			event.Action,
			event.Success,
			event.IPAddress,
			event.UserAgent,
			event.Details,
			event.CreatedAt,
		}
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"audit", "tbl_audit_log"},
		[]string{"id", "actor_id", "actor_type", "resource", "action", "success", "ip_address", "user_agent", "details", "created_at"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
