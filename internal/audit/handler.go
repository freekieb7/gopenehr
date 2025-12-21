package audit

import (
	"strconv"

	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/oauth"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/pkg/audit"
	"github.com/freekieb7/gopenehr/pkg/web/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Settings     *config.Settings
	Telemetry    *telemetry.Telemetry
	AuditService *Service
	OAuthService *oauth.Service
	AuditSink    *Sink
}

func NewHandler(settings *config.Settings, telemetry *telemetry.Telemetry, auditService *Service, oauthService *oauth.Service, auditSink *Sink) Handler {
	return Handler{
		Settings:     settings,
		Telemetry:    telemetry,
		AuditService: auditService,
		OAuthService: oauthService,
		AuditSink:    auditSink,
	}
}

func (h *Handler) RegisterRoutes(c *fiber.App) {
	v1 := c.Group("/audit/v1")

	v1.Use(middleware.APIKeyProtected(h.Settings.APIKey))

	v1.Use(middleware.RequestID())
	v1.Use(middleware.Recover(h.Telemetry))
	v1.Use(middleware.Telemetry(h.Telemetry))

	var validateToken middleware.ValidateTokenFunc = nil
	if h.OAuthService.Enabled() {
		validateToken = h.OAuthService.ValidateToken
	}

	v1.Get("/logs",
		middleware.Audit(h.AuditSink.Enqueue, audit.ResourceAudit, audit.ActionRead),
		middleware.JWTProtected([]string{oauth.ScopeAuditRead.String()}, validateToken),
		h.ListLogEntries,
	)
}

func (h *Handler) ListLogEntries(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := middleware.AuditFrom(c)

	h.Telemetry.Logger.InfoContext(ctx, "Listing log entries")

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
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to list log entries", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Status:  "internal_error",
			Message: "Failed to list log entries",
		})
	}

	auditCtx.Success()

	return c.JSON(fiber.Map{
		"events": result.Events,
		"token":  result.NextToken,
	})
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
	Status  string `json:"status"`
	Details any    `json:"details,omitempty"`
}

func SendErrorResponse(c *fiber.Ctx, auditCtx *audit.Context, errorRes ErrorResponse) error {
	auditCtx.Fail(errorRes.Status, errorRes.Message)
	return c.Status(errorRes.Code).JSON(map[string]ErrorResponse{
		"error": errorRes,
	})
}
