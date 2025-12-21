package migration

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var _ Migration = (*SetupAudit)(nil)

type SetupAudit struct{}

func (m *SetupAudit) Version() uint64 {
	return 20251113195001
}

func (m *SetupAudit) Name() string {
	return "Setup Audit Schema and Tables"
}

func (m *SetupAudit) Up(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `CREATE SCHEMA audit;`)
	if err != nil {
		return fmt.Errorf("failed to create audit schema: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE audit.tbl_audit_log (
			id UUID PRIMARY KEY,
			data JSONB NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create audit log table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE INDEX idx_audit_log_created_at ON audit.tbl_audit_log USING btree (created_at DESC);
	`)
	if err != nil {
		return fmt.Errorf("failed to create index on audit log table: %w", err)
	}

	return nil
}

func (m *SetupAudit) Down(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `DROP SCHEMA audit CASCADE;`)
	if err != nil {
		return fmt.Errorf("failed to drop audit schema: %w", err)
	}

	return nil
}
