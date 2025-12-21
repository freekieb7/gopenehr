package database

import (
	"context"
	"slices"

	"github.com/freekieb7/gopenehr/internal/database/migration"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoRows   = pgx.ErrNoRows
	ErrTxClosed = pgx.ErrTxClosed
)

type Database struct {
	*pgxpool.Pool
	Migrations       []migration.Migration
	MigrationVersion uint64 // Latest migration version
}

func New() *Database {
	migrations := []migration.Migration{
		&migration.SetupOpenEHR{},
		&migration.SetupAudit{},
		&migration.SetupWebhook{},
		&migration.SetupTenant{},
	}
	slices.SortFunc(migrations, func(migration1, migration2 migration.Migration) int {
		if migration1.Version() < migration2.Version() {
			return -1
		} else if migration1.Version() > migration2.Version() {
			return 1
		}
		panic("duplicate migration version detected")
	})

	return &Database{
		Migrations:       migrations,
		MigrationVersion: migrations[len(migrations)-1].Version(),
	}
}

func (db *Database) Connect(ctx context.Context, url string) error {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return err
	}

	// Ping the database to ensure connection is valid
	if err := pool.Ping(ctx); err != nil {
		return err
	}

	db.Pool = pool
	return nil
}

func (db *Database) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
