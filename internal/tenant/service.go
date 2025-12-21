package tenant

import (
	"context"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/google/uuid"
)

type Service struct {
	DB *database.Database
}

func NewService(db *database.Database) *Service {
	return &Service{
		DB: db,
	}
}

type Tenant struct {
	ID   uuid.UUID
	Name string
}

type Subscription struct {
	EHRLimit int
}

func (s *Service) CreateTenant(ctx context.Context, name string, subscription Subscription) (Tenant, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return Tenant{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err != database.ErrTxClosed {
			fmt.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	tenant := Tenant{
		Name: name,
	}

	row := tx.QueryRow(ctx, `INSERT INTO tenant.tbl_tenant (name) VALUES ($1) RETURNING id;`, name)
	err = row.Scan(&tenant.ID)
	if err != nil {
		return Tenant{}, fmt.Errorf("failed to insert tenant: %w", err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO tenant.tbl_subscription (tenant_id, ehr_limit, ehr_count) VALUES ($1, $2, $3);`, tenant.ID, subscription.EHRLimit, 0)
	if err != nil {
		return Tenant{}, fmt.Errorf("failed to insert tenant subscription: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Tenant{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tenant, nil
}
