package openehr

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/audit"
	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/openehr/model"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/internal/webhook"
	"github.com/freekieb7/gopenehr/pkg/utils"
	"github.com/freekieb7/gopenehr/pkg/web/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	Settings       *config.Settings
	Logger         *slog.Logger
	OpenEHRService *Service
	AuditService   *audit.Service
	WebhookService *webhook.Service
}

func NewHandler(settings *config.Settings, logger *slog.Logger, openEHRService *Service, auditService *audit.Service, webhookService *webhook.Service) Handler {
	return Handler{
		Settings:       settings,
		Logger:         logger,
		OpenEHRService: openEHRService,
		AuditService:   auditService,
		WebhookService: webhookService,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/openehr/v1")
	v1.Use(middleware.APIKeyProtected(h.Settings.APIKey))

	v1.Options("", h.SystemInfo)

	v1.Get("/ehr", h.GetEHRBySubject)
	v1.Post("/ehr", h.CreateEHR)
	v1.Get("/ehr/:ehr_id", h.GetEHR)
	v1.Put("/ehr/:ehr_id", h.CreateEHRWithID)

	v1.Get("/ehr/:ehr_id/ehr_status", h.GetEHRStatus)
	v1.Put("/ehr/:ehr_id/ehr_status", h.UpdateEhrStatus)
	v1.Get("/ehr/:ehr_id/ehr_status/:version_uid", h.GetEHRStatusByVersionID)

	v1.Get("/ehr/:ehr_id/versioned_ehr_status", h.GetVersionedEHRStatus)
	v1.Get("/ehr/:ehr_id/versioned_ehr_status/revision_history", h.GetVersionedEHRStatusRevisionHistory)
	v1.Get("/ehr/:ehr_id/versioned_ehr_status/version", h.GetVersionedEHRStatusVersion)
	v1.Get("/ehr/:ehr_id/versioned_ehr_status/version/:version_uid", h.GetVersionedEHRStatusVersionByID)

	v1.Post("/ehr/:ehr_id/composition", h.CreateComposition)
	v1.Get("/ehr/:ehr_id/composition/:uid_based_id", h.GetComposition)
	v1.Put("/ehr/:ehr_id/composition/:uid_based_id", h.UpdateComposition)
	v1.Delete("/ehr/:ehr_id/composition/:uid_based_id", h.DeleteComposition)

	v1.Get("/ehr/:ehr_id/versioned_composition/:versioned_object_id", h.GetVersionedCompositionByID)
	v1.Get("/ehr/:ehr_id/versioned_composition/:versioned_object_id/revision_history", h.GetVersionedCompositionRevisionHistory)
	v1.Get("/ehr/:ehr_id/versioned_composition/:versioned_object_id/version", h.GetVersionedCompositionVersionAtTime)
	v1.Get("/ehr/:ehr_id/versioned_composition/:versioned_object_id/version/:version_uid", h.GetVersionedCompositionVersionByID)

	v1.Post("/ehr/:ehr_id/directory", h.CreateDirectory)
	v1.Put("/ehr/:ehr_id/directory", h.UpdateDirectory)
	v1.Delete("/ehr/:ehr_id/directory", h.DeleteDirectory)
	v1.Get("/ehr/:ehr_id/directory", h.GetFolderInDirectoryVersionAtTime)
	v1.Get("/ehr/:ehr_id/directory/:version_uid", h.GetFolderInDirectoryVersion)

	v1.Post("/ehr/:ehr_id/contribution", h.CreateContribution)
	v1.Get("/ehr/:ehr_id/contribution/:contribution_uid", h.GetContribution)

	v1.Get("/ehr/:ehr_id/tags", h.GetEHRTags)
	v1.Get("/ehr/:ehr_id/composition/:uid_based_id/tags", h.GetCompositionTags)
	v1.Put("/ehr/:ehr_id/composition/:uid_based_id/tags", h.UpdateCompositionTags)
	v1.Delete("/ehr/:ehr_id/composition/:uid_based_id/tags", h.DeleteCompositionTagByKey)
	v1.Get("/ehr/:ehr_id/ehr_status/tags", h.GetEHRStatusTags)
	v1.Get("/ehr/:ehr_id/ehr_status/:version_uid/tags", h.GetEHRStatusVersionTags)
	v1.Put("/ehr/:ehr_id/ehr_status/:version_uid/tags", h.UpdateEHRStatusVersionTags)
	v1.Delete("/ehr/:ehr_id/ehr_status/:version_uid/tags/:key", h.DeleteEHRStatusVersionTagByKey)

	v1.Post("/demographic/agent", h.CreateAgent)
	v1.Get("/demographic/agent/:uid_based_id", h.GetAgent)
	v1.Put("/demographic/agent/:uid_based_id", h.UpdateAgent)
	v1.Delete("/demographic/agent/:uid_based_id", h.DeleteAgent)

	v1.Post("/demographic/group", h.CreateGroup)
	v1.Get("/demographic/group/:uid_based_id", h.GetGroup)
	v1.Put("/demographic/group/:uid_based_id", h.UpdateGroup)
	v1.Delete("/demographic/group/:uid_based_id", h.DeleteGroup)

	v1.Post("/demographic/person", h.CreatePerson)
	v1.Get("/demographic/person/:uid_based_id", h.GetPerson)
	v1.Put("/demographic/person/:uid_based_id", h.UpdatePerson)
	v1.Delete("/demographic/person/:uid_based_id", h.DeletePerson)

	v1.Post("/demographic/organisation", h.CreateOrganisation)
	v1.Get("/demographic/organisation/:uid_based_id", h.GetOrganisation)
	v1.Put("/demographic/organisation/:uid_based_id", h.UpdateOrganisation)
	v1.Delete("/demographic/organisation/:uid_based_id", h.DeleteOrganisation)

	v1.Post("/demographic/role", h.CreateRole)
	v1.Get("/demographic/role/:uid_based_id", h.GetRole)
	v1.Put("/demographic/role/:uid_based_id", h.UpdateRole)
	v1.Delete("/demographic/role/:uid_based_id", h.DeleteRole)

	v1.Get("/demographic/versioned_party/:versioned_object_id", h.GetVersionedParty)
	v1.Get("/demographic/versioned_party/:versioned_object_id/revision_history", h.GetVersionedPartyRevisionHistory)
	v1.Get("/demographic/versioned_party/:versioned_object_id/version", h.GetVersionedPartyVersionAtTime)
	v1.Get("/demographic/versioned_party/:versioned_object_id/version/:version_id", h.GetVersionedPartyVersion)

	v1.Post("/demographic/contribution", h.CreateDemographicContribution)
	v1.Get("/demographic/contribution/:contribution_uid", h.GetDemographicContribution)

	v1.Get("/demographic/tags", h.GetDemographicTags)
	v1.Get("/demographic/agent/:uid_based_id/tags", h.GetAgentTags)
	v1.Put("/demographic/agent/:uid_based_id/tags", h.UpdateAgentTags)
	v1.Delete("/demographic/agent/:uid_based_id/tags/:key", h.DeleteAgentTagByKey)
	v1.Get("/demographic/group/:uid_based_id/tags", h.GetGroupTags)
	v1.Put("/demographic/group/:uid_based_id/tags", h.UpdateGroupTags)
	v1.Delete("/demographic/group/:uid_based_id/tags/:key", h.DeleteGroupTagByKey)
	v1.Get("/demographic/person/:uid_based_id/tags", h.GetPersonTags)
	v1.Put("/demographic/person/:uid_based_id/tags", h.UpdatePersonTags)
	v1.Delete("/demographic/person/:uid_based_id/tags/:key", h.DeletePersonTagByKey)
	v1.Get("/demographic/organisation/:uid_based_id/tags", h.GetOrganisationTags)
	v1.Put("/demographic/organisation/:uid_based_id/tags", h.UpdateOrganisationTags)
	v1.Delete("/demographic/organisation/:uid_based_id/tags/:key", h.DeleteOrganisationTagByKey)
	v1.Get("/demographic/role/:uid_based_id/tags", h.GetRoleTags)
	v1.Put("/demographic/role/:uid_based_id/tags", h.UpdateRoleTags)
	v1.Delete("/demographic/role/:uid_based_id/tags/:key", h.DeleteRoleTagByKey)

	v1.Get("/query/aql", h.ExecuteAdHocAQL)
	v1.Post("/query/aql", h.ExecuteAdHocAQLPost)
	v1.Get("/query/:qualified_query_name", h.ExecuteStoredAQL)
	v1.Post("/query/:qualified_query_name", h.ExecuteStoredAQLPost)
	v1.Get("/query/:qualified_query_name/:version", h.ExecuteStoredAQLVersion)
	v1.Post("/query/:qualified_query_name/:version", h.ExecuteStoredAQLVersionPost)

	v1.Get("/definition/template/adl1.4", h.GetTemplatesADL14)
	v1.Post("/definition/template/adl1.4", h.UploadTemplateADL14)
	v1.Get("/definition/template/adl1.4/:template_id", h.GetTemplateADL14ByID)

	v1.Get("/definition/template/adl2", h.GetTemplatesADL2)
	v1.Post("/definition/template/adl2", h.UploadTemplateADL2)
	v1.Get("/definition/template/adl2/:template_id", h.GetTemplateADL2ByID)
	v1.Get("/definition/template/adl2/:template_id/:version", h.GetTemplateADL2AtVersion)

	v1.Get("/definition/query/:qualified_query_name", h.ListStoredQueries)
	v1.Put("/definition/query/:qualified_query_name", h.StoreQuery)
	v1.Put("/definition/query/:qualified_query_name/:version", h.StoreQueryVersion)
	v1.Get("/definition/query/:qualified_query_name/:version", h.GetStoredQueryAtVersion)

	v1.Delete("/admin/ehr/all", h.DeleteMultipleEHRs)
	v1.Delete("/admin/ehr/:ehr_id", h.DeleteEHRByID)
}

func (h *Handler) SystemInfo(c *fiber.Ctx) error {
	response := map[string]any{
		"solution":              "gopenEHR",
		"version":               h.Settings.Version,
		"vendor":                "freekieb7",
		"restapi_specs_version": "development",
		"conformance_profile":   "STANDARD",
		"endpoints": []string{
			"/ehr",
			"/demographic",
			"/definition",
			"/query",
			"/admin",
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) GetEHRBySubject(c *fiber.Ctx) error {
	ctx := c.Context()

	subjectID := c.Query("subject_id")
	if subjectID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "subject_id query parameters are required",
			Status:  "bad_request",
		})
	}

	subjectNamespace := c.Query("subject_namespace")
	if subjectNamespace == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "subject_namespace query parameters are required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHR,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"subject_id":        subjectID,
				"subject_namespace": subjectNamespace,
				"outcome":           outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetEHRBySubject", "error", err)
		}
	}()

	ehr, err := h.OpenEHRService.GetEHRBySubject(ctx, subjectID, subjectNamespace)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given subject",
				Status:  "not_found",
			})
		}

		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get EHR by subject",
			Status:  "error",
		})
	}
	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(ehr)
}

