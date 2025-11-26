package handler

import (
	"log/slog"
	"strconv"

	"github.com/freekieb7/gopenehr/internal/audit"
	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Logger       *slog.Logger
	AuditService *audit.Service
}

func NewHandler(logger *slog.Logger, auditService *audit.Service) Handler {
	return Handler{
		Logger:       logger,
		AuditService: auditService,
	}
}

func (h *Handler) RegisterRoutes(c *fiber.App) {
	auditGroup := c.Group("/audit/v1")
	auditGroup.Get("/logs", h.ListLogEntries)
}

func (h *Handler) ListLogEntries(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse pagination parameters
	pageSize := 25 // default
	if pageSizeParam := c.Query("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	token := c.Query("token")

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHR,
			Action:    audit.ActionDelete,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"page_size": pageSize,
				"token":     token,
				"outcome":   outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Delete Multiple EHRs", "error", err)
		}
	}()

	// Use the new paginated method
	listReq := audit.ListLogEntriesRequest{
		PageSize: pageSize,
		Token:    token,
	}

	result, err := h.AuditService.ListLogEntriesPaginated(ctx, listReq)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to list log entries", "error", err)
		outcome = "failure"
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list log entries",
		})
	}

	outcome = "success"
	return c.JSON(fiber.Map{
		"entries": result.LogEntries,
		"token":   result.NextToken,
	})
}
