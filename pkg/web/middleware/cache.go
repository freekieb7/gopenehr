package middleware

import "github.com/gofiber/fiber/v2"

func NoCache(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")
	c.Set("Surrogate-Control", "no-store")
	return c.Next()
}
