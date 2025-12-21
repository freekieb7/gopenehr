package migration

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var _ Migration = (*SetupTenant)(nil)

type SetupTenant struct{}

func (m *SetupTenant) Version() uint64 {
	return 20251113195003
}

func (m *SetupTenant) Name() string {
	return "Setup Tenant Schema and Tables"
}

func (m *SetupTenant) Up(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `CREATE SCHEMA tenant;`)
	if err != nil {
		return fmt.Errorf("failed to create tenant schema: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE tenant.tbl_tenant (
			id UUID PRIMARY KEY DEFAULT uuidv4(),
			name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tenant table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE tenant.tbl_subscription (
			tenant_id UUID PRIMARY KEY REFERENCES tenant.tbl_tenant(id) ON DELETE CASCADE,
			ehr_count INT NOT NULL,
			ehr_limit INT NOT NULL,
			token_limit INT NOT NULL,
			token_count INT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tenant subscription table: %w", err)
	}

	return nil
}

func (m *SetupTenant) Down(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `DROP SCHEMA tenant CASCADE;`)
	if err != nil {
		return fmt.Errorf("failed to drop tenant schema: %w", err)
	}

	return nil
}
