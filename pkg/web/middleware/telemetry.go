package middleware

import (
	"log/slog"
	"time"

	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
)

func Telemetry(tel *telemetry.Telemetry) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		ctx := c.UserContext()

		// ---- Start span ----
		ctx, span := tel.Tracing.Tracer.Start(ctx, c.Method()+" "+c.Path())
		defer span.End()

		// ---- Run business logic ----
		err := c.Next()

		// ---- Response data ----
		duration := time.Since(start).Seconds()
		status := c.Response().StatusCode()
		traceID := span.SpanContext().TraceID().String()

		// ---- Attributes ----
		attrs := []attribute.KeyValue{
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.Int("http.status", status),
			attribute.Float64("duration_s", duration),
		}
		span.SetAttributes(attrs...)

		// ---- Metrics ----
		tel.Metrics.Requests.Add(ctx, 1)
		tel.Metrics.Duration.Record(ctx, duration)

		// ---- Logging ----
		tel.Logger.Info("request",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.Float64("duration_s", duration),
			slog.String("trace_id", traceID),
		)

		return err
	}
}
