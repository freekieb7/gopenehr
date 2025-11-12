package health

import (
	"context"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
)

type Checker struct {
	Version string
	DB      *database.Database
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

	return ComponentStatus{
		Status:              status,
		LatencyMicroSeconds: uint64(latency.Microseconds()),
		LastChecked:         time.Now().UTC(),
		Critical:            true,
		Message:             message,
	}
}
