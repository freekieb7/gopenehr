package cli

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/freekieb7/gopenehr/internal/database"
)

type Migrator struct {
	DB     *database.Database
	Logger *slog.Logger
}

func NewMigrator(db *database.Database, logger *slog.Logger) *Migrator {
	return &Migrator{
		DB:     db,
		Logger: logger,
	}
}

func (m *Migrator) Run(ctx context.Context, args []string) error {
	cmd := "help"
	if len(args) >= 1 {
		cmd = args[0]
	}

	switch cmd {
	case "up":
		step := 0
		if len(args) >= 2 {
			var err error
			step, err = strconv.Atoi(args[1])
			if err != nil || step < 0 {
				return fmt.Errorf("invalid step value: %s", args[1])
			}
		}
		return m.MigrateUp(ctx, step)
	case "down":
		step := 0
		if len(args) >= 2 {
			var err error
			step, err = strconv.Atoi(args[1])
			if err != nil || step < 0 {
				return fmt.Errorf("invalid step value: %s", args[1])
			}
		}
		return m.MigrateDown(ctx, step)
	case "help":
		fallthrough
	default:
		fmt.Println("Usage: migrate [command] [step]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  up [step]     - Apply pending migrations (all or specified number)")
		fmt.Println("  down [step]   - Rollback applied migrations (all or specified number)")
		return nil
	}
}

func (m *Migrator) CreateMigrationTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS public.tbl_migration (
		version BIGINT PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TIMESTAMP NOT NULL
	);`

	if _, err := m.DB.Exec(ctx, query); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

func (m *Migrator) MigrateUp(ctx context.Context, step int) error {
	if err := m.CreateMigrationTable(ctx); err != nil {
		return fmt.Errorf("failed to ensure migrations table exists: %w", err)
	}

	// Apply migrations in order
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			m.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Apply migrations in order
	count := 0
	for _, migration := range m.DB.Migrations {
		if step > 0 && count >= step {
			break
		}

		var applied bool
		row := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM public.tbl_migration WHERE version=$1)", migration.Version())
		err = row.Scan(&applied)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", migration.Name(), err)
		}

		if applied {
			m.Logger.InfoContext(ctx, "Skipping migration (already applied)", slog.String("migration", migration.Name()))
			continue
		}

		err = migration.Up(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Name(), err)
		}

		_, err = tx.Exec(ctx, "INSERT INTO public.tbl_migration (version, name, applied_at) VALUES ($1, $2, NOW())", migration.Version(), migration.Name())
		if err != nil {
			return fmt.Errorf("failed to record applied migration %s: %w", migration.Name(), err)
		}

		m.Logger.InfoContext(ctx, "Applied migration successfully", slog.String("migration", migration.Name()))
		count++
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	if count > 0 {
		m.Logger.InfoContext(ctx, "Applied migrations successfully",
			slog.Int("count", count))
	} else {
		m.Logger.InfoContext(ctx, "No migrations to apply")
	}

	return nil
}

func (m *Migrator) MigrateDown(ctx context.Context, step int) error {
	if err := m.CreateMigrationTable(ctx); err != nil {
		return fmt.Errorf("failed to ensure migrations table exists: %w", err)
	}

	// Get applied migrations from database in reverse order (newest first)
	rows, err := m.DB.Query(ctx, "SELECT version, name FROM public.tbl_migration ORDER BY version DESC")
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()

	appliedVersions := map[uint64]string{}
	for rows.Next() {
		var version uint64
		var name string
		if err := rows.Scan(&version, &name); err != nil {
			return fmt.Errorf("failed to scan migration version and name: %w", err)
		}
		appliedVersions[version] = name
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating migrations: %w", err)
	}

	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			m.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Roll back migrations in reverse order, reading only needed files
	count := 0
	for version, name := range appliedVersions {
		if step > 0 && count >= step {
			break
		}

		var applied bool
		for _, migration := range m.DB.Migrations {
			if migration.Version() != version {
				continue
			}

			err = migration.Down(ctx, tx)
			if err != nil {
				return fmt.Errorf("failed to roll back migration %s: %w", name, err)
			}

			_, err = m.DB.Exec(ctx, "DELETE FROM public.tbl_migration WHERE version=$1", version)
			if err != nil {
				return fmt.Errorf("failed to remove rolled back migration %s: %w", name, err)
			}

			m.Logger.InfoContext(ctx, "Rolled back migration successfully", slog.String("migration", name))
			applied = true
		}

		if !applied {
			return fmt.Errorf("migration %s (version %d) not found in migration list", name, version)
		}

		count++
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	if count > 0 {
		m.Logger.InfoContext(ctx, "Rolled back migrations successfully",
			slog.Int("count", count))
	} else {
		m.Logger.InfoContext(ctx, "No migrations to rollback")
	}

	return nil
}
