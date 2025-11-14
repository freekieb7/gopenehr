package handler

import (
	"context"
	"time"

	"github.com/freekieb7/gopenehr/internal/health"
	"github.com/freekieb7/gopenehr/pkg/web/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	HealthChecker *health.Checker
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	healthGroup := app.Group("/health")
	healthGroup.Use(middleware.NoCache)

	healthGroup.Get("/", h.HandleHealth)
}

func (h *Handler) HandleHealth(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	status := h.HealthChecker.CheckHealth(ctx)
	return c.Status(fiber.StatusOK).JSON(status)
}
