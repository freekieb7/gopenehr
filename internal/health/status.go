package health

import "time"

type ServiceStatus string

var (
	// ServiceStatusOK indicates that the service is healthy.
	ServiceStatusHealthy   ServiceStatus = "healthy"
	ServiceStatusUnhealthy ServiceStatus = "unhealthy"
)

type Status struct {
	Status     ServiceStatus              `json:"status"`
	Timestamp  time.Time                  `json:"timestamp"`
	Version    string                     `json:"version"`
	Components map[string]ComponentStatus `json:"components"`
}

type ComponentStatus struct {
	Status              ServiceStatus `json:"status"`
	Message             string        `json:"message,omitempty"`
	LatencyMicroSeconds uint64        `json:"latency_microseconds,omitempty"`
	LastChecked         time.Time     `json:"last_checked,omitempty"`
	Critical            bool          `json:"critical"`
}
