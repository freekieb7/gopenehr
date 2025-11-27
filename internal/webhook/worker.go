package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
)

const (
	MaxRetries     = 10
	InitialBackoff = time.Minute
	BatchSize      = 50
	PollInterval   = 2 * time.Second
)

type DeliveryJob struct {
	ID           string
	EventID      string
	EventType    string
	Payload      []byte
	URL          string
	Secret       string
	AttemptCount int
}

type Worker struct {
	Logger *slog.Logger
	DB     *database.Database
	Client *http.Client
}

func NewWorker(logger *slog.Logger, db *database.Database, client *http.Client) Worker {
	if client.Timeout == 0 {
		client.Timeout = 10 * time.Second
	}

	return Worker{
		Logger: logger,
		DB:     db,
		Client: client,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	w.Logger.Info("Webhook worker started")

	for {
		select {
		case <-ctx.Done():
			w.Logger.Info("Webhook worker shutting down")
			return nil

		case <-ticker.C:
			if err := w.processBatch(ctx); err != nil {
				w.Logger.ErrorContext(ctx, "Batch processing failed", "error", err)
			}
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) error {
	tx, err := w.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil && err != database.ErrTxClosed {
			w.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	rows, err := tx.Query(ctx, `
		SELECT
			d.id,
			e.id,
			e.type,
			e.payload,
			s.url,
			s.secret,
			d.attempt_count
		FROM webhook.tbl_delivery d
		JOIN webhook.tbl_event e ON e.id = d.event_id
		JOIN webhook.tbl_subscription s ON s.id = d.subscription_id
		WHERE
			d.status IN ('pending','retry')
			AND (d.next_attempt_at IS NULL OR d.next_attempt_at <= now())
		ORDER BY d.created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`, BatchSize)
	if err != nil {
		return err
	}
	defer rows.Close()

	var jobs []DeliveryJob
	for rows.Next() {
		var j DeliveryJob
		if err := rows.Scan(&j.ID, &j.EventID, &j.EventType, &j.Payload, &j.URL, &j.Secret, &j.AttemptCount); err != nil {
			return err
		}
		jobs = append(jobs, j)
	}

	if len(jobs) == 0 {
		return nil
	}

	ids := make([]string, len(jobs))
	for i := range jobs {
		ids[i] = jobs[i].ID
	}

	_, err = tx.Exec(ctx, `
		UPDATE webhook.tbl_delivery
		SET status = 'processing', updated_at = NOW()
		WHERE id = ANY($1) 
	`, ids)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	for _, job := range jobs {
		if err := w.handleJob(ctx, job); err != nil {
			w.Logger.ErrorContext(ctx, "Job failed", "jobID", job.ID, "error", err)
		}
	}

	return nil
}

func (w *Worker) handleJob(ctx context.Context, job DeliveryJob) error {
	sig := sign(job.Secret, job.Payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, job.URL, bytes.NewReader(job.Payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Event", job.EventType)
	req.Header.Set("X-Webhook-ID", job.EventID)
	req.Header.Set("X-Webhook-Signature", sig)
	req.Header.Set("Idempotency-Key", job.ID)

	resp, err := w.Client.Do(req)
	if err != nil {
		return w.fail(ctx, job, "network_error", nil, err)
	}
	defer func() {
		_, err = io.Copy(io.Discard, resp.Body)
		if err != nil {
			w.Logger.Warn("Failed to drain response body", "jobID", job.ID, "error", err)
		}
		err = resp.Body.Close()
		if err != nil {
			w.Logger.Warn("Failed to close response body", "jobID", job.ID, "error", err)
		}
	}()

	if resp.StatusCode >= 500 {
		return w.fail(ctx, job, "server_error", resp, nil)
	}

	if resp.StatusCode >= 400 {
		return w.fail(ctx, job, "client_error", resp, nil)
	}

	return w.success(ctx, job, resp)
}

func (w *Worker) success(ctx context.Context, job DeliveryJob, resp *http.Response) error {
	_, err := w.DB.Exec(ctx, `
		UPDATE webhook.tbl_delivery
		SET
			status = 'delivered',
			last_attempt_at = now(),
			last_response_code = $2,
			updated_at = NOW()
		WHERE id = $1
	`, job.ID, resp.StatusCode)
	return err
}

func (w *Worker) fail(ctx context.Context, job DeliveryJob, reason string, resp *http.Response, cause error) error {
	status := "retry"
	if job.AttemptCount+1 >= MaxRetries || reason == "client_error" {
		status = "dead"
	}

	next := time.Now().Add(backoff(job.AttemptCount))

	var code *int
	var body *string

	if resp != nil {
		b, _ := io.ReadAll(resp.Body)
		s := string(b)
		body = &s

		c := resp.StatusCode
		code = &c
	}

	w.Logger.Warn("Webhook failed",
		"jobID", job.ID,
		"reason", reason,
		"attempt", job.AttemptCount,
		"status", status,
		"error", cause,
	)

	_, err := w.DB.Exec(ctx, `
		UPDATE webhook.tbl_delivery
		SET
			attempt_count = attempt_count + 1,
			status = $2,
			next_attempt_at = $3,
			last_attempt_at = now(),
			last_response_code = $4,
			last_response_body = $5,
			updated_at = NOW()
		WHERE id = $1
	`, job.ID, status, next, code, body)

	return err
}

func backoff(retries int) time.Duration {
	if retries > 10 {
		retries = 10
	}
	return InitialBackoff*time.Duration(1<<retries) + jitter()
}

func jitter() time.Duration {
	return time.Duration(rand.Int63n(int64(3 * time.Second)))
}

func sign(secret string, body []byte) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}