func (h *Handler) CreateEHR(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	// response type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Check for optional EHR_STATUS in the request body
	ehrID := uuid.New()
	ehrStatus := NewEHRStatus(uuid.New())
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&ehrStatus); err != nil {
			if err, ok := err.(util.ValidateError); ok {
				return SendErrorResponse(c, ErrorResponse{
					Code:    fiber.StatusBadRequest,
					Message: "Validation error in request body",
					Status:  "bad_request",
					Details: err,
				})
			}
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid request body",
				Status:  "bad_request",
			})
		}
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHR,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":  ehrID.String(),
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for CreateEHR", "error", err)
		}
	}()

	// Create EHR
	ehr, err := h.OpenEHRService.CreateEHR(ctx, ehrID, ehrStatus)
	if err != nil {
		if err == ErrEHRStatusAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "EHR Status with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to create EHR", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeEHRCreated, map[string]any{"ehr": ehr})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for EHR creation", "error", err)
	}

	// Determine response
	c.Set("ETag", "\""+ehrID.String()+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String())

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(ehr)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + ehrID.String() + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(ehr)
	}
}

func (h *Handler) GetEHR(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHR,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":  ehrID,
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetEHR", "error", err)
		}
	}()

	ehr, err := h.OpenEHRService.GetEHR(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given ID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get EHR by ID", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(ehr)
}

func (h *Handler) CreateEHRWithID(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id format",
			Status:  "bad_request",
		})
	}

	// response type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Check for optional EHR_STATUS in the request body
	ehrStatus := NewEHRStatus(uuid.New())
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&ehrStatus); err != nil {
			if err, ok := err.(util.ValidateError); ok {
				return SendErrorResponse(c, ErrorResponse{
					Code:    fiber.StatusBadRequest,
					Message: "Validation error in request body",
					Status:  "bad_request",
					Details: err,
				})
			}

			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid request body",
				Status:  "bad_request",
			})
		}
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHR,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":  ehrID.String(),
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for CreateEHRWithID", "error", err)
		}
	}()

	// Create EHR with specified ID and EHR_STATUS
	ehr, err := h.OpenEHRService.CreateEHR(ctx, ehrID, ehrStatus)
	if err != nil {
		if err == ErrEHRAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "EHR with the given ID already exists",
				Status:  "conflict",
			})
		}
		if err == ErrEHRStatusAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "EHR Status with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to create EHR", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeEHRCreated, map[string]any{"ehr": ehr})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for EHR creation", "error", err)
	}

	// Determine response
	c.Set("ETag", "\""+ehrID.String()+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String())

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(ehr)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + ehrID.String() + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(ehr)
	}
}

func (h *Handler) GetEHRStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id format",
			Status:  "bad_request",
		})
	}

	// Optional filter: version_at_time
	var filterAtTime time.Time
	if atTimeStr := c.Query("version_at_time"); atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid version_at_time format. Use RFC3339 format.",
				Status:  "bad_request",
			})
		}
		filterAtTime = parsedTime
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHRStatus,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":          ehrID,
				"version_at_time": filterAtTime,
				"outcome":         outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetEHRStatus", "error", err)
		}
	}()

	ehrStatus, err := h.OpenEHRService.GetEHRStatus(ctx, ehrID, filterAtTime, "")
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR Status not found for the given EHR ID at the specified time",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get EHR Status at time", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(ehrStatus)
}

func (h *Handler) UpdateEhrStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id format",
			Status:  "bad_request",
		})
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "If-Match header is required",
			Status:  "bad_request",
		})
	}

	// response type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	var requestEhrStatus model.EHR_STATUS
	if err := c.BodyParser(&requestEhrStatus); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHRStatus,
			Action:    audit.ActionUpdate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":  ehrID,
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for UpdateEhrStatus", "error", err)
		}
	}()

	// Check collision using If-Match header
	currentEHRStatus, err := h.OpenEHRService.GetEHRStatus(ctx, ehrID, time.Time{}, "")
	if err != nil {
		if err == ErrEHRStatusNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR Status not found for the given EHR ID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get current EHR Status", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	currentEHRStatusID, ok := currentEHRStatus.UID.V.Value.(*model.OBJECT_VERSION_ID)
	if !ok {
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Current EHR Status UID is not of type OBJECT_VERSION_ID",
			Status:  "error",
		})
	}

	// Check collision using If-Match header
	if currentEHRStatusID.Value != ifMatch {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "EHR Status has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Proceed to update EHR Status
	updatedEHRStatus, err := h.OpenEHRService.UpdateEHRStatus(ctx, ehrID, requestEhrStatus)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR Status not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrEHRStatusAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "EHR Status with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err == ErrEHRStatusVersionLowerOrEqualToCurrent {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "EHR Status version in request body must be incremented",
				Status:  "bad_request",
			})
		}
		if err == ErrInvalidEHRStatusUIDMismatch {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "EHR Status UID HIER_OBJECT_ID in request body does not match current EHR Status UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to update EHR Status", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeEHRStatusUpdated, map[string]any{"ehr_status": updatedEHRStatus})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for EHR Status updated", "error", err)
	}

	// Determine response
	updatedEHRStatusID := updatedEHRStatus.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedEHRStatusID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String()+"/ehr_status/"+updatedEHRStatusID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedEHRStatus)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusOK).JSON(`{"uid":"` + updatedEHRStatusID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedEHRStatus)
	}
}

func (h *Handler) GetEHRStatusByVersionID(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id format",
			Status:  "bad_request",
		})
	}

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version_uid parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHRStatus,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":      ehrID,
				"version_uid": versionUID,
				"outcome":     outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetEHRStatusByVersionID", "error", err)
		}
	}()

	ehrStatus, err := h.OpenEHRService.GetEHRStatus(ctx, ehrID, time.Time{}, versionUID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR Status not found for the given EHR ID and version UID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get EHR Status at time", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(ehrStatus)
}

func (h *Handler) GetVersionedEHRStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHRStatus,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":  ehrID,
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedEHRStatus", "error", err)
		}
	}()

	versionedStatus, err := h.OpenEHRService.GetVersionedEHRStatus(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned EHR Status not found for the given EHR ID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(versionedStatus)
}

func (h *Handler) GetVersionedEHRStatusRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceVersionedEHRStatus,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":  ehrID,
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedEHRStatusRevisionHistory", "error", err)
		}
	}()

	revisionHistory, err := h.OpenEHRService.GetVersionedEHRStatusRevisionHistory(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned EHR Status revision history not found for the given EHR ID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status revision history", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(revisionHistory)
}

func (h *Handler) GetVersionedEHRStatusVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id parameter format",
			Status:  "bad_request",
		})
	}

	var filterAtTime time.Time
	atTimeStr := c.Query("version_at_time")
	if atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid version_at_time format. Use RFC3339 format.",
				Status:  "bad_request",
			})
		}
		filterAtTime = parsedTime
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHRStatus,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":          ehrID.String(),
				"version_at_time": filterAtTime,
				"outcome":         outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedEHRStatusVersion", "error", err)
		}
	}()

	ehrStatusVersionJSON, err := h.OpenEHRService.GetVersionedEHRStatusVersionAsJSON(ctx, ehrID, filterAtTime, "")
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned EHR Status version not found for the given EHR ID at the specified time",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status version at time", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(ehrStatusVersionJSON)
}

func (h *Handler) GetVersionedEHRStatusVersionByID(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id parameter format",
			Status:  "bad_request",
		})
	}

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version_uid parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceEHRStatus,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":      ehrID.String(),
				"version_uid": versionUID,
				"outcome":     outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedEHRStatusVersionByID", "error", err)
		}
	}()

	ehrStatusVersionJSON, err := h.OpenEHRService.GetVersionedEHRStatusVersionAsJSON(ctx, ehrID, time.Time{}, versionUID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned EHR Status version not found for the given EHR ID and version UID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status version by ID", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(ehrStatusVersionJSON)
}

