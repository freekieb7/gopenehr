package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/freekieb7/gopenehr/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoRows = sql.ErrNoRows
)

type Database struct {
	*pgxpool.Pool
}

func NewDatabase() Database {
	return Database{}
}

func (db *Database) Connect(ctx context.Context, dsn string) error {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse postgres dsn: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	db.Pool = pool

	return nil
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

type ActiveQuery struct {
	ID        uuid.UUID
	Name      string
	AQL       string
	TableName string
	LastRun   time.Time
	NextRun   time.Time
	CreatedAt time.Time
}

func (db *Database) Migrate(ctx context.Context) error {
	// todo
	return nil
}

type CreateActiveQueryParams struct {
	Name string
	AQL  string
	SQL  string
}

func (db *Database) CreateActiveQuery(ctx context.Context, params CreateActiveQueryParams) (ActiveQuery, error) {
	activeQuery := ActiveQuery{
		ID:        uuid.New(),
		Name:      params.Name,
		AQL:       params.AQL,
		TableName: "active_query_" + util.RandomLowerAlphaString(10),
		LastRun:   time.Now(),
		NextRun:   time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return activeQuery, err
	}
	defer tx.Rollback(ctx)

	// Insert active query metadata
	if _, err := tx.Exec(ctx, "INSERT INTO tbl_active_query (id, name, aql, table_name, last_run, next_run, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", activeQuery.ID, activeQuery.Name, activeQuery.AQL, activeQuery.TableName, activeQuery.LastRun, activeQuery.NextRun, activeQuery.CreatedAt); err != nil {
		return activeQuery, err
	}

	if _, err := tx.Exec(ctx, "CREATE MATERIALIZED VIEW "+activeQuery.TableName+" AS "+params.SQL); err != nil {
		return activeQuery, err
	}

	if err := tx.Commit(ctx); err != nil {
		return activeQuery, err
	}

	return activeQuery, nil
}

func (db *Database) GetAllActiveQueries(ctx context.Context) ([]ActiveQuery, error) {
	rows, err := db.Query(ctx, "SELECT id, name, aql, table_name, last_run, next_run, created_at FROM tbl_active_query")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activeQueries := make([]ActiveQuery, 0)
	for rows.Next() {
		var activeQuery ActiveQuery
		if err := rows.Scan(&activeQuery.ID, &activeQuery.Name, &activeQuery.AQL, &activeQuery.TableName, &activeQuery.LastRun, &activeQuery.NextRun, &activeQuery.CreatedAt); err != nil {
			return nil, err
		}
		activeQueries = append(activeQueries, activeQuery)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activeQueries, nil
}

func (db *Database) GetActiveQueryByID(ctx context.Context, id uuid.UUID) (ActiveQuery, error) {
	var activeQuery ActiveQuery
	row := db.QueryRow(ctx, "SELECT id, name, aql, table_name, last_run, next_run, created_at FROM tbl_active_query WHERE id = $1", id)
	if err := row.Scan(&activeQuery.ID, &activeQuery.Name, &activeQuery.AQL, &activeQuery.TableName, &activeQuery.LastRun, &activeQuery.NextRun, &activeQuery.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return activeQuery, ErrNoRows
		}
		return activeQuery, err
	}

	return activeQuery, nil
}

func (db *Database) GetActiveQueryByName(ctx context.Context, name string) (ActiveQuery, error) {
	var activeQuery ActiveQuery
	row := db.QueryRow(ctx, "SELECT id, name, aql, table_name, last_run, next_run, created_at FROM tbl_active_query WHERE name = $1", name)
	if err := row.Scan(&activeQuery.ID, &activeQuery.Name, &activeQuery.AQL, &activeQuery.TableName, &activeQuery.LastRun, &activeQuery.NextRun, &activeQuery.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return activeQuery, ErrNoRows
		}
		return activeQuery, err
	}

	return activeQuery, nil
}

func (db *Database) SyncActiveQuery(ctx context.Context, id uuid.UUID) error {
	activeQuery, err := db.GetActiveQueryByID(ctx, id)
	if err != nil {
		return err
	}

	if _, err := db.Exec(ctx, "REFRESH MATERIALIZED VIEW "+activeQuery.TableName); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, "UPDATE tbl_active_query SET last_run = $1, next_run = $2 WHERE id = $3", time.Now(), time.Now().Add(10*time.Minute), id); err != nil {
		return err
	}

	return nil
}
