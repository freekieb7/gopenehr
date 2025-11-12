package health

import (
	"context"
	"time"
)

type Checker struct {
	Version string
}

func (c *Checker) CheckHealth(ctx context.Context) Status {
	return Status{
		Status:     ServiceStatusHealthy,
		Timestamp:  time.Now(),
		Version:    c.Version,
		Components: map[string]ComponentStatus{},
	}
}