func (h *Handler) CreateComposition(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id parameter format",
			Status:  "bad_request",
		})
	}

	// response type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Parse new composition from request body
	var requestComposition model.COMPOSITION
	if err := c.BodyParser(&requestComposition); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceComposition,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":  ehrID.String(),
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for CreateComposition", "error", err)
		}
	}()

	// Check if EHR exists
	exists, err := h.OpenEHRService.ExistsEHR(ctx, ehrID)
	if err != nil {
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to check if EHR exists", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}
	if !exists {
		outcome = "not_found"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotFound,
			Message: "EHR not found for the given EHR ID",
			Status:  "not_found",
		})
	}

	// Create Composition
	newComposition, err := h.OpenEHRService.CreateComposition(ctx, ehrID, requestComposition)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Composition with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to create Composition", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeCompositionCreated, map[string]any{"composition": newComposition})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Composition created", "error", err)
	}

	// Determine response
	compositionID := newComposition.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+compositionID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String()+"/composition/"+compositionID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(newComposition)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + compositionID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(newComposition)
	}
}

func (h *Handler) GetComposition(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id parameter format",
			Status:  "bad_request",
		})
	}

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceComposition,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":       ehrID.String(),
				"uid_based_id": uidBasedID,
				"outcome":      outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetComposition", "error", err)
		}
	}()

	composition, err := h.OpenEHRService.GetComposition(ctx, ehrID, uidBasedID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Composition not found for the given UID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Composition by ID", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(composition)
}

func (h *Handler) UpdateComposition(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id parameter format",
			Status:  "bad_request",
		})
	}

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "If-Match header is required",
			Status:  "bad_request",
		})
	}

	// response type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceComposition,
			Action:    audit.ActionUpdate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":       ehrID.String(),
				"uid_based_id": uidBasedID,
				"outcome":      outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for UpdateComposition", "error", err)
		}
	}()

	// Check collision using If-Match header
	currentComposition, err := h.OpenEHRService.GetComposition(ctx, ehrID, strings.Split(uidBasedID, "::")[0])
	if err != nil {
		if err == ErrCompositionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Composition not found for the given EHR ID and UID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get current Composition", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Composition
	if currentComposition.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != ifMatch {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Composition has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Parse updated composition from request body
	var requestComposition model.COMPOSITION
	if err := c.BodyParser(&requestComposition); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			outcome = "validation_error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	// Proceed to update Composition
	composition, err := h.OpenEHRService.UpdateComposition(ctx, ehrID, requestComposition)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Composition not found for the given UID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Composition with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err == ErrCompositionVersionLowerOrEqualToCurrent {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Composition version in request body is lower or equal to the current version",
				Status:  "bad_request",
			})
		}
		if err == ErrInvalidCompositionUIDMismatch {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Composition UID HIER_OBJECT_ID in request body does not match current Composition UID",
				Status:  "bad_request",
			})
		}
		if err == ErrCompositionUIDNotProvided {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Composition UID must be provided in the request body for update",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "validation_error"
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to update Composition", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeCompositionUpdated, map[string]any{"composition": composition})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Composition updated", "error", err)
	}

	// Determine response
	compositionID := composition.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+compositionID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String()+"/composition/"+compositionID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(composition)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + compositionID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(composition)
	}
}

func (h *Handler) DeleteComposition(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}

	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id parameter format",
			Status:  "bad_request",
		})
	}

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	if strings.Count(uidBasedID, "::") != 2 {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Composition by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceComposition,
			Action:    audit.ActionDelete,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":       ehrID.String(),
				"uid_based_id": uidBasedID,
				"outcome":      outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for DeleteComposition", "error", err)
		}
	}()

	// Check existence before deletion
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentComposition, err := h.OpenEHRService.GetComposition(ctx, ehrID, versionedObjectID)
	if err != nil {
		if err == ErrCompositionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Composition not found for the given UID and EHR ID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Composition by ID before deletion", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Composition
	if currentComposition.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != uidBasedID {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Composition has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Proceed to delete Composition
	if err := h.OpenEHRService.DeleteComposition(ctx, ehrID, uuid.MustParse(currentComposition.UID.V.Value.(*model.OBJECT_VERSION_ID).Value)); err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Composition not found for the given UID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to delete Composition by ID", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeCompositionDeleted, map[string]any{"composition": currentComposition})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Composition deleted", "error", err)
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) GetVersionedCompositionByID(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "versioned_object_id parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceVersionedComposition,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":              ehrID,
				"versioned_object_id": versionedObjectID,
				"outcome":             outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedCompositionByID", "error", err)
		}
	}()

	versionedComposition, err := h.OpenEHRService.GetVersionedComposition(ctx, ehrID, versionedObjectID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Composition not found for the given versioned object ID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Composition by ID", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(versionedComposition)
}

func (h *Handler) GetVersionedCompositionRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "versioned_object_id parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceVersionedComposition,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":              ehrID,
				"versioned_object_id": versionedObjectID,
				"outcome":             outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedCompositionRevisionHistory", "error", err)
		}
	}()

	revisionHistory, err := h.OpenEHRService.GetVersionedCompositionRevisionHistory(ctx, ehrID, versionedObjectID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Composition revision history not found for the given versioned object ID",
				Status:  "not_found",
			})
		}

		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Composition revision history", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(revisionHistory)
}

func (h *Handler) GetVersionedCompositionVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "versioned_object_id parameter is required",
			Status:  "bad_request",
		})
	}

	var filterAtTime time.Time
	atTimeStr := c.Query("version_at_time")
	if atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid version_at_time format. Use RFC3339 format.",
				Status:  "bad_request",
			})
		}
		filterAtTime = parsedTime
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceVersionedComposition,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":              ehrID,
				"versioned_object_id": versionedObjectID,
				"version_at_time":     filterAtTime,
				"outcome":             outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedCompositionVersionAtTime", "error", err)
		}
	}()

	versionJSON, err := h.OpenEHRService.GetVersionedCompositionVersionJSON(ctx, ehrID, versionedObjectID, filterAtTime, "")
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Composition version not found for the given versioned object ID at the specified time",
				Status:  "not_found",
			})
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Composition version at time", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	outcome = "success"
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(versionJSON)
}

func (h *Handler) GetVersionedCompositionVersionByID(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "versioned_object_id parameter is required",
			Status:  "bad_request",
		})
	}

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version_uid parameter is required",
			Status:  "bad_request",
		})
	}

	if versionedObjectID != strings.Split(versionUID, "::")[0] {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version_uid does not match the versioned_object_id",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceVersionedComposition,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":              ehrID,
				"versioned_object_id": versionedObjectID,
				"version_uid":         versionUID,
				"outcome":             outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedCompositionVersionByID", "error", err)
		}
	}()

	versionJSON, err := h.OpenEHRService.GetVersionedCompositionVersionJSON(ctx, ehrID, versionedObjectID, time.Time{}, versionUID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Composition version not found for the given versioned object ID and version UID",
				Status:  "not_found",
			})
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Composition version by ID", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(versionJSON)
}

func (h *Handler) CreateDirectory(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id parameter format",
			Status:  "bad_request",
		})
	}

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	if len(c.Body()) <= 0 {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Request body is required",
			Status:  "bad_request",
		})
	}

	var requestDirectory model.FOLDER
	if err := c.BodyParser(&requestDirectory); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceDirectory,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":  ehrID.String(),
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for CreateDirectory", "error", err)
		}
	}()

	directory, err := h.OpenEHRService.CreateDirectory(ctx, ehrID, requestDirectory)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Directory with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to create Directory", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeDirectoryCreated, map[string]any{"directory": directory})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Directory created", "error", err)
	}

	// Determine response
	directoryID := directory.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+directoryID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String()+"/directory/"+directoryID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(directory)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + directoryID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(directory)
	}
}

func (h *Handler) UpdateDirectory(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id parameter format",
			Status:  "bad_request",
		})
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "If-Match header is required",
			Status:  "bad_request",
		})
	}

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceDirectory,
			Action:    audit.ActionUpdate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":  ehrID.String(),
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for UpdateDirectory", "error", err)
		}
	}()

	// Check collision using If-Match header
	currentDirectory, err := h.OpenEHRService.GetDirectory(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get current Directory", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Directory",
			Status:  "error",
		})
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Directory
	if currentDirectory.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != ifMatch {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Directory has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Parse updated directory from request body
	var requestDirectory model.FOLDER
	if err := c.BodyParser(&requestDirectory); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			outcome = "validation_error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}
		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	directory, err := h.OpenEHRService.UpdateDirectory(ctx, ehrID, requestDirectory)
	if err != nil {
		if err == ErrDirectoryNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Directory with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err == ErrDirectoryVersionLowerOrEqualToCurrent {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Directory version in request body is lower or equal to the current version",
				Status:  "bad_request",
			})
		}
		if err == ErrInvalidDirectoryUIDMismatch {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Directory UID HIER_OBJECT_ID in request body does not match current Directory UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "validation_error"
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to update Directory", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Directory",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeDirectoryUpdated, map[string]any{"directory": directory})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Directory updated", "error", err)
	}

	// Determine response
	updatedDirectoryID := directory.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedDirectoryID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String()+"/directory/"+updatedDirectoryID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(directory)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusOK).JSON(`{"uid":"` + updatedDirectoryID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(directory)
	}
}

