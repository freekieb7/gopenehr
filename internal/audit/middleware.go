package audit

import (
	"github.com/gofiber/fiber/v2"
)

func AuditLoggedMiddleware(sink *Sink, resource Resource, action Action) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := sink.NewContext(c, resource, action)

		// Attach to Fiber locals
		c.Locals(string(ContextKey), ctx)

		// Ensure event is always logged
		defer ctx.Commit()

		// Continue request
		return c.Next()
	}
}
