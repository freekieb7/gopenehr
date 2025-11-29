package telemetry

import "go.opentelemetry.io/otel/metric"

type Metrics struct {
	Meter    metric.Meter
	Requests metric.Int64Counter
	Duration metric.Float64Histogram
}

func NewMetrics(m metric.Meter) *Metrics {
	req, _ := m.Int64Counter("requests_total")
	dur, _ := m.Float64Histogram("request_duration_seconds")

	return &Metrics{
		Meter:    m,
		Requests: req,
		Duration: dur,
	}
}
