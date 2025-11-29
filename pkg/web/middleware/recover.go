package middleware

import (
	"fmt"
	"log/slog"

	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

func Recover(tel *telemetry.Telemetry) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				ctx := c.UserContext()

				if span := trace.SpanFromContext(ctx); span != nil {
					span.RecordError(fmt.Errorf("%v", r))
				}

				tel.Logger.Error("panic", slog.Any("panic", r), slog.String("path", c.Path()))

				err := c.Status(500).JSON(fiber.Map{
					"error": "internal_error",
				})
				if err != nil {
					tel.Logger.Error("failed to send panic response", slog.Any("error", err))
				}
			}
		}()
		return c.Next()
	}
}
