package tenant

import (
	intAudit "github.com/freekieb7/gopenehr/internal/audit"
	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/oauth"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/pkg/audit"
	"github.com/freekieb7/gopenehr/pkg/web/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Settings      *config.Settings
	Telemetry     *telemetry.Telemetry
	TenantService *Service
	OAuthService  *oauth.Service
	AuditSink     *intAudit.Sink
}

func NewHandler(settings *config.Settings, telemetry *telemetry.Telemetry, tenantService *Service, oauthService *oauth.Service, auditSink *intAudit.Sink) *Handler {
	return &Handler{
		Settings:      settings,
		Telemetry:     telemetry,
		OAuthService:  oauthService,
		TenantService: tenantService,
		AuditSink:     auditSink,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/tenants/v1")

	v1.Use(middleware.APIKeyProtected(h.Settings.APIKey))

	v1.Use(middleware.RequestID())
	v1.Use(middleware.Recover(h.Telemetry))
	v1.Use(middleware.Telemetry(h.Telemetry))

	var validateToken middleware.ValidateTokenFunc = nil
	if h.OAuthService.Enabled() {
		validateToken = h.OAuthService.ValidateToken
	}

	v1.Post("/", middleware.Audit(h.AuditSink.Enqueue, audit.ResourceTenant, audit.ActionCreate), middleware.JWTProtected([]string{oauth.ScopeWebhookManage.String()}, validateToken), h.CreateTenant)
}

type CreateTenantRequest struct {
	Name     string `json:"name"`
	EHRLimit int    `json:"ehr_limit"`
}

func (h *Handler) CreateTenant(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := middleware.AuditFrom(c)

	accepts := c.Accepts("application/json")
	if accepts == "" {
		return fiber.ErrNotAcceptable
	}

	var req CreateTenantRequest
	err := c.BodyParser(&req)
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Status:  "bad_request",
			Message: "invalid request body",
		})
	}

	tenant, err := h.TenantService.CreateTenant(ctx, req.Name, Subscription{
		EHRLimit: req.EHRLimit,
	})
	if err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to create tenant", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Status:  "internal_server_error",
			Message: "failed to create tenant",
		})
	}

	auditCtx.Event.Details["tenant_id"] = tenant.ID.String()
	auditCtx.Event.Details["tenant_name"] = tenant.Name
	auditCtx.Event.Details["ehr_limit"] = req.EHRLimit
	auditCtx.Success()

	return c.Status(fiber.StatusCreated).JSON(tenant)
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
