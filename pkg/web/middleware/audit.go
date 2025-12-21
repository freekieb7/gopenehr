package middleware

import (
	"net"
	"time"

	"github.com/freekieb7/gopenehr/pkg/audit"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const AuditContextKey string = "audit_context"

type SendFunc func(event audit.Event)

func Audit(send SendFunc, resource audit.Resource, action audit.Action) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := net.ParseIP(c.IP())

		auditCtx := audit.Context{
			Event: audit.Event{
				ID:        uuid.New(),
				Resource:  string(resource),
				Action:    string(action),
				IPAddress: ip,
				UserAgent: c.Get("User-Agent"),
				Details:   make(map[string]any),
				CreatedAt: time.Now().UTC(),
			},
		}

		// Attach to Fiber locals
		c.Locals(AuditContextKey, &auditCtx)

		// Ensure event is always logged
		defer func() {
			if !auditCtx.Event.Success {
				if _, ok := auditCtx.Event.Details["outcome"]; !ok {
					auditCtx.Event.Details["outcome"] = "failure"
				}
			}

			// Send event to sink
			send(auditCtx.Event)
		}()

		// Continue request
		return c.Next()
	}
}

func AuditFrom(c *fiber.Ctx) *audit.Context {
	raw := c.Locals(AuditContextKey)
	if raw == nil {
		panic("audit context not found in fiber context")
	}
	_, ok := raw.(*audit.Context)
	if !ok {
		panic("audit context has wrong type")
	}
	return raw.(*audit.Context)
}
