package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TenantID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Get("X-Tenant-ID", uuid.Nil.String())
		id, err := uuid.Parse(tenantID)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		c.Locals("tenant_id", id)
		return c.Next()
	}
}

func TenantFrom(c *fiber.Ctx) uuid.UUID {
	tenantID, ok := c.Locals("tenant_id").(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return tenantID
}
