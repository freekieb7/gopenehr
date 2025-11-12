package handler

import (
	"context"
	"time"

	"github.com/freekieb7/gopenehr/internal/health"
	"github.com/freekieb7/gopenehr/internal/web"
	"github.com/freekieb7/gopenehr/internal/web/middleware"
	"github.com/gofiber/fiber/v2"
)

type Health struct {
	HealthChecker *health.Checker
}

func (h *Health) RegisterRoutes(app *web.Server) {
	healthGroup := app.Fiber.Group("/health")
	healthGroup.Use(middleware.NoCache)

	healthGroup.Get("/", h.HandleHealth)
}

func (h *Health) HandleHealth(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	status := h.HealthChecker.CheckHealth(ctx)
	return c.Status(fiber.StatusOK).JSON(status)
}
