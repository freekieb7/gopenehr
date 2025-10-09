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

type PreparedTable struct {
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

type CreatePreparedTableParams struct {
	Name string
	AQL  string
	SQL  string
}

func (db *Database) CreatePreparedTable(ctx context.Context, params CreatePreparedTableParams) (PreparedTable, error) {
	preparedTable := PreparedTable{
		ID:        uuid.New(),
		Name:      params.Name,
		AQL:       params.AQL,
		TableName: "prepared_table_" + util.RandomLowerAlphaString(10),
		LastRun:   time.Now(),
		NextRun:   time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return preparedTable, err
	}
	defer tx.Rollback(ctx)

	// Insert prepared table metadata
	if _, err := tx.Exec(ctx, "INSERT INTO tbl_prepared_table (id, name, aql, table_name, last_run, next_run, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", preparedTable.ID, preparedTable.Name, preparedTable.AQL, preparedTable.TableName, preparedTable.LastRun, preparedTable.NextRun, preparedTable.CreatedAt); err != nil {
		return preparedTable, err
	}

	if _, err := tx.Exec(ctx, "CREATE MATERIALIZED VIEW "+preparedTable.TableName+" AS "+params.SQL); err != nil {
		return preparedTable, err
	}

	if err := tx.Commit(ctx); err != nil {
		return preparedTable, err
	}

	return preparedTable, nil
}

func (db *Database) GetAllPreparedTables(ctx context.Context) ([]PreparedTable, error) {
	rows, err := db.Query(ctx, "SELECT id, name, aql, table_name, last_run, next_run, created_at FROM tbl_prepared_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	preparedTables := make([]PreparedTable, 0)
	for rows.Next() {
		var preparedTable PreparedTable
		if err := rows.Scan(&preparedTable.ID, &preparedTable.Name, &preparedTable.AQL, &preparedTable.TableName, &preparedTable.LastRun, &preparedTable.NextRun, &preparedTable.CreatedAt); err != nil {
			return nil, err
		}
		preparedTables = append(preparedTables, preparedTable)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return preparedTables, nil
}

func (db *Database) GetPreparedTableByID(ctx context.Context, id uuid.UUID) (PreparedTable, error) {
	var preparedTable PreparedTable
	row := db.QueryRow(ctx, "SELECT id, name, aql, table_name, last_run, next_run, created_at FROM tbl_prepared_table WHERE id = $1", id)
	if err := row.Scan(&preparedTable.ID, &preparedTable.Name, &preparedTable.AQL, &preparedTable.TableName, &preparedTable.LastRun, &preparedTable.NextRun, &preparedTable.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return preparedTable, ErrNoRows
		}
		return preparedTable, err
	}

	return preparedTable, nil
}

func (db *Database) GetPreparedTableByName(ctx context.Context, name string) (PreparedTable, error) {
	var preparedTable PreparedTable
	row := db.QueryRow(ctx, "SELECT id, name, aql, table_name, last_run, next_run, created_at FROM tbl_prepared_table WHERE name = $1", name)
	if err := row.Scan(&preparedTable.ID, &preparedTable.Name, &preparedTable.AQL, &preparedTable.TableName, &preparedTable.LastRun, &preparedTable.NextRun, &preparedTable.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return preparedTable, ErrNoRows
		}
		return preparedTable, err
	}

	return preparedTable, nil
}

func (db *Database) SyncPreparedTable(ctx context.Context, id uuid.UUID) error {
	preparedTable, err := db.GetPreparedTableByID(ctx, id)
	if err != nil {
		return err
	}

	if _, err := db.Exec(ctx, "REFRESH MATERIALIZED VIEW "+preparedTable.TableName); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, "UPDATE tbl_prepared_table SET last_run = $1, next_run = $2 WHERE id = $3", time.Now(), time.Now().Add(10*time.Minute), id); err != nil {
		return err
	}

	return nil
}
