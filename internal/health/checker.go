package health

import (
	"context"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
)

type Checker struct {
	Version                string
	TargetMigrationVersion uint64
	DB                     *database.Database
}

func NewChecker(version string, targetMigrationVersion uint64, db *database.Database) Checker {
	return Checker{
		Version:                version,
		TargetMigrationVersion: targetMigrationVersion,
		DB:                     db,
	}
}

func (c *Checker) CheckHealth(ctx context.Context) Status {
	status := ServiceStatusHealthy
	components := make(map[string]ComponentStatus)

	// Check database health
	dbStatus := c.CheckDatabaseHealth(ctx)
	components["database"] = dbStatus

	for _, component := range components {
		if component.Critical && component.Status == ServiceStatusUnhealthy {
			status = ServiceStatusUnhealthy
			break
		}
	}

	return Status{
		Status:     status,
		Timestamp:  time.Now().UTC(),
		Version:    c.Version,
		Components: components,
	}
}

func (c *Checker) CheckDatabaseHealth(ctx context.Context) ComponentStatus {
	status := ServiceStatusHealthy
	message := "database connection healthy"
	start := time.Now()

	if err := c.DB.Ping(ctx); err != nil {
		status = ServiceStatusUnhealthy
		message = "database connection failed: " + err.Error()
	}

	latency := time.Since(start)

	// Check if migrations table exists
	err := c.DB.QueryRow(ctx, `SELECT 1 FROM information_schema.tables WHERE table_schema='public' AND table_name='tbl_migration'`).Scan(new(int))
	if err != nil {
		status = ServiceStatusUnhealthy
		message = "migrations table check failed: " + err.Error()
	}

	// Check if there is a migration migration is applied based on current time
	err = c.DB.QueryRow(ctx, `SELECT 1 FROM public.tbl_migration WHERE version = $1 FROM public.tbl_migration)`, c.TargetMigrationVersion).Scan(new(int))
	if err != nil {
		status = ServiceStatusUnhealthy
		message = "required migration not applied"
	}

	return ComponentStatus{
		Status:              status,
		LatencyMicroSeconds: uint64(latency.Microseconds()),
		LastChecked:         time.Now().UTC(),
		Critical:            true,
		Message:             message,
	}
}