func (h *Handler) DeleteDirectory(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id parameter format",
			Status:  "bad_request",
		})
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "If-Match header is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceDirectory,
			Action:    audit.ActionDelete,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":  ehrID.String(),
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for DeleteDirectory", "error", err)
		}
	}()

	// Check collision using If-Match header
	currentDirectory, err := h.OpenEHRService.GetDirectory(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get current Directory", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Directory",
			Status:  "error",
		})
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Directory
	if currentDirectory.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != ifMatch {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Directory has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	if err := h.OpenEHRService.DeleteDirectory(ctx, ehrID, uuid.MustParse(currentDirectory.UID.V.Value.(*model.OBJECT_VERSION_ID).Value)); err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to delete Directory", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Directory",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeDirectoryDeleted, map[string]any{"directory": currentDirectory})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Directory deleted", "error", err)
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) GetFolderInDirectoryVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}

	var filterAtTime time.Time
	atTimeStr := c.Query("version_at_time")
	if atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid version_at_time format. Use RFC3339 format.",
				Status:  "bad_request",
			})
		}
		filterAtTime = parsedTime
	}

	path := c.Query("path")
	var pathParts []string
	if path != "" {
		for part := range strings.SplitSeq(path, "/") {
			pathParts = append(pathParts, strings.TrimSpace(part))
		}
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceFolder,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":          ehrID,
				"version_at_time": filterAtTime,
				"path":            path,
				"outcome":         outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetFolderInDirectoryVersionAtTime", "error", err)
		}
	}()

	folder, err := h.OpenEHRService.GetFolderInDirectoryVersion(ctx, ehrID, filterAtTime, "", pathParts)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID at the specified time",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory version not found at the specified time for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrFolderNotFoundInDirectory {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Folder not found in Directory version at the specified time for the given path",
				Status:  "not_found",
			})
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Folder in Directory version at time", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Folder in Directory version at time",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(folder)
}

func (h *Handler) GetFolderInDirectoryVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version_uid parameter is required",
			Status:  "bad_request",
		})
	}

	path := c.Query("path")
	var pathParts []string
	if path != "" {
		for part := range strings.SplitSeq(path, "/") {
			pathParts = append(pathParts, strings.TrimSpace(part))
		}
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceFolder,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":      ehrID,
				"version_uid": versionUID,
				"path":        path,
				"outcome":     outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetFolderInDirectoryVersion", "error", err)
		}
	}()

	folder, err := h.OpenEHRService.GetFolderInDirectoryVersion(ctx, ehrID, time.Time{}, versionUID, pathParts)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory version not found at the specified time for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrFolderNotFoundInDirectory {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Folder not found in Directory version for the given path",
				Status:  "not_found",
			})
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Folder in Directory version", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Folder in Directory version",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(folder)
}

func (h *Handler) CreateContribution(c *fiber.Ctx) error {
	c.Status(fiber.StatusNotImplemented)
	return nil

	// ctx := c.Context()

	// c.Accepts("application/json")

	// ehrIDStr := c.Params("ehr_id")
	// if ehrIDStr == "" {
	// 	return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	// }
	// ehrID, err := uuid.Parse(ehrIDStr)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	// }

	// // return type
	// returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	// if !returnType.IsValid() {
	// 	return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	// }

	// // Parse new contribution from request body
	// var newContribution openehr.CONTRIBUTION
	// if err := c.BodyParser(&newContribution); err != nil {
	// 	if err, ok := err.(util.ValidateError); ok {
	// 		return c.Status(fiber.StatusBadRequest).JSON(err)
	// 	}

	// 	return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	// }

	// contribution, err := h.OpenEHRService.CreateContribution(ctx, ehrID, newContribution)
	// if err != nil {
	// 	if err == ErrContributionAlreadyExists {
	// 		return c.Status(fiber.StatusConflict).SendString("Contribution with the given UID already exists")
	// 	}

	// 	h.Logger.ErrorContext(ctx, "Failed to create Contribution", "error", err)
	// 	c.Status(http.StatusInternalServerError)
	// 	return nil
	// }

	// // Determine response
	// contributionID := contribution.UID.Value
	// c.Set("ETag", "\""+contributionID+"\"")
	// c.Set("Location", c.Protocol()+ "://" + c.Hostname() + "/openehr/v1/ehr/"+ehrID.String()+"/contribution/"+contributionID)

	// c.Status(fiber.StatusCreated)
	// switch returnType {
	// case ReturnTypeMinimal:
	// 	return nil
	// case ReturnTypeRepresentation:
	// 	return c.JSON(contribution)
	// case ReturnTypeIdentifier:
	// 	return c.JSON(`{"uid":"` + contributionID + `"}`)
	// default:
	// 	h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
	// 	return c.JSON(contribution)
	// }
}

func (h *Handler) GetContribution(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ehr_id parameter format",
			Status:  "bad_request",
		})
	}

	contributionUID := c.Params("contribution_uid")
	if contributionUID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "contribution_uid parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceContribution,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"ehr_id":           ehrID,
				"contribution_uid": contributionUID,
				"outcome":          outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetContribution", "error", err)
		}
	}()

	contribution, err := h.OpenEHRService.GetContribution(ctx, contributionUID, utils.Some(ehrID))
	if err != nil {
		if err == ErrContributionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Contribution not found for the given EHR ID and Contribution UID",
				Status:  "not_found",
			})
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to get Contribution by ID", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Contribution by ID",
			Status:  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(contribution)
}

func (h *Handler) GetEHRTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get EHR Tags not implemented yet")
}

func (h *Handler) GetCompositionTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Composition Tags not implemented yet")
}

func (h *Handler) UpdateCompositionTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Update Composition Tags not implemented yet")
}

func (h *Handler) DeleteCompositionTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Delete Composition Tag By Key not implemented yet")
}

func (h *Handler) GetEHRStatusTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get EHR Status Tags not implemented yet")
}

func (h *Handler) GetEHRStatusVersionTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get EHR Status Version Tags not implemented yet")
}

func (h *Handler) UpdateEHRStatusVersionTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Update EHR Status Version Tags not implemented yet")
}

func (h *Handler) DeleteEHRStatusVersionTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Delete EHR Status Version Tag By Key not implemented yet")
}

func (h *Handler) CreateAgent(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Parse new agent from request body
	var newAgent model.AGENT
	if err := c.BodyParser(&newAgent); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Agent", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceAgent,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"agent_uid": newAgent.UID,
				"outcome":   outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for CreateAgent", "error", err)
		}
	}()

	agent, err := h.OpenEHRService.CreateAgent(ctx, newAgent)
	if err != nil {
		if err == ErrAgentAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Agent with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}
		outcome = "error"
		h.Logger.ErrorContext(ctx, "Failed to create Agent", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create Agent",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeAgentCreated, map[string]any{"agent": agent})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Agent created", "error", err)
	}

	// Determine response
	agentID := agent.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+agentID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/agent/"+agentID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(agent)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + agentID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(agent)
	}
}

func (h *Handler) GetAgent(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceAgent,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"uid_based_id": uidBasedID,
				"outcome":      outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetAgent", "error", err)
		}
	}()

	if strings.Count(uidBasedID, "::") == 3 {
		// Is version id
		agent, err := h.OpenEHRService.GetAgentAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == ErrAgentNotFound {
				outcome = "not_found"
				return SendErrorResponse(c, ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Agent not found for the given agent ID",
					Status:  "not_found",
				})
			}

			h.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)
			outcome = "error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to get Agent by ID",
				Status:  "error",
			})
		}

		outcome = "success"
		return c.Status(fiber.StatusOK).JSON(agent)
	}

	// Is versioned party id
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	agent, err := h.OpenEHRService.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrAgentNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Agent not found for the given agent ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Agent by ID",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(agent)
}

