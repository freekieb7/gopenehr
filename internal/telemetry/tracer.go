package telemetry

import "go.opentelemetry.io/otel/trace"

type Tracing struct {
	Tracer trace.Tracer
}

func NewTracing(t trace.Tracer) *Tracing {
	return &Tracing{Tracer: t}
}
