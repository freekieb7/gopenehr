package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoRows   = pgx.ErrNoRows
	ErrTxClosed = pgx.ErrTxClosed
)

type Database struct {
	*pgxpool.Pool
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