func (h *Handler) UpdateAgent(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionID := c.Params("uid_based_id")
	if versionID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "If-Match header is required",
			Status:  "bad_request",
		})
	}

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Ensure Agent exists before update
	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceAgent,
			Action:    audit.ActionUpdate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_party_id": versionedPartyID,
				"outcome":            outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for UpdateAgent", "error", err)
		}
	}()

	// Check collision using If-Match header
	currentAgent, err := h.OpenEHRService.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrAgentNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Agent not found for the given agent ID",
				Status:  "not_found",
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to get current Agent by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Agent by ID",
			Status:  "error",
		})
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Agent
	if currentAgent.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != ifMatch {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Agent has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Parse updated agent from request body
	var requestAgent model.AGENT
	if err := c.BodyParser(&requestAgent); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Agent", "error", err)
		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	// Proceed to update Agent
	updatedAgent, err := h.OpenEHRService.UpdateAgent(ctx, versionedPartyID, requestAgent)
	if err != nil {
		if err == ErrAgentNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Agent not found for the given agent ID",
				Status:  "not_found",
			})
		}
		if err == ErrAgentVersionLowerOrEqualToCurrent {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Agent version is lower than or equal to the current version",
				Status:  "conflict",
			})
		}
		if err == ErrInvalidAgentUIDMismatch {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Agent UID HIER_OBJECT_ID in request body does not match current Agent UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "validation_error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to update Agent", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Agent",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeAgentUpdated, map[string]any{"agent": updatedAgent})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Agent updated", "error", err)
	}

	// Determine response
	updatedAgentID := updatedAgent.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedAgentID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/agent/"+updatedAgentID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedAgent)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusOK).JSON(`{"uid":"` + updatedAgentID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedAgent)
	}
}

func (h *Handler) DeleteAgent(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	if strings.Count(uidBasedID, "::") != 2 {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Agent by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	// Check existence before deletion
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceAgent,
			Action:    audit.ActionDelete,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_party_id": versionedPartyID,
				"outcome":            outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for DeleteAgent", "error", err)
		}
	}()

	currentAgent, err := h.OpenEHRService.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrAgentNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Agent not found for the given agent ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Agent by ID before deletion", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Agent by ID before deletion",
			Status:  "error",
		})
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Agent
	if currentAgent.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != uidBasedID {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Agent has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Proceed to delete Agent
	if err := h.OpenEHRService.DeleteAgent(ctx, uidBasedID); err != nil {
		if err == ErrAgentNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Agent not found for the given agent ID",
				Status:  "not_found",
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Agent by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Agent by ID",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeAgentDeleted, map[string]any{"agent": currentAgent})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Agent deleted", "error", err)
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) CreateGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Parse new group from request body
	var newGroup model.GROUP
	if err := c.BodyParser(&newGroup); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Group", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceGroup,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"group_uid": newGroup.UID,
				"outcome":   outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for CreateGroup", "error", err)
		}
	}()

	group, err := h.OpenEHRService.CreateGroup(ctx, newGroup)
	if err != nil {
		if err == ErrGroupAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Group with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to create Group", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create Group",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeGroupCreated, map[string]any{"group": group})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Group created", "error", err)
	}

	// Determine response
	groupID := group.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+groupID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/group/"+groupID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(group)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + groupID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(group)
	}
}

func (h *Handler) GetGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceGroup,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"uid_based_id": uidBasedID,
				"outcome":      outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetGroup", "error", err)
		}
	}()

	if strings.Count(uidBasedID, "::") == 3 {
		// Is version id
		group, err := h.OpenEHRService.GetGroupAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == ErrGroupNotFound {
				outcome = "not_found"
				return SendErrorResponse(c, ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Group not found for the given group ID",
					Status:  "not_found",
				})
			}
			h.Logger.ErrorContext(ctx, "Failed to get Group by ID", "error", err)
			outcome = "error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to get Group by ID",
				Status:  "error",
			})
		}

		outcome = "success"
		return c.Status(fiber.StatusOK).JSON(group)

	}
	// Is versioned party id
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	group, err := h.OpenEHRService.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrGroupNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Group not found for the given group ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Group by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Group by ID",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(group)
}

func (h *Handler) UpdateGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionID := c.Params("uid_based_id")
	if versionID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "If-Match header is required",
			Status:  "bad_request",
		})
	}

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Ensure Agent exists before update
	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceGroup,
			Action:    audit.ActionUpdate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_party_id": versionedPartyID,
				"outcome":            outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for UpdateGroup", "error", err)
		}
	}()

	// Check collision using If-Match header
	currentGroup, err := h.OpenEHRService.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrGroupNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Group not found for the given group ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Group by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Group by ID",
			Status:  "error",
		})
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Group
	if currentGroup.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != ifMatch {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Group has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Parse updated group from request body
	var requestGroup model.GROUP
	if err := c.BodyParser(&requestGroup); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Group", "error", err)
		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	// Proceed to update Group
	updatedGroup, err := h.OpenEHRService.UpdateGroup(ctx, versionedPartyID, requestGroup)
	if err != nil {
		if err == ErrGroupNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Group not found for the given group ID",
				Status:  "not_found",
			})
		}
		if err == ErrGroupVersionLowerOrEqualToCurrent {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Group version is lower than or equal to the current version",
				Status:  "conflict",
			})
		}
		if err == ErrInvalidGroupUIDMismatch {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Group UID HIER_OBJECT_ID in request body does not match current Group UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "validation_error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to update Agent", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Agent",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeGroupUpdated, map[string]any{"group": updatedGroup})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Group updated", "error", err)
	}

	// Determine response
	updatedGroupID := updatedGroup.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedGroupID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/group/"+updatedGroupID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedGroup)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusOK).JSON(`{"uid":"` + updatedGroupID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedGroup)
	}
}

func (h *Handler) DeleteGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	if !strings.Contains(uidBasedID, "::") {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Group by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	// Check existence before deletion
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceGroup,
			Action:    audit.ActionDelete,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_party_id": versionedPartyID,
				"outcome":            outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for DeleteGroup", "error", err)
		}
	}()

	currentGroup, err := h.OpenEHRService.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrGroupNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Group not found for the given group ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Group by ID before deletion", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Group by ID before deletion",
			Status:  "error",
		})
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Group
	if currentGroup.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != uidBasedID {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Group has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Proceed to delete Group
	if err := h.OpenEHRService.DeleteGroup(ctx, uidBasedID); err != nil {
		if err == ErrGroupNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Group not found for the given group ID",
				Status:  "not_found",
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Group by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Group by ID",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeGroupDeleted, map[string]any{"group": currentGroup})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Group deleted", "error", err)
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) CreatePerson(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Parse new person from request body
	var newPerson model.PERSON
	if err := c.BodyParser(&newPerson); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Person", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourcePerson,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"person_uid": newPerson.UID,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for CreatePerson", "error", err)
		}
	}()

	person, err := h.OpenEHRService.CreatePerson(ctx, newPerson)
	if err != nil {
		if err == ErrPersonAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Person with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to create Person", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create Person",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypePersonCreated, map[string]any{"person": person})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Person created", "error", err)
	}

	// Determine response
	personID := person.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+personID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/person/"+personID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(person)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + personID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(person)
	}
}

func (h *Handler) GetPerson(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourcePerson,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"person_uid": uidBasedID,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetPerson", "error", err)
		}
	}()

	if strings.Count(uidBasedID, "::") == 3 {
		// Is version id
		person, err := h.OpenEHRService.GetPersonAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == ErrPersonNotFound {
				outcome = "not_found"
				return SendErrorResponse(c, ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Person not found for the given person ID",
					Status:  "not_found",
				})
			}
			h.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)
			outcome = "error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to get Person",
				Status:  "error",
			})
		}

		outcome = "success"
		return c.Status(fiber.StatusOK).JSON(person)
	}

	// Is versioned party id
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	person, err := h.OpenEHRService.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrPersonNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Person not found for the given person ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Person",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(person)
}

func (h *Handler) UpdatePerson(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionID := c.Params("uid_based_id")
	if versionID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "If-Match header is required",
			Status:  "bad_request",
		})
	}

	// Parse return type from Prefer header
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Ensure Person exists before update
	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourcePerson,
			Action:    audit.ActionUpdate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_party_id": versionedPartyID,
				"outcome":            outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for UpdatePerson", "error", err)
		}
	}()

	// Check collision using If-Match header
	currentPerson, err := h.OpenEHRService.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrPersonNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Person not found for the given person ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Person by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Person",
			Status:  "error",
		})
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Person
	if currentPerson.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != ifMatch {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Person has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Parse updated person from request body
	var requestPerson model.PERSON
	if err := c.BodyParser(&requestPerson); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Person", "error", err)
		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	// Proceed to update Person
	updatePerson, err := h.OpenEHRService.UpdatePerson(ctx, versionedPartyID, requestPerson)
	if err != nil {
		if err == ErrPersonNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Person not found for the given person ID",
				Status:  "not_found",
			})
		}
		if err == ErrPersonVersionLowerOrEqualToCurrent {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Person version is lower than or equal to the current version",
				Status:  "conflict",
			})
		}
		if err == ErrInvalidPersonUIDMismatch {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Person UID HIER_OBJECT_ID in request body does not match current Person UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "validation_error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to update Person", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Person",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypePersonUpdated, map[string]any{"person": updatePerson})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Person updated", "error", err)
	}

	// Determine response
	updatedPersonID := updatePerson.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedPersonID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/person/"+updatedPersonID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatePerson)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusOK).JSON(`{"uid":"` + updatedPersonID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatePerson)
	}
}

func (h *Handler) DeletePerson(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	if strings.Count(uidBasedID, "::") != 3 {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Person by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	// Check existence before deletion
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourcePerson,
			Action:    audit.ActionDelete,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_party_id": versionedPartyID,
				"outcome":            outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for DeletePerson", "error", err)
		}
	}()

	currentPerson, err := h.OpenEHRService.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrPersonNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Person not found for the given person ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Person by ID before deletion", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Person by ID before deletion",
			Status:  "error",
		})
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Person
	if currentPerson.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != uidBasedID {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Person has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Proceed to delete Person
	if err := h.OpenEHRService.DeletePerson(ctx, uidBasedID); err != nil {
		if err == ErrPersonNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Person not found for the given person ID",
				Status:  "not_found",
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Person by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Person by ID",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypePersonDeleted, map[string]any{"person": currentPerson})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Person deleted", "error", err)
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) CreateOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Parse new organisation from request body
	var newOrganisation model.ORGANISATION
	if err := c.BodyParser(&newOrganisation); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Organisation", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceOrganisation,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"organisation_uid": newOrganisation.UID,
				"outcome":          outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for CreateOrganisation", "error", err)
		}
	}()

	organisation, err := h.OpenEHRService.CreateOrganisation(ctx, newOrganisation)
	if err != nil {
		if err == ErrOrganisationAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Organisation with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to create Organisation", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create Organisation",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeOrganisationCreated, map[string]any{"organisation": organisation})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Organisation created", "error", err)
	}

	// Determine response
	organisationID := organisation.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+organisationID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/organisation/"+organisationID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(organisation)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + organisationID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(organisation)
	}
}

func (h *Handler) GetOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	if strings.Count(uidBasedID, "::") == 3 {
		// Is version id
		organisation, err := h.OpenEHRService.GetOrganisationAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == ErrOrganisationNotFound {
				return SendErrorResponse(c, ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Organisation not found for the given organisation ID",
					Status:  "not_found",
				})
			}
			h.Logger.ErrorContext(ctx, "Failed to get Organisation by ID", "error", err)
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to get Organisation by ID",
				Status:  "error",
			})
		}

		return c.Status(fiber.StatusOK).JSON(organisation)
	}

	// Is versioned party id
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceOrganisation,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"organisation_uid": uidBasedID,
				"outcome":          outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetOrganisation", "error", err)
		}
	}()

	organisation, err := h.OpenEHRService.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrOrganisationNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Organisation not found for the given organisation ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Organisation by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Organisation by ID",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(organisation)
}

func (h *Handler) UpdateOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionID := c.Params("uid_based_id")
	if versionID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "If-Match header is required",
			Status:  "bad_request",
		})
	}

	// Parse return type from Prefer header
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Ensure Organisation exists before update
	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceOrganisation,
			Action:    audit.ActionUpdate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_party_id": versionedPartyID,
				"outcome":            outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for UpdateOrganisation", "error", err)
		}
	}()

	// Check collision using If-Match header
	currentOrganisation, err := h.OpenEHRService.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrOrganisationNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Organisation not found for the given organisation ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Organisation by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Organisation by ID",
			Status:  "error",
		})
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Organisation
	if currentOrganisation.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != ifMatch {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Organisation has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Parse updated organisation from request body
	var requestOrganisation model.ORGANISATION
	if err := c.BodyParser(&requestOrganisation); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Organisation", "error", err)
		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	// Proceed to update Organisation
	organisation, err := h.OpenEHRService.UpdateOrganisation(ctx, versionedPartyID, requestOrganisation)
	if err != nil {
		if err == ErrOrganisationNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Organisation not found for the given organisation ID",
				Status:  "not_found",
			})
		}
		if err == ErrOrganisationVersionLowerOrEqualToCurrent {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Organisation version is lower than or equal to the current version",
				Status:  "conflict",
			})
		}
		if err == ErrInvalidOrganisationUIDMismatch {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Organisation with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "validation_error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to update Organisation", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Organisation",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeOrganisationUpdated, map[string]any{"organisation": organisation})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Organisation updated", "error", err)
	}

	// Determine response
	organisationID := organisation.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+organisationID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/organisation/"+organisationID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(organisation)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusOK).JSON(`{"uid":"` + organisationID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(organisation)
	}
}

