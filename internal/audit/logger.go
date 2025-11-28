package audit

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
)

// Logger is an async audit logger.
type Logger struct {
	Logger *slog.Logger
	DB     *database.Database

	Queue      chan Event
	BatchSize  int
	FlushAfter time.Duration

	dropped int64
}

// NewLogger creates a new audit logger with sane defaults.
func NewLogger(db *database.Database) Logger {
	return Logger{
		DB:         db,
		Queue:      make(chan Event, 1000),
		BatchSize:  100,
		FlushAfter: 2 * time.Second,
	}
}

// Log queues an event for asynchronous processing.
func (l *Logger) Log(event Event) {
	select {
	case l.Queue <- event:
	default:
		l.dropped++
		l.Logger.Warn("Event queue is full, dropping event", "event", event, "dropped_total", l.dropped)
	}
}

// Start runs the background worker for flushing the audit queue.
func (l *Logger) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(l.FlushAfter)
		defer ticker.Stop()

		var batch []Event

		flush := func() {
			if len(batch) == 0 {
				return
			}
			l.flushWithRetry(ctx, batch)
			batch = batch[:0]
		}

		for {
			select {
			case e := <-l.Queue:
				batch = append(batch, e)
				if len(batch) >= l.BatchSize {
					flush()
				}
			case <-ticker.C:
				flush()
			case <-ctx.Done():
				l.drainAndFlush(ctx, batch)
				return
			}
		}
	}()
}

// flushWithRetry writes a batch of events to the database with retry/backoff.
func (l *Logger) flushWithRetry(ctx context.Context, batch []Event) {
	if len(batch) == 0 {
		return
	}

	maxRetries := 3
	backoff := 100 * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := l.flush(ctx, batch); err == nil {
			return
		} else {
			l.Logger.Warn("Audit flush failed, retrying", "attempt", attempt, "error", err)
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	l.Logger.Error("Audit flush failed after retries, events may be lost", "count", len(batch))
}

// flush writes a batch of events to the database once.
func (l *Logger) flush(ctx context.Context, batch []Event) error {
	tx, err := l.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, database.ErrTxClosed) {
			l.Logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	stmt, err := tx.Prepare(ctx, "insert_audit_log", `
		INSERT INTO audit.tbl_audit_log
		(actor_id, actor_type, resource, action, success, ip_address, user_agent, details, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`)
	if err != nil {
		return err
	}

	for _, event := range batch {
		_, err := tx.Exec(ctx, stmt.Name,
			event.ActorID,
			event.ActorType,
			event.Resource,
			event.Action,
			event.Success,
			event.IPAddress,
			event.UserAgent,
			event.Details,
			event.CreatedAt,
		)
		if err != nil {
			l.Logger.Error("Failed to execute statement for audit log", "error", err, "event", event)
			continue
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// drainAndFlush empties the queue and flushes remaining events.
func (l *Logger) drainAndFlush(ctx context.Context, batch []Event) {
	for {
		select {
		case e := <-l.Queue:
			batch = append(batch, e)
		default:
			l.flushWithRetry(ctx, batch)
			return
		}
	}
}
