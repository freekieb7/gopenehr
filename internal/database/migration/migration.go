package migration

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Migration interface {
	Version() uint64
	Name() string
	Up(ctx context.Context, tx pgx.Tx) error
	Down(ctx context.Context, tx pgx.Tx) error
}