func (h *Handler) DeleteOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	if strings.Count(uidBasedID, "::") != 2 {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Organisation by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	// Check existence before deletion
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "error"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceOrganisation,
			Action:    audit.ActionDelete,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_party_id": versionedPartyID,
				"outcome":            outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for DeleteOrganisation", "error", err)
		}
	}()

	currentOrganisation, err := h.OpenEHRService.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrOrganisationNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Organisation not found for the given organisation ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Organisation by ID before deletion", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Organisation by ID before deletion",
			Status:  "error",
		})
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Organisation
	if currentOrganisation.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != uidBasedID {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Organisation has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Proceed to delete Organisation
	if err := h.OpenEHRService.DeleteOrganisation(ctx, uidBasedID); err != nil {
		if err == ErrOrganisationNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Organisation not found for the given organisation ID",
				Status:  "not_found",
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Organisation by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Organisation by ID",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeOrganisationDeleted, map[string]any{"organisation": currentOrganisation})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Organisation deleted", "error", err)
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) CreateRole(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	// Parse return type from Prefer header
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Parse new role from request body
	var newRole model.ROLE
	if err := c.BodyParser(&newRole); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Role", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceRole,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"role_uid": newRole.UID,
				"outcome":  outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for CreateRole", "error", err)
		}
	}()

	// Proceed to create Role
	role, err := h.OpenEHRService.CreateRole(ctx, newRole)
	if err != nil {
		if err == ErrRoleAlreadyExists {
			outcome = "conflict"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Role with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err, ok := err.(util.ValidateError); ok {
			outcome = "bad_request"
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to create Role", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create Role",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeRoleCreated, map[string]any{"role": role})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Role created", "error", err)
	}

	// Determine response
	roleID := role.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+roleID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/role/"+roleID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(role)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + roleID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(role)
	}
}

func (h *Handler) GetRole(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceRole,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"role_uid": uidBasedID,
				"outcome":  outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetRole", "error", err)
		}
	}()

	if strings.Count(uidBasedID, "::") == 3 {
		// Is version id
		role, err := h.OpenEHRService.GetRoleAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == ErrRoleNotFound {
				outcome = "not_found"
				return SendErrorResponse(c, ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Role not found for the given role ID",
					Status:  "not_found",
				})
			}
			h.Logger.ErrorContext(ctx, "Failed to get Role by ID", "error", err)
			outcome = "error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to get Role by ID",
				Status:  "error",
			})
		}

		outcome = "success"
		return c.Status(fiber.StatusOK).JSON(role)
	}

	// Is versioned party id
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	role, err := h.OpenEHRService.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrRoleNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Role not found for the given role ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Role by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Role by ID",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(role)
}

func (h *Handler) UpdateRole(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionID := c.Params("uid_based_id")
	if versionID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "If-Match header is required",
			Status:  "bad_request",
		})
	}

	// Parse return type from Prefer header
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header value",
			Status:  "bad_request",
		})
	}

	// Ensure Role exists before update
	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceRole,
			Action:    audit.ActionUpdate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_party_id": versionedPartyID,
				"outcome":            outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for UpdateRole", "error", err)
		}
	}()

	// Check collision using If-Match header
	currentRole, err := h.OpenEHRService.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrRoleNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Role not found for the given role ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Role by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Role by ID",
			Status:  "error",
		})
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Role
	if currentRole.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != ifMatch {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Role has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Parse updated role from request body
	var requestRole model.ROLE
	if err := c.BodyParser(&requestRole); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			outcome = "validation_error"
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Role", "error", err)
		outcome = "bad_request"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	// Proceed to update Role
	updatedRole, err := h.OpenEHRService.UpdateRole(ctx, versionedPartyID, requestRole)
	if err != nil {
		if err == ErrRoleNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Role not found for the given role ID",
				Status:  "not_found",
			})
		}
		if err == ErrRoleVersionLowerOrEqualToCurrent {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Role version is lower than or equal to the current version",
				Status:  "conflict",
			})
		}
		if err == ErrInvalidRoleUIDMismatch {
			outcome = "bad_request"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Role UID HIER_OBJECT_ID in request body does not match current Role UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			outcome = "validation_error"
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to update Role", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Role",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeRoleUpdated, map[string]any{"role": updatedRole})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Role updated", "error", err)
	}

	// Determine response
	updatedRoleID := updatedRole.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedRoleID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/role/"+updatedRoleID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedRole)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusOK).JSON(`{"uid":"` + updatedRoleID + `"}`)
	default:
		h.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedRole)
	}
}

func (h *Handler) DeleteRole(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "uid_based_id parameter is required",
			Status:  "bad_request",
		})
	}

	if !strings.Contains(uidBasedID, "::") {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Role by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	// Check existence before deletion
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceRole,
			Action:    audit.ActionDelete,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_party_id": versionedPartyID,
				"outcome":            outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for DeleteRole", "error", err)
		}
	}()

	currentRole, err := h.OpenEHRService.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrRoleNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Role not found for the given role ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Role by ID before deletion", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Role by ID before deletion",
			Status:  "error",
		})
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Role
	if currentRole.UID.V.Value.(*model.OBJECT_VERSION_ID).Value != uidBasedID {
		outcome = "precondition_failed"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Role has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	// Proceed to delete Role
	if err := h.OpenEHRService.DeleteRole(ctx, uidBasedID); err != nil {
		if err == ErrRoleNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Role not found for the given role ID",
				Status:  "not_found",
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Role by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Role by ID",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeRoleDeleted, map[string]any{"role": currentRole})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Role deleted", "error", err)
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) GetVersionedParty(c *fiber.Ctx) error {
	ctx := c.Context()

	versionedObjectIDStr := c.Params("versioned_object_id")
	if versionedObjectIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "versioned_object_id parameter is required",
			Status:  "bad_request",
		})
	}
	versionedObjectID, err := uuid.Parse(versionedObjectIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid versioned_object_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceVersionedParty,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_object_id": versionedObjectID,
				"outcome":             outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedParty", "error", err)
		}
	}()

	party, err := h.OpenEHRService.GetVersionedParty(ctx, versionedObjectID)
	if err != nil {
		if err == ErrVersionedPartyNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party not found for the given ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Versioned Party by ID",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(party)
}

func (h *Handler) GetVersionedPartyRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	versionedObjectIDStr := c.Params("versioned_object_id")
	if versionedObjectIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "versioned_object_id parameter is required",
			Status:  "bad_request",
		})
	}
	versionedObjectID, err := uuid.Parse(versionedObjectIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid versioned_object_id format",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceVersionedParty,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_object_id": versionedObjectID,
				"outcome":             outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedPartyRevisionHistory", "error", err)
		}
	}()

	revisionHistory, err := h.OpenEHRService.GetVersionedPartyRevisionHistory(ctx, versionedObjectID)
	if err != nil {
		if err == ErrRevisionHistoryNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party not found for the given ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party Revision History by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Versioned Party Revision History by ID",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(revisionHistory)
}

func (h *Handler) GetVersionedPartyVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()

	versionedObjectIDStr := c.Params("versioned_object_id")
	if versionedObjectIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "versioned_object_id parameter is required",
			Status:  "bad_request",
		})
	}
	versionedObjectID, err := uuid.Parse(versionedObjectIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid versioned_object_id format",
			Status:  "bad_request",
		})
	}

	// Parse version_at_time query parameter
	var filterAtTime time.Time
	atTimeStr := c.Query("version_at_time")
	if atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid version_at_time format. Use RFC3339 format.",
				Status:  "bad_request",
			})
		}
		filterAtTime = parsedTime
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceVersionedPartyVersion,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_object_id": versionedObjectID,
				"outcome":             outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedPartyVersionAtTime", "error", err)
		}
	}()

	partyVersionJSON, err := h.OpenEHRService.GetVersionedPartyVersionJSON(ctx, versionedObjectID, filterAtTime, "")
	if err != nil {
		if err == ErrVersionedPartyNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party not found for the given ID",
				Status:  "not_found",
			})
		}
		if err == ErrVersionedPartyVersionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party version not found for the given time",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party Version at Time by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Versioned Party Version at Time by ID",
			Status:  "error",
		})
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(partyVersionJSON)
}

