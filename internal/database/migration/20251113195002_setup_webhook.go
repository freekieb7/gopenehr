package migration

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var _ Migration = (*SetupWebhook)(nil)

type SetupWebhook struct{}

func (m *SetupWebhook) Version() uint64 {
	return 20251113195002
}

func (m *SetupWebhook) Name() string {
	return "Setup Webhook Schema and Tables"
}

func (m *SetupWebhook) Up(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `CREATE SCHEMA webhook;`)
	if err != nil {
		return fmt.Errorf("failed to create webhook schema: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE webhook.tbl_event (
			id UUID PRIMARY KEY DEFAULT uuidv4(),
			type TEXT NOT NULL,  -- e.g., 'ehr_created', 'ehr_deleted'
			payload JSONB NOT NULL,    -- The payload sent to the webhook
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create webhook event table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE webhook.tbl_subscription (
			id UUID PRIMARY KEY DEFAULT uuidv4(),
			url TEXT NOT NULL,
			secret TEXT NOT NULL,      
			event_types TEXT[] NOT NULL,
			is_active BOOLEAN NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE (url)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create webhook subscription table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE webhook.tbl_delivery (
			id UUID PRIMARY KEY DEFAULT uuidv4(),
			event_id UUID NOT NULL REFERENCES webhook.tbl_event(id) ON DELETE CASCADE,
			subscription_id UUID NOT NULL REFERENCES webhook.tbl_subscription(id) ON DELETE CASCADE,
			status TEXT NOT NULL,
			attempt_count INT NOT NULL,
			next_attempt_at TIMESTAMPTZ,
			last_attempt_at TIMESTAMPTZ,
			last_response_code INT,
			last_response_body TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE (event_id, subscription_id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create webhook delivery table: %w", err)
	}

	return nil
}

func (m *SetupWebhook) Down(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `DROP SCHEMA webhook CASCADE;`)
	if err != nil {
		return fmt.Errorf("failed to drop webhook schema: %w", err)
	}

	return nil
}
