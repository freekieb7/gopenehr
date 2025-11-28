package middleware

import (
	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/gofiber/fiber/v2"
)

func APIKeyProtected(secret string) func(c *fiber.Ctx) error {
	enabled := secret != ""

	if !enabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return func(c *fiber.Ctx) error {
		apiKey := c.Get(config.API_KEY_HEADER)
		if apiKey != secret {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		return c.Next()
	}
}