func (h *Handler) GetVersionedPartyVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	versionedObjectIDStr := c.Params("versioned_object_id")
	if versionedObjectIDStr == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "versioned_object_id parameter is required",
			Status:  "bad_request",
		})
	}
	versionedObjectID, err := uuid.Parse(versionedObjectIDStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid versioned_object_id format",
			Status:  "bad_request",
		})
	}

	versionID := c.Params("version_id")
	if versionID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version_id parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceVersionedPartyVersion,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"versioned_object_id": versionedObjectID,
				"version_id":          versionID,
				"outcome":             outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for GetVersionedPartyVersion", "error", err)
		}
	}()

	partyVersionJSON, err := h.OpenEHRService.GetVersionedPartyVersionJSON(ctx, versionedObjectID, time.Time{}, versionID)
	if err != nil {
		if err == ErrVersionedPartyNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party not found for the given ID",
				Status:  "not_found",
			})
		}
		if err == ErrVersionedPartyVersionNotFound {
			outcome = "not_found"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party version not found for the given version ID",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party Version by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Versioned Party Version by ID",
			Status:  "error",
		})
	}

	outcome = "success"
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(partyVersionJSON)
}

func (h *Handler) CreateDemographicContribution(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Create Demographic Contribution not implemented yet")

	// ctx := c.Context()

	// c.Accepts("application/json")

	// // Parse return type from Prefer header
	// returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	// if !returnType.IsValid() {
	// 	return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	// }

	// // Parse new contribution from request body
	// var newContribution openehr.CONTRIBUTION
	// if err := c.BodyParser(&newContribution); err != nil {
	// 	h.Logger.ErrorContext(ctx, "Failed to parse request body for new Demographic Contribution", "error", err)
	// 	return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	// }

	// validateErr := newContribution.Validate("$")
	// if len(validateErr.Errs) > 0 {
	// 	return c.Status(fiber.StatusBadRequest).JSON(validateErr)
	// }

	// contribution, err := h.OpenEHRService.CreateContribution(ctx, newContribution)
	// if err != nil {
	// 	h.Logger.ErrorContext(ctx, "Failed to create Demographic Contribution", "error", err)
	// 	c.Status(http.StatusInternalServerError)
	// 	return nil
	// }

	// // Determine response
	// contributionID := contribution.UID.Value
	// c.Set("ETag", "\""+contributionID+"\"")
	// c.Set("Location", c.Protocol()+ "://" + c.Hostname() + "/openehr/v1/demographic/contribution/"+contributionID)

	// c.Status(fiber.StatusCreated)
	// switch returnType {
	// case ReturnTypeMinimal:
	// 	return nil
	// case ReturnTypeRepresentation:
	// 	return c.JSON(contribution)
	// case ReturnTypeIdentifier:
	// 	return c.JSON(`{"uid":"` + contributionID + `"}`)
	// default:
	// 	h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
	// 	return c.JSON(contribution)
	// }
}

func (h *Handler) GetDemographicContribution(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Demographic Contribution not implemented yet")
}

func (h *Handler) GetDemographicTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Demographic Tags not implemented yet")
}

func (h *Handler) GetAgentTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Agent Tags not implemented yet")
}

func (h *Handler) UpdateAgentTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Update Agent Tags not implemented yet")
}

func (h *Handler) DeleteAgentTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Delete Agent Tag By Key not implemented yet")
}

func (h *Handler) GetGroupTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Group Tags not implemented yet")
}

func (h *Handler) UpdateGroupTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Update Group Tags not implemented yet")
}

func (h *Handler) DeleteGroupTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Delete Group Tag By Key not implemented yet")
}

func (h *Handler) GetPersonTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Person Tags not implemented yet")
}

func (h *Handler) UpdatePersonTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Update Person Tags not implemented yet")
}

func (h *Handler) DeletePersonTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Delete Person Tag By Key not implemented yet")
}

func (h *Handler) GetOrganisationTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Organisation Tags not implemented yet")
}

func (h *Handler) UpdateOrganisationTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Update Organisation Tags not implemented yet")
}

func (h *Handler) DeleteOrganisationTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Delete Organisation Tag By Key not implemented yet")
}

func (h *Handler) GetRoleTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Role Tags not implemented yet")
}

func (h *Handler) UpdateRoleTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Update Role Tags not implemented yet")
}

func (h *Handler) DeleteRoleTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Delete Role Tag By Key not implemented yet")
}

func (h *Handler) ExecuteAdHocAQL(c *fiber.Ctx) error {
	ctx := c.Context()

	query := c.Query("q")
	if query == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "q query parameter is required",
			Status:  "bad_request",
		})
	}

	ehrID := c.Query("ehr_id")
	if ehrID != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	fetch := c.Query("fetch")
	if fetch != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	offset := c.Query("offset")
	if offset != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	queryParametersStr := c.Query("query_parameters")
	queryParametersURLValues, err := url.ParseQuery(queryParametersStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid query_parameters format",
			Status:  "bad_request",
		})
	}

	queryParameters := make(map[string]any)
	for key, values := range queryParametersURLValues {
		if len(values) > 0 {
			queryParameters[key] = values[0]
		}
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceQuery,
			Action:    audit.ActionExecute,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"query":      query,
				"parameters": queryParameters,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Execute Ad Hoc AQL", "error", err)
		}
	}()

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), query, queryParameters); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Ad Hoc AQL", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Ad Hoc AQL",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeQueryExecuted, map[string]any{
		"query":      query,
		"parameters": queryParameters,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Execute Ad Hoc AQL", "error", err)
	}

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return nil
}

type AdHocAQLRequest struct {
	Query           string         `json:"q"`
	EHRID           string         `json:"ehr_id,omitempty"`
	Fetch           string         `json:"fetch,omitempty"`
	Offset          int            `json:"offset,omitempty"`
	QueryParameters map[string]any `json:"query_parameters,omitempty"`
}

func (h *Handler) ExecuteAdHocAQLPost(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	var aqlRequest AdHocAQLRequest
	if err := c.BodyParser(&aqlRequest); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse AQL request body", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	if aqlRequest.Query == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "query field is required in the request body",
			Status:  "bad_request",
		})
	}

	if aqlRequest.EHRID != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL Post with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Fetch != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL Post with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Offset != 0 {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL Post with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceQuery,
			Action:    audit.ActionExecute,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"query":      aqlRequest.Query,
				"parameters": aqlRequest.QueryParameters,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Execute Ad Hoc AQL Post", "error", err)
		}
	}()

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), aqlRequest.Query, aqlRequest.QueryParameters); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Ad Hoc AQL Post", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Ad Hoc AQL Post",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err := h.WebhookService.RegisterEvent(ctx, webhook.EventTypeQueryExecuted, map[string]any{
		"query":      aqlRequest.Query,
		"parameters": aqlRequest.QueryParameters,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Execute Ad Hoc AQL Post", "error", err)
	}

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) ExecuteStoredAQL(c *fiber.Ctx) error {
	ctx := c.Context()

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	ehrID := c.Query("ehr_id")
	if ehrID != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	fetch := c.Query("fetch")
	if fetch != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	offset := c.Query("offset")
	if offset != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	queryParametersStr := c.Query("query_parameters")
	queryParametersURLValues, err := url.ParseQuery(queryParametersStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid query_parameters format",
			Status:  "bad_request",
		})
	}

	queryParameters := make(map[string]any)
	for key, values := range queryParametersURLValues {
		if len(values) > 0 {
			queryParameters[key] = values[0]
		}
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceQuery,
			Action:    audit.ActionExecute,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"query_name": name,
				"parameters": queryParameters,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Execute Stored AQL", "error", err)
		}
	}()

	// Retrieve stored query by name
	storedQuery, err := h.OpenEHRService.GetQueryByName(ctx, name, "")
	if err != nil {
		if err == ErrQueryNotFound {
			outcome = "error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Stored query not found for the given name",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get stored query by name", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get stored query by name",
			Status:  "error",
		})
	}

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), storedQuery.Query, nil); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Stored AQL", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Stored AQL",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeQueryExecuted, map[string]any{
		"query_name": name,
		"parameters": queryParameters,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Execute Stored AQL", "error", err)
	}

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return nil
}

type StoredAQLRequest struct {
	EHRID           string         `json:"ehr_id,omitempty"`
	Fetch           string         `json:"fetch,omitempty"`
	Offset          int            `json:"offset,omitempty"`
	QueryParameters map[string]any `json:"query_parameters,omitempty"`
}

