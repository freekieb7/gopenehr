package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/telemetry"
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

type Sender struct {
	Logger *telemetry.Logger
	DB     *database.Database
	Client *http.Client

	jobs   []DeliveryJob
	jobIds []string
}

func NewSender(logger *telemetry.Logger, db *database.Database, client *http.Client) *Sender {
	if client.Timeout == 0 {
		client.Timeout = 10 * time.Second
	}

	return &Sender{
		Logger: logger,
		DB:     db,
		Client: client,
		jobs:   make([]DeliveryJob, BatchSize),
		jobIds: make([]string, BatchSize),
	}
}

func (s *Sender) Start(ctx context.Context) error {
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	s.Logger.Info("Webhook worker started")

	for {
		select {
		case <-ctx.Done():
			s.Logger.Info("Webhook worker shutting down")
			return nil

		case <-ticker.C:
			if err := s.processBatch(ctx); err != nil {
				s.Logger.ErrorContext(ctx, "Batch processing failed", "error", err)
			}
		}
	}
}

func (s *Sender) processBatch(ctx context.Context) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
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

	n := 0
	for rows.Next() {
		if err := rows.Scan(&s.jobs[n].ID, &s.jobs[n].EventID, &s.jobs[n].EventType, &s.jobs[n].Payload, &s.jobs[n].URL, &s.jobs[n].Secret, &s.jobs[n].AttemptCount); err != nil {
			return err
		}
		s.jobIds[n] = s.jobs[n].ID
		n++
	}

	if n == 0 {
		return nil
	}

	_, err = tx.Exec(ctx, `
		UPDATE webhook.tbl_delivery
		SET status = 'processing', updated_at = NOW()
		WHERE id = ANY($1) 
	`, s.jobIds[:n])
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	var wg sync.WaitGroup
	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for idx := range n {
		wg.Go(func() {
			if err := s.handleJob(ctx, s.jobs[idx]); err != nil {
				s.Logger.ErrorContext(ctx, "Job failed", "jobID", s.jobs[idx].ID, "error", err)
			}
		})
	}

	wg.Wait()

	return nil
}

func (s *Sender) handleJob(ctx context.Context, job DeliveryJob) error {
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

	resp, err := s.Client.Do(req)
	if err != nil {
		return s.fail(ctx, job, "network_error", nil, err)
	}
	defer func() {
		_, err = io.Copy(io.Discard, resp.Body)
		if err != nil {
			s.Logger.Warn("Failed to drain response body", "jobID", job.ID, "error", err)
		}
		err = resp.Body.Close()
		if err != nil {
			s.Logger.Warn("Failed to close response body", "jobID", job.ID, "error", err)
		}
	}()

	if resp.StatusCode >= 500 {
		return s.fail(ctx, job, "server_error", resp, nil)
	}

	if resp.StatusCode >= 400 {
		return s.fail(ctx, job, "client_error", resp, nil)
	}

	return s.success(ctx, job, resp)
}

func (s *Sender) success(ctx context.Context, job DeliveryJob, resp *http.Response) error {
	_, err := s.DB.Exec(ctx, `
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

func (s *Sender) fail(ctx context.Context, job DeliveryJob, reason string, resp *http.Response, cause error) error {
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

	s.Logger.Warn("Webhook failed",
		"jobID", job.ID,
		"reason", reason,
		"attempt", job.AttemptCount,
		"status", status,
		"error", cause,
	)

	_, err := s.DB.Exec(ctx, `
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
