package webhook

import (
	"github.com/freekieb7/gopenehr/internal/audit"
	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/oauth"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/pkg/utils"
	"github.com/freekieb7/gopenehr/pkg/web/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	Settings       *config.Settings
	Logger         *telemetry.Logger
	AuditSink      *audit.Sink
	OAuthService   *oauth.Service
	WebhookService *Service
}

func NewHandler(settings *config.Settings, logger *telemetry.Logger, auditSink *audit.Sink, oauthService *oauth.Service, webhookService *Service) Handler {
	return Handler{
		Settings:       settings,
		Logger:         logger,
		AuditSink:      auditSink,
		OAuthService:   oauthService,
		WebhookService: webhookService,
	}
}

func (h *Handler) RegisterRoutes(a *fiber.App) {
	v1 := a.Group("/webhooks/v1")
	v1.Use(middleware.APIKeyProtected(h.Settings.APIKey))
	v1.Use(oauth.JWTProtectedMiddleware(h.OAuthService, []oauth.Scope{oauth.ScopeWebhookManage}))

	v1.Get("", audit.AuditLoggedMiddleware(h.AuditSink, audit.ResourceWebhook, audit.ActionRead), h.HandleListSubscriptions)
	v1.Post("", audit.AuditLoggedMiddleware(h.AuditSink, audit.ResourceWebhook, audit.ActionCreate), h.HandleSubscribe)
	v1.Patch("/:subscription_id", audit.AuditLoggedMiddleware(h.AuditSink, audit.ResourceWebhook, audit.ActionUpdate), h.HandleUpdateSubscription)
	v1.Delete("/:subscription_id", audit.AuditLoggedMiddleware(h.AuditSink, audit.ResourceWebhook, audit.ActionDelete), h.HandleUnsubscribe)
}

type SubscribeRequest struct {
	URL        string   `json:"url"`
	EventTypes []string `json:"event_types"`
}

type UpdateSubscriptionRequest struct {
	EventTypes utils.Optional[[]string] `json:"event_types"`
}

type SubscribeResponse struct {
	SubscriptionID string   `json:"subscription_id"`
	URL            string   `json:"url"`
	Secret         string   `json:"secret"`
	EventTypes     []string `json:"event_types"`
	IsActive       bool     `json:"is_active"`
	CreatedAt      string   `json:"created_at"`
}

func (h *Handler) HandleListSubscriptions(c *fiber.Ctx) error {
	ctx := c.Context()

	subscriptions, err := h.WebhookService.ListSubscriptions(ctx)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to list subscriptions", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to list subscriptions",
			Status:  "error",
		})
	}

	resp := make([]SubscribeResponse, len(subscriptions))
	for i, sub := range subscriptions {
		var eventTypeStrs []string
		for _, et := range sub.EventTypes {
			eventTypeStrs = append(eventTypeStrs, string(et))
		}
		resp[i] = SubscribeResponse{
			SubscriptionID: sub.ID.String(),
			URL:            sub.URL,
			Secret:         sub.Secret,
			EventTypes:     eventTypeStrs,
			IsActive:       sub.IsActive,
			CreatedAt:      sub.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.JSON(resp)
}

func (h *Handler) HandleSubscribe(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	var req SubscribeRequest
	if c.BodyParser(&req) != nil {
		auditCtx.Fail("bad_request", "Invalid request body")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "error",
		})
	}

	if req.URL == "" {
		auditCtx.Fail("bad_request", "Subscription URL is required")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Subscription URL is required",
			Status:  "error",
		})
	}
	if len(req.EventTypes) == 0 {
		auditCtx.Fail("bad_request", "At least one event type is required")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "At least one event type is required",
			Status:  "error",
		})
	}

	exists, err := h.WebhookService.ExistsSubscriptionWithURL(ctx, req.URL)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to check existing subscription", "url", req.URL, "error", err)

		auditCtx.Fail("error", "Failed to check existing subscription")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create subscription",
			Status:  "error",
		})
	}
	if exists {
		auditCtx.Fail("already_exists", "Subscription with the given URL already exists")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusConflict,
			Message: "Subscription with the given URL already exists",
			Status:  "error",
		})
	}

	// Call the webhook service to create the subscription
	subscription, err := h.WebhookService.Subscribe(ctx, req.URL, req.EventTypes)
	if err != nil {
		if err == ErrInvalidEventType {
			auditCtx.Fail("invalid_event_type", "Invalid event type")
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid event type",
				Status:  "error",
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to create subscription", "url", req.URL, "event_types", req.EventTypes, "error", err)

		auditCtx.Fail("error", "Failed to create subscription")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create subscription",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusCreated).JSON(SubscribeResponse{
		SubscriptionID: subscription.ID.String(),
		URL:            subscription.URL,
		Secret:         subscription.Secret,
		EventTypes:     req.EventTypes,
		IsActive:       subscription.IsActive,
		CreatedAt:      subscription.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *Handler) HandleUpdateSubscription(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	subscriptionIDStr := c.Params("subscription_id")
	if subscriptionIDStr == "" {
		auditCtx.Fail("bad_request", "subscription_id path parameter is required")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "subscription_id path parameter is required",
			Status:  "error",
		})
	}
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		auditCtx.Fail("bad_request", "invalid subscription_id format")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "invalid subscription_id format",
			Status:  "error",
		})
	}

	// For simplicity, let's assume we only allow updating the event types
	var req UpdateSubscriptionRequest
	if c.BodyParser(&req) != nil {
		auditCtx.Fail("bad_request", "Invalid request body")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "error",
		})
	}

	// Check if subscription exists
	exists, err := h.WebhookService.ExistsSubscription(ctx, subscriptionID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to check subscription existence", "subscription_id", subscriptionID, "error", err)

		auditCtx.Fail("error", "Failed to check subscription existence")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update subscription",
			Status:  "error",
		})
	}
	if !exists {
		auditCtx.Fail("not_found", "Subscription not found")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotFound,
			Message: "Subscription not found",
			Status:  "error",
		})
	}

	// Call the webhook service to update the subscription
	err = h.WebhookService.UpdateSubscription(ctx, subscriptionID, req.EventTypes)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to update subscription", "subscription_id", subscriptionID, "error", err)

		auditCtx.Fail("error", "Failed to update subscription")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update subscription",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) HandleUnsubscribe(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	subscriptionIDStr := c.Params("subscription_id")
	if subscriptionIDStr == "" {
		auditCtx.Fail("bad_request", "subscription_id query parameter is required")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "subscription_id query parameter is required",
			Status:  "error",
		})
	}
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		auditCtx.Fail("bad_request", "invalid subscription_id format")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "invalid subscription_id format",
			Status:  "error",
		})
	}

	err = h.WebhookService.Unsubscribe(ctx, subscriptionID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to unsubscribe", "subscription_id", subscriptionID, "error", err)

		auditCtx.Fail("error", "Failed to unsubscribe")
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to unsubscribe",
			Status:  "error",
		})
	}

	auditCtx.Success()
	c.Status(fiber.StatusNoContent)
	return nil
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
	Status  string `json:"status"`
	Details any    `json:"details,omitempty"`
}

func SendErrorResponse(c *fiber.Ctx, errorRes ErrorResponse) error {
	return c.Status(errorRes.Code).JSON(map[string]ErrorResponse{
		"error": errorRes,
	})
}
