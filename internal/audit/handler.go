package audit

import (
	"strconv"

	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/oauth"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/pkg/web/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Settings     *config.Settings
	Logger       *telemetry.Logger
	AuditService *Service
	OAuthService *oauth.Service
	AuditLogger  *Logger
}

func NewHandler(settings *config.Settings, logger *telemetry.Logger, auditService *Service, oauthService *oauth.Service, auditLogger *Logger) Handler {
	return Handler{
		Settings:     settings,
		Logger:       logger,
		AuditService: auditService,
		OAuthService: oauthService,
		AuditLogger:  auditLogger,
	}
}

func (h *Handler) RegisterRoutes(c *fiber.App) {
	v1 := c.Group("/audit/v1")
	v1.Use(middleware.APIKeyProtected(h.Settings.APIKey))

	v1.Get("/logs",
		middleware.JWTProtected(h.OAuthService, []oauth.Scope{oauth.ScopeAuditRead}), // Example use
		Middleware(h.AuditLogger, ResourceEHR, ActionCreate),
		h.ListLogEntries,
	)
}

func (h *Handler) ListLogEntries(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := From(c)

	h.Logger.InfoContext(ctx, "Listing log entries")

	// Parse pagination parameters
	pageSize := 25 // default
	if pageSizeParam := c.Query("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	token := c.Query("token")

	// Use the new paginated method
	listReq := ListEventsRequest{
		PageSize: pageSize,
		Token:    token,
	}

	result, err := h.AuditService.ListEventsPaginated(ctx, listReq)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to list log entries", "error", err)
		auditCtx.Fail("internal_error", "Failed to list log entries")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list log entries",
		})
	}

	auditCtx.Success()

	return c.JSON(fiber.Map{
		"events": result.Events,
		"token":  result.NextToken,
	})
}
