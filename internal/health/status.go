package health

import "time"

var (
	// ServiceStatusOK indicates that the service is healthy.
	ServiceStatusHealthy   string = "healthy"
	ServiceStatusUnhealthy string = "unhealthy"
)

type Status struct {
	Status     string                     `json:"status"`
	Timestamp  time.Time                  `json:"timestamp"`
	Version    string                     `json:"version"`
	Components map[string]ComponentStatus `json:"components"`
}

type ComponentStatus struct {
	Status      string    `json:"status"`
	Message     string    `json:"message,omitempty"`
	LatencyMs   int64     `json:"latency_ms,omitempty"`
	LastChecked time.Time `json:"last_checked,omitempty"`
	Critical    bool      `json:"critical"`
}
