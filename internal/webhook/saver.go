package webhook

import (
	"context"
	"log/slog"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
)

type Saver struct {
	Logger *slog.Logger
	DB     *database.Database

	Queue      chan Event
	BatchSize  int
	FlushEvery time.Duration

	dropped uint64
}

func NewSaver(logger *slog.Logger, db *database.Database) Saver {
	return Saver{
		Logger:     logger,
		DB:         db,
		Queue:      make(chan Event, 200),
		BatchSize:  100,
		FlushEvery: 1 * time.Second,
	}
}

func (s *Saver) Start(ctx context.Context) {
	go s.worker(ctx)
}

func (s *Saver) Enqueue(eventType EventType, payload map[string]any) {
	select {
	case s.Queue <- Event{
		Type: eventType,
		Data: payload,
	}:
	default:
		s.dropped++
		s.Logger.Warn("webhook queue full, dropping event", "type", eventType, "dropped_total", s.dropped)
	}
}

func (s *Saver) worker(ctx context.Context) {
	ticker := time.NewTicker(s.FlushEvery)
	defer ticker.Stop()

	batch := make([]Event, 0, s.BatchSize)

	flush := func() {
		if len(batch) > 0 {
			s.flushWithRetry(ctx, batch)
			batch = batch[:0]
		}
	}

	for {
		select {
		case ev := <-s.Queue:
			batch = append(batch, ev)
			if len(batch) == s.BatchSize {
				flush()
			}

		case <-ticker.C:
			flush()

		case <-ctx.Done():
			s.drainAndFlush(ctx, batch)
			return
		}
	}
}

func (s *Saver) flushWithRetry(ctx context.Context, batch []Event) {
	backoff := 150 * time.Millisecond

	for i := 0; i < 3; i++ {
		if err := s.flush(ctx, batch); err == nil {
			return
		}

		s.Logger.Warn("flush failed, retrying", "attempt", i+1)
		time.Sleep(backoff)
		backoff *= 2
	}

	s.Logger.Error("batch permanently failed", "count", len(batch))
}

func (s *Saver) drainAndFlush(ctx context.Context, batch []Event) {
	for {
		select {
		case ev := <-s.Queue:
			batch = append(batch, ev)
		default:
			s.flushWithRetry(ctx, batch)
			return
		}
	}
}

func (s *Saver) flush(ctx context.Context, batch []Event) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx) // safe when already committed
	}()

	var eventID string

	for _, event := range batch {
		// create event
		if err := tx.QueryRow(ctx, `
			INSERT INTO webhook.tbl_event (type, payload)
			VALUES ($1, $2)
			RETURNING id
		`, event.Type, event.Data).Scan(&eventID); err != nil {

			s.Logger.Error("failed inserting webhook event", "error", err)
			continue
		}

		// create delivery rows
		if _, err := tx.Exec(ctx, `
			INSERT INTO webhook.tbl_delivery 
				(event_id, subscription_id, attempt_count, status, created_at)
			SELECT e.id, s.id, 0, 'pending', NOW()
			FROM webhook.tbl_event e
			CROSS JOIN webhook.tbl_subscription s
			WHERE e.id = $1
			  AND s.is_active = TRUE
			  AND $2 = ANY(s.event_types)
		`, eventID, event.Type); err != nil {

			s.Logger.Error("failed inserting delivery rows", "event_id", eventID, "error", err)
		}
	}

	return tx.Commit(ctx)
}
