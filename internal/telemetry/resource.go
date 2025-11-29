package telemetry

import (
	sdkResource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewResource(name string, version string) *sdkResource.Resource {
	return sdkResource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(name),
		semconv.ServiceVersionKey.String(version),
	)
}
