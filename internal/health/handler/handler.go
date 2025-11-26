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

func NewHandler(healthChecker *health.Checker) Handler {
	return Handler{
		HealthChecker: healthChecker,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	healthGroup := app.Group("/health")
	healthGroup.Use(middleware.NoCache)

	healthGroup.Get("/healthz", h.HandleLiveness)
	healthGroup.Get("/readyz", h.HandleReadiness)
	healthGroup.Get("/startup", h.HandleStartup)
}

// LIVENESS — basic "is process alive?"
func (h *Handler) HandleLiveness(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

// READINESS — checks DB + migration version
func (h *Handler) HandleReadiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	status := h.HealthChecker.CheckHealth(ctx)
	if status.Status != health.ServiceStatusHealthy {
		return c.SendStatus(fiber.StatusServiceUnavailable)
	}
	return c.SendStatus(fiber.StatusOK)
}

// STARTUP — only used if you need a delayed boot phase
func (h *Handler) HandleStartup(c *fiber.Ctx) error {
	return c.SendString("started")
}
