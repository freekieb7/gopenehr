package async

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/freekieb7/gopenehr/internal/telemetry"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Event represents a generic unit of work
type Event any

// FlushFunc is a callback to persist a batch of events
type FlushFunc func(ctx context.Context, batch []Event) error

// WorkerConfig defines the behavior of the async worker
type WorkerConfig struct {
	Logger     *telemetry.Logger
	QueueSize  int
	BatchSize  int
	FlushEvery time.Duration
	Metrics    *WorkerMetrics
	Tracer     trace.Tracer // NEW: OTEL tracer for batch flush spans
}

// WorkerMetrics holds OTEL metrics for a worker
type WorkerMetrics struct {
	QueueDepth    metric.Int64UpDownCounter
	DroppedEvents metric.Int64Counter
	FlushFailures metric.Int64Counter
	BatchSize     metric.Int64Histogram
}

// Worker is a reusable async batch worker
type Worker struct {
	cfg     WorkerConfig
	queue   chan Event
	dropped atomic.Uint64
	flushFn FlushFunc
}

// NewWorker creates a new async batch worker
func NewWorker(cfg WorkerConfig, flushFn FlushFunc) *Worker {
	return &Worker{
		cfg:     cfg,
		queue:   make(chan Event, cfg.QueueSize),
		flushFn: flushFn,
	}
}

// Enqueue adds an event to the worker queue
func (w *Worker) Enqueue(ev Event) {
	select {
	case w.queue <- ev:
		if w.cfg.Metrics != nil {
			w.cfg.Metrics.QueueDepth.Add(context.Background(), 1)
		}
	default:
		w.dropped.Add(1)
		if w.cfg.Metrics != nil {
			w.cfg.Metrics.DroppedEvents.Add(context.Background(), 1)
		}
		w.cfg.Logger.Warn("queue full, dropping event", "dropped_total", w.dropped.Load())
	}
}

// Start launches the background worker
func (w *Worker) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.cfg.FlushEvery)
		defer ticker.Stop()

		batch := make([]Event, 0, w.cfg.BatchSize)

		flush := func() {
			if len(batch) == 0 {
				return
			}

			if w.cfg.Metrics != nil {
				w.cfg.Metrics.BatchSize.Record(context.Background(), int64(len(batch)))
			}

			// // ---- TRACE SPAN START ----
			// var span trace.Span
			// if w.cfg.Tracer != nil {
			// 	ctx, span = w.cfg.Tracer.Start(ctx, "worker.flush.batch",
			// 		trace.WithAttributes(attribute.Int("batch.size", len(batch))),
			// 	)
			// 	defer span.End()
			// }

			w.flushWithRetry(ctx, batch)

			// if span != nil {
			// 	span.SetAttributes(attribute.Int("batch.final_size", len(batch)))
			// }

			batch = batch[:0]
		}

		for {
			select {
			case ev := <-w.queue:
				batch = append(batch, ev)
				if w.cfg.Metrics != nil {
					w.cfg.Metrics.QueueDepth.Add(context.Background(), -1)
				}
				if len(batch) >= w.cfg.BatchSize {
					flush()
				}
			case <-ticker.C:
				flush()
			case <-ctx.Done():
				timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				w.drainAndFlush(timeoutCtx, batch)
				return
			}
		}
	}()
}

// flushWithRetry retries the flush function up to 3 times with exponential backoff
func (w *Worker) flushWithRetry(ctx context.Context, batch []Event) {
	backoff := 150 * time.Millisecond
	for i := 0; i < 3; i++ {
		if err := w.flushFn(ctx, batch); err == nil {
			return
		} else {
			if w.cfg.Tracer != nil {
				trace.SpanFromContext(ctx).RecordError(err)
			}
			w.cfg.Logger.Warn("flush failed, retrying", "attempt", i+1, "error", err)
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	if w.cfg.Metrics != nil {
		w.cfg.Metrics.FlushFailures.Add(ctx, int64(len(batch)))
	}

	w.cfg.Logger.Error("batch permanently failed", "count", len(batch))
}

// drainAndFlush drains the queue and flushes remaining events
func (w *Worker) drainAndFlush(ctx context.Context, batch []Event) {
	for {
		select {
		case ev := <-w.queue:
			batch = append(batch, ev)
		default:
			w.flushWithRetry(ctx, batch)
			return
		}
	}
}
