package webhook

import (
	"log/slog"
	"net/http"

	"github.com/freekieb7/gopenehr/internal/audit"
	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/pkg/utils"
	"github.com/freekieb7/gopenehr/pkg/web/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	Settings       *config.Settings
	Logger         *slog.Logger
	AuditService   *audit.Service
	WebhookService *Service
}

func NewHandler(settings *config.Settings, logger *slog.Logger, auditService *audit.Service, webhookService *Service) Handler {
	return Handler{
		Settings:       settings,
		Logger:         logger,
		AuditService:   auditService,
		WebhookService: webhookService,
	}
}

func (h *Handler) RegisterRoutes(a *fiber.App) {
	v1 := a.Group("/webhooks/v1")
	v1.Use(middleware.APIKeyProtected(h.Settings.APIKey))

	v1.Get("", h.HandleListSubscriptions)
	v1.Post("", h.HandleSubscribe)
	v1.Patch("/:subscription_id", h.HandleUpdateSubscription)
	v1.Delete("/:subscription_id", h.HandleUnsubscribe)
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

	var req SubscribeRequest
	if c.BodyParser(&req) != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "error",
		})
	}

	if req.URL == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Subscription URL is required",
			Status:  "error",
		})
	}
	if len(req.EventTypes) == 0 {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "At least one event type is required",
			Status:  "error",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceWebhook,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"url":         req.URL,
				"event_types": req.EventTypes,
				"outcome":     outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetEHRBySubject", "error", err)
		}
	}()

	exists, err := h.WebhookService.ExistsSubscriptionWithURL(ctx, req.URL)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to check existing subscription", "url", req.URL, "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create subscription",
			Status:  "error",
		})
	}
	if exists {
		outcome = "already_exists"
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
			outcome = "invalid_event_type"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid event type",
				Status:  "error",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to create subscription", "url", req.URL, "event_types", req.EventTypes, "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create subscription",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(http.StatusCreated).JSON(SubscribeResponse{
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

	outcome := "unknown"
	auditDetails := map[string]any{}
	defer func() {
		auditDetails["outcome"] = outcome
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceWebhook,
			Action:    audit.ActionUpdate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details:   auditDetails,
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for UpdateSubscription", "error", err)
		}
	}()

	subscriptionIDStr := c.Params("subscription_id")
	if subscriptionIDStr == "" {
		outcome = "invalid_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "subscription_id path parameter is required",
			Status:  "error",
		})
	}
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		outcome = "invalid_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "invalid subscription_id format",
			Status:  "error",
		})
	}

	auditDetails["subscription_id"] = subscriptionID

	// For simplicity, let's assume we only allow updating the event types
	var req UpdateSubscriptionRequest
	if c.BodyParser(&req) != nil {
		outcome = "invalid_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "error",
		})
	}

	auditDetails["event_types"] = req.EventTypes

	// Check if subscription exists
	exists, err := h.WebhookService.ExistsSubscription(ctx, subscriptionID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to check subscription existence", "subscription_id", subscriptionID, "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update subscription",
			Status:  "error",
		})
	}
	if !exists {
		outcome = "not_found"
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
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update subscription",
			Status:  "error",
		})
	}
	outcome = "success"

	return c.SendStatus(http.StatusNoContent)
}

func (h *Handler) HandleUnsubscribe(c *fiber.Ctx) error {
	ctx := c.Context()

	subscriptionIDStr := c.Params("subscription_id")
	if subscriptionIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "subscription_id query parameter is required",
			Status:  "error",
		})
	}
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "invalid subscription_id format",
			Status:  "error",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceWebhook,
			Action:    audit.ActionDelete,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"subscription_id": subscriptionID,
				"outcome":         outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Unsubscribe", "error", err)
		}
	}()

	err = h.WebhookService.Unsubscribe(ctx, subscriptionID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to unsubscribe", "subscription_id", subscriptionID, "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to unsubscribe",
			Status:  "error",
		})
	}

	outcome = "success"
	c.Status(http.StatusNoContent)
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