func (h *Handler) ExecuteStoredAQLPost(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	var aqlRequest StoredAQLRequest
	if err := c.BodyParser(&aqlRequest); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse Stored AQL request body", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	if aqlRequest.EHRID != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Post with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Fetch != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Post with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Offset != 0 {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Post with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceQuery,
			Action:    audit.ActionExecute,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"query_name": name,
				"parameters": aqlRequest.QueryParameters,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Execute Stored AQL Post", "error", err)
		}
	}()

	// Retrieve stored query by name
	storedQuery, err := h.OpenEHRService.GetQueryByName(ctx, name, "")
	if err != nil {
		if err == ErrQueryNotFound {
			outcome = "error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Stored query not found for the given name",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get stored query by name", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get stored query by name",
			Status:  "error",
		})
	}

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), storedQuery.Query, aqlRequest.QueryParameters); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Stored AQL Post", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Stored AQL Post",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeQueryExecuted, map[string]any{
		"query_name": name,
		"parameters": aqlRequest.QueryParameters,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Execute Stored AQL Post", "error", err)
	}

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) ExecuteStoredAQLVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	version := c.Params("version")
	if version == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version path parameter is required",
			Status:  "bad_request",
		})
	}

	ehrID := c.Query("ehr_id")
	if ehrID != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	fetch := c.Query("fetch")
	if fetch != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	offset := c.Query("offset")
	if offset != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	queryParametersStr := c.Query("query_parameters")
	queryParametersURLValues, err := url.ParseQuery(queryParametersStr)
	if err != nil {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid query_parameters format",
			Status:  "bad_request",
		})
	}

	queryParameters := make(map[string]any)
	for key, values := range queryParametersURLValues {
		if len(values) > 0 {
			queryParameters[key] = values[0]
		}
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceQuery,
			Action:    audit.ActionExecute,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"query_name": name,
				"version":    version,
				"parameters": queryParameters,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Execute Stored AQL Version", "error", err)
		}
	}()

	// Retrieve stored query by name and version
	storedQuery, err := h.OpenEHRService.GetQueryByName(ctx, name, version)
	if err != nil {
		if err == ErrQueryNotFound {
			outcome = "error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Stored query not found for the given name and version",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get stored query by name and version", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get stored query by name and version",
			Status:  "error",
		})
	}

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), storedQuery.Query, queryParameters); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Stored AQL Version", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Stored AQL Version",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeQueryExecuted, map[string]any{
		"query_name": name,
		"version":    version,
		"parameters": queryParameters,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Execute Stored AQL Version", "error", err)
	}

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return nil
}

type StoredAQLVersionRequest struct {
	EHRID           string         `json:"ehr_id,omitempty"`
	Fetch           string         `json:"fetch,omitempty"`
	Offset          int            `json:"offset,omitempty"`
	QueryParameters map[string]any `json:"query_parameters,omitempty"`
}

func (h *Handler) ExecuteStoredAQLVersionPost(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	version := c.Params("version")
	if version == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version path parameter is required",
			Status:  "bad_request",
		})
	}

	var aqlRequest StoredAQLVersionRequest
	if err := c.BodyParser(&aqlRequest); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse Stored AQL Version request body", "error", err)
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}

	if aqlRequest.EHRID != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version Post with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Fetch != "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version Post with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Offset != 0 {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version Post with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceQuery,
			Action:    audit.ActionExecute,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"query_name": name,
				"version":    version,
				"parameters": aqlRequest.QueryParameters,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Execute Stored AQL Version Post", "error", err)
		}
	}()

	// Retrieve stored query by name and version
	storedQuery, err := h.OpenEHRService.GetQueryByName(ctx, name, version)
	if err != nil {
		if err == ErrQueryNotFound {
			outcome = "error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Stored query not found for the given name and version",
				Status:  "not_found",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get stored query by name and version", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get stored query by name and version",
			Status:  "error",
		})
	}

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), storedQuery.Query, aqlRequest.QueryParameters); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Stored AQL Version Post", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Stored AQL Version Post",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeQueryExecuted, map[string]any{
		"query_name": name,
		"version":    version,
		"parameters": aqlRequest.QueryParameters,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Execute Stored AQL Version Post", "error", err)
	}

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) GetTemplatesADL14(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Templates ADL1.4 not implemented yet")
}

func (h *Handler) UploadTemplateADL14(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Upload Template ADL1.4 not implemented yet")
}

func (h *Handler) GetTemplateADL14ByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Template ADL1.4 By ID not implemented yet")
}

func (h *Handler) GetTemplatesADL2(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Templates ADL2 not implemented yet")
}

func (h *Handler) UploadTemplateADL2(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Upload Template ADL2 not implemented yet")
}

func (h *Handler) GetTemplateADL2ByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Template ADL2 By ID not implemented yet")
}

func (h *Handler) GetTemplateADL2AtVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Template ADL2 At Version not implemented yet")
}

func (h *Handler) ListStoredQueries(c *fiber.Ctx) error {
	ctx := c.Context()

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceQuery,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"query_name": name,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for List Stored Queries", "error", err)
		}
	}()

	queries, err := h.OpenEHRService.ListStoredQueries(ctx, name)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to list stored queries", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to list stored queries",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(queries)
}

func (h *Handler) StoreQuery(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("text/plain")

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	queryType := c.Query("query_type")
	if queryType != "" && queryType != "AQL" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Unsupported query_type. Only 'AQL' is supported.",
			Status:  "bad_request",
		})
	}

	query := string(c.Body())
	if query == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Query in request body is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceQuery,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"query_name": name,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Store Query", "error", err)
		}
	}()

	// Check if query with the same name already exists
	_, err := h.OpenEHRService.GetQueryByName(ctx, name, "")
	if err == nil {
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusConflict,
			Message: "Query with the given name already exists, system cannot update without knowing the target version, please use Store Query Version endpoint instead",
			Status:  "conflict",
		})
	}
	if err != ErrQueryNotFound {
		h.Logger.ErrorContext(ctx, "Failed to check existing query by name", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to store query",
			Status:  "error",
		})
	}

	// Store the query
	err = h.OpenEHRService.StoreQuery(ctx, name, "1.0.0", query)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to store query", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to store query",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeQueryStored, map[string]any{
		"query_name": name,
		"version":    "1.0.0",
		"query":      query,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Store Query", "error", err)
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) StoreQueryVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("text/plain")

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	version := c.Params("version")
	if version == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version path parameter is required",
			Status:  "bad_request",
		})
	}

	queryType := c.Query("query_type")
	if queryType != "" && queryType != "AQL" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Unsupported query_type. Only 'AQL' is supported.",
			Status:  "bad_request",
		})
	}

	query := string(c.Body())
	if query == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Query in request body is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceQuery,
			Action:    audit.ActionCreate,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"query_name": name,
				"version":    version,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Store Query Version", "error", err)
		}
	}()

	// Store the new version of the query
	err := h.OpenEHRService.StoreQuery(ctx, name, version, query)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to store query version", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to store query version",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeQueryStored, map[string]any{
		"query_name": name,
		"version":    version,
		"query":      query,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Store Query Version", "error", err)
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) GetStoredQueryAtVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	version := c.Params("version")
	if version == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version path parameter is required",
			Status:  "bad_request",
		})
	}

	outcome := "unknown"
	defer func() {
		if err := h.AuditService.LogEvent(ctx, audit.LogEventRequest{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  audit.ResourceQuery,
			Action:    audit.ActionRead,
			Success:   outcome == "success",
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Details: map[string]any{
				"query_name": name,
				"version":    version,
				"outcome":    outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Get Stored Query At Version", "error", err)
		}
	}()

	storedQuery, err := h.OpenEHRService.GetQueryByName(ctx, name, version)
	if err != nil {
		if err == ErrQueryNotFound {
			outcome = "error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Stored query not found for the given name and version",
				Status:  "error",
			})
		}
		h.Logger.ErrorContext(ctx, "Failed to get stored query by name and version", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get stored query by name and version",
			Status:  "error",
		})
	}

	outcome = "success"
	return c.Status(fiber.StatusOK).JSON(storedQuery)
}

func (h *Handler) DeleteEHRByID(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "ehr_id parameter is required",
			Status:  "bad_request",
		})
	}

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
				"ehr_id":  ehrID,
				"outcome": outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Delete EHR By ID", "error", err)
		}
	}()

	err := h.OpenEHRService.DeleteEHR(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			outcome = "error"
			return SendErrorResponse(c, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR with the given ID not found",
				Status:  "error",
			})
		}

		h.Logger.ErrorContext(ctx, "Failed to delete EHR by ID", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete EHR by ID",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeEHRDeleted, map[string]any{
		"ehr_id": ehrID,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Delete EHR By ID", "error", err)
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) DeleteMultipleEHRs(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse multiple ehr_id query parameters
	// Example: ?ehr_id=7d44b88c-4199-4bad-97dc-d78268e01398&ehr_id=297c3e91-7c17-4497-85dd-01e05aaae44e
	var ehrIDList []string
	for key, value := range c.Context().QueryArgs().All() {
		if string(key) == "ehr_id" {
			ehrIDList = append(ehrIDList, string(value))
		}
	}

	if len(ehrIDList) == 0 {
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "At least one ehr_id query parameter is required",
			Status:  "bad_request",
		})
	}

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
				"ehr_id_list": ehrIDList,
				"outcome":     outcome,
			},
		}); err != nil {
			h.Logger.ErrorContext(ctx, "Failed to log audit event for Delete Multiple EHRs", "error", err)
		}
	}()

	err := h.OpenEHRService.DeleteEHRBulk(ctx, ehrIDList)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to delete multiple EHRs", "error", err)
		outcome = "error"
		return SendErrorResponse(c, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete multiple EHRs",
			Status:  "error",
		})
	}

	outcome = "success"

	// Register event
	err = h.WebhookService.RegisterEvent(ctx, webhook.EventTypeEHRDeleted, map[string]any{
		"ehr_id_list": ehrIDList,
	})
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to notify webhook for Delete Multiple EHRs", "error", err)
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

type ReturnType string

const (
	ReturnTypeMinimal        ReturnType = "return=minimal"
	ReturnTypeRepresentation ReturnType = "return=representation"
	ReturnTypeIdentifier     ReturnType = "return=identifier"
)

func (r ReturnType) IsValid() bool {
	switch r {
	case ReturnTypeMinimal, ReturnTypeRepresentation, ReturnTypeIdentifier:
		return true
	default:
		return false
	}
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
