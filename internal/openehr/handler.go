package openehr

import (
	"net/url"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/audit"
	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/openehr/model"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/internal/webhook"
	"github.com/freekieb7/gopenehr/pkg/utils"
	"github.com/freekieb7/gopenehr/pkg/web/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	Settings       *config.Settings
	Telemetry      *telemetry.Telemetry
	OpenEHRService *Service
	AuditService   *audit.Service
	WebhookService *webhook.Service
	AuditLogger    *audit.Logger
	WebhookSaver   *webhook.Saver
}

func NewHandler(settings *config.Settings, telemetry *telemetry.Telemetry, openEHRService *Service, auditService *audit.Service, webhookService *webhook.Service, auditLogger *audit.Logger, webhookSaver *webhook.Saver) Handler {
	return Handler{
		Settings:       settings,
		Telemetry:      telemetry,
		OpenEHRService: openEHRService,
		AuditService:   auditService,
		WebhookService: webhookService,
		AuditLogger:    auditLogger,
		WebhookSaver:   webhookSaver,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/openehr/v1")
	v1.Use(middleware.NoCache)

	v1.Use(middleware.Telemetry(h.Telemetry))
	v1.Use(middleware.APIKeyProtected(h.Settings.APIKey))

	v1.Options("", h.SystemInfo)

	v1.Get("/ehr", audit.Middleware(h.AuditLogger, audit.ResourceEHR, audit.ActionRead), h.GetEHRBySubject)
	v1.Post("/ehr", audit.Middleware(h.AuditLogger, audit.ResourceEHR, audit.ActionCreate), h.CreateEHR)
	v1.Get("/ehr/:ehr_id", audit.Middleware(h.AuditLogger, audit.ResourceEHR, audit.ActionRead), h.GetEHR)
	v1.Put("/ehr/:ehr_id", audit.Middleware(h.AuditLogger, audit.ResourceEHR, audit.ActionCreate), h.CreateEHRWithID)

	v1.Get("/ehr/:ehr_id/ehr_status", audit.Middleware(h.AuditLogger, audit.ResourceEHRStatus, audit.ActionRead), h.GetEHRStatus)
	v1.Put("/ehr/:ehr_id/ehr_status", audit.Middleware(h.AuditLogger, audit.ResourceEHRStatus, audit.ActionUpdate), h.UpdateEhrStatus)
	v1.Get("/ehr/:ehr_id/ehr_status/:version_uid", audit.Middleware(h.AuditLogger, audit.ResourceEHRStatus, audit.ActionRead), h.GetEHRStatusByVersionID)

	v1.Get("/ehr/:ehr_id/versioned_ehr_status", audit.Middleware(h.AuditLogger, audit.ResourceVersionedEHRStatus, audit.ActionRead), h.GetVersionedEHRStatus)
	v1.Get("/ehr/:ehr_id/versioned_ehr_status/revision_history", audit.Middleware(h.AuditLogger, audit.ResourceVersionedEHRStatus, audit.ActionRead), h.GetVersionedEHRStatusRevisionHistory)
	v1.Get("/ehr/:ehr_id/versioned_ehr_status/version", audit.Middleware(h.AuditLogger, audit.ResourceVersionedEHRStatusVersion, audit.ActionRead), h.GetVersionedEHRStatusVersion)
	v1.Get("/ehr/:ehr_id/versioned_ehr_status/version/:version_uid", audit.Middleware(h.AuditLogger, audit.ResourceVersionedEHRStatusVersion, audit.ActionRead), h.GetVersionedEHRStatusVersionByID)

	v1.Post("/ehr/:ehr_id/composition", audit.Middleware(h.AuditLogger, audit.ResourceComposition, audit.ActionCreate), h.CreateComposition)
	v1.Get("/ehr/:ehr_id/composition/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceComposition, audit.ActionRead), h.GetComposition)
	v1.Put("/ehr/:ehr_id/composition/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceComposition, audit.ActionUpdate), h.UpdateComposition)
	v1.Delete("/ehr/:ehr_id/composition/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceComposition, audit.ActionDelete), h.DeleteComposition)

	v1.Get("/ehr/:ehr_id/versioned_composition/:versioned_object_id", audit.Middleware(h.AuditLogger, audit.ResourceVersionedComposition, audit.ActionRead), h.GetVersionedCompositionByID)
	v1.Get("/ehr/:ehr_id/versioned_composition/:versioned_object_id/revision_history", audit.Middleware(h.AuditLogger, audit.ResourceVersionedComposition, audit.ActionRead), h.GetVersionedCompositionRevisionHistory)
	v1.Get("/ehr/:ehr_id/versioned_composition/:versioned_object_id/version", audit.Middleware(h.AuditLogger, audit.ResourceVersionedCompositionVersion, audit.ActionRead), h.GetVersionedCompositionVersionAtTime)
	v1.Get("/ehr/:ehr_id/versioned_composition/:versioned_object_id/version/:version_uid", audit.Middleware(h.AuditLogger, audit.ResourceVersionedCompositionVersion, audit.ActionRead), h.GetVersionedCompositionVersionByID)

	v1.Post("/ehr/:ehr_id/directory", audit.Middleware(h.AuditLogger, audit.ResourceDirectory, audit.ActionCreate), h.CreateDirectory)
	v1.Put("/ehr/:ehr_id/directory", audit.Middleware(h.AuditLogger, audit.ResourceDirectory, audit.ActionUpdate), h.UpdateDirectory)
	v1.Delete("/ehr/:ehr_id/directory", audit.Middleware(h.AuditLogger, audit.ResourceDirectory, audit.ActionDelete), h.DeleteDirectory)
	v1.Get("/ehr/:ehr_id/directory", audit.Middleware(h.AuditLogger, audit.ResourceDirectory, audit.ActionRead), h.GetFolderInDirectoryVersionAtTime)
	v1.Get("/ehr/:ehr_id/directory/:version_uid", audit.Middleware(h.AuditLogger, audit.ResourceDirectory, audit.ActionRead), h.GetFolderInDirectoryVersion)

	v1.Post("/ehr/:ehr_id/contribution", audit.Middleware(h.AuditLogger, audit.ResourceContribution, audit.ActionCreate), h.CreateContribution)
	v1.Get("/ehr/:ehr_id/contribution/:contribution_uid", audit.Middleware(h.AuditLogger, audit.ResourceContribution, audit.ActionRead), h.GetContribution)

	// v1.Get("/ehr/:ehr_id/tags", h.GetEHRTags)
	// v1.Get("/ehr/:ehr_id/composition/:uid_based_id/tags", h.GetCompositionTags)
	// v1.Put("/ehr/:ehr_id/composition/:uid_based_id/tags", h.UpdateCompositionTags)
	// v1.Delete("/ehr/:ehr_id/composition/:uid_based_id/tags", h.DeleteCompositionTagByKey)
	// v1.Get("/ehr/:ehr_id/ehr_status/tags", h.GetEHRStatusTags)
	// v1.Get("/ehr/:ehr_id/ehr_status/:version_uid/tags", h.GetEHRStatusVersionTags)
	// v1.Put("/ehr/:ehr_id/ehr_status/:version_uid/tags", h.UpdateEHRStatusVersionTags)
	// v1.Delete("/ehr/:ehr_id/ehr_status/:version_uid/tags/:key", h.DeleteEHRStatusVersionTagByKey)

	v1.Post("/demographic/agent", audit.Middleware(h.AuditLogger, audit.ResourceAgent, audit.ActionCreate), h.CreateAgent)
	v1.Get("/demographic/agent/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceAgent, audit.ActionRead), h.GetAgent)
	v1.Put("/demographic/agent/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceAgent, audit.ActionUpdate), h.UpdateAgent)
	v1.Delete("/demographic/agent/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceAgent, audit.ActionDelete), h.DeleteAgent)

	v1.Post("/demographic/group", audit.Middleware(h.AuditLogger, audit.ResourceGroup, audit.ActionCreate), h.CreateGroup)
	v1.Get("/demographic/group/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceGroup, audit.ActionRead), h.GetGroup)
	v1.Put("/demographic/group/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceGroup, audit.ActionUpdate), h.UpdateGroup)
	v1.Delete("/demographic/group/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceGroup, audit.ActionDelete), h.DeleteGroup)

	v1.Post("/demographic/person", audit.Middleware(h.AuditLogger, audit.ResourcePerson, audit.ActionCreate), h.CreatePerson)
	v1.Get("/demographic/person/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourcePerson, audit.ActionRead), h.GetPerson)
	v1.Put("/demographic/person/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourcePerson, audit.ActionUpdate), h.UpdatePerson)
	v1.Delete("/demographic/person/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourcePerson, audit.ActionDelete), h.DeletePerson)

	v1.Post("/demographic/organisation", audit.Middleware(h.AuditLogger, audit.ResourceOrganisation, audit.ActionCreate), h.CreateOrganisation)
	v1.Get("/demographic/organisation/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceOrganisation, audit.ActionRead), h.GetOrganisation)
	v1.Put("/demographic/organisation/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceOrganisation, audit.ActionUpdate), h.UpdateOrganisation)
	v1.Delete("/demographic/organisation/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceOrganisation, audit.ActionDelete), h.DeleteOrganisation)

	v1.Post("/demographic/role", audit.Middleware(h.AuditLogger, audit.ResourceRole, audit.ActionCreate), h.CreateRole)
	v1.Get("/demographic/role/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceRole, audit.ActionRead), h.GetRole)
	v1.Put("/demographic/role/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceRole, audit.ActionUpdate), h.UpdateRole)
	v1.Delete("/demographic/role/:uid_based_id", audit.Middleware(h.AuditLogger, audit.ResourceRole, audit.ActionDelete), h.DeleteRole)

	v1.Get("/demographic/versioned_party/:versioned_object_id", audit.Middleware(h.AuditLogger, audit.ResourceVersionedParty, audit.ActionRead), h.GetVersionedParty)
	v1.Get("/demographic/versioned_party/:versioned_object_id/revision_history", audit.Middleware(h.AuditLogger, audit.ResourceVersionedParty, audit.ActionRead), h.GetVersionedPartyRevisionHistory)
	v1.Get("/demographic/versioned_party/:versioned_object_id/version", audit.Middleware(h.AuditLogger, audit.ResourceVersionedPartyVersion, audit.ActionRead), h.GetVersionedPartyVersionAtTime)
	v1.Get("/demographic/versioned_party/:versioned_object_id/version/:version_id", audit.Middleware(h.AuditLogger, audit.ResourceVersionedPartyVersion, audit.ActionRead), h.GetVersionedPartyVersion)

	v1.Post("/demographic/contribution", audit.Middleware(h.AuditLogger, audit.ResourceContribution, audit.ActionCreate), h.CreateDemographicContribution)
	v1.Get("/demographic/contribution/:contribution_uid", audit.Middleware(h.AuditLogger, audit.ResourceContribution, audit.ActionRead), h.GetDemographicContribution)

	// v1.Get("/demographic/tags", h.GetDemographicTags)
	// v1.Get("/demographic/agent/:uid_based_id/tags", h.GetAgentTags)
	// v1.Put("/demographic/agent/:uid_based_id/tags", h.UpdateAgentTags)
	// v1.Delete("/demographic/agent/:uid_based_id/tags/:key", h.DeleteAgentTagByKey)
	// v1.Get("/demographic/group/:uid_based_id/tags", h.GetGroupTags)
	// v1.Put("/demographic/group/:uid_based_id/tags", h.UpdateGroupTags)
	// v1.Delete("/demographic/group/:uid_based_id/tags/:key", h.DeleteGroupTagByKey)
	// v1.Get("/demographic/person/:uid_based_id/tags", h.GetPersonTags)
	// v1.Put("/demographic/person/:uid_based_id/tags", h.UpdatePersonTags)
	// v1.Delete("/demographic/person/:uid_based_id/tags/:key", h.DeletePersonTagByKey)
	// v1.Get("/demographic/organisation/:uid_based_id/tags", h.GetOrganisationTags)
	// v1.Put("/demographic/organisation/:uid_based_id/tags", h.UpdateOrganisationTags)
	// v1.Delete("/demographic/organisation/:uid_based_id/tags/:key", h.DeleteOrganisationTagByKey)
	// v1.Get("/demographic/role/:uid_based_id/tags", h.GetRoleTags)
	// v1.Put("/demographic/role/:uid_based_id/tags", h.UpdateRoleTags)
	// v1.Delete("/demographic/role/:uid_based_id/tags/:key", h.DeleteRoleTagByKey)

	v1.Get("/query/aql", audit.Middleware(h.AuditLogger, audit.ResourceQuery, audit.ActionExecute), h.ExecuteAdHocAQL)
	v1.Post("/query/aql", audit.Middleware(h.AuditLogger, audit.ResourceQuery, audit.ActionExecute), h.ExecuteAdHocAQLPost)
	v1.Get("/query/:qualified_query_name", audit.Middleware(h.AuditLogger, audit.ResourceQuery, audit.ActionExecute), h.ExecuteStoredAQL)
	v1.Post("/query/:qualified_query_name", audit.Middleware(h.AuditLogger, audit.ResourceQuery, audit.ActionExecute), h.ExecuteStoredAQLPost)
	v1.Get("/query/:qualified_query_name/:version", audit.Middleware(h.AuditLogger, audit.ResourceQuery, audit.ActionExecute), h.ExecuteStoredAQLVersion)
	v1.Post("/query/:qualified_query_name/:version", audit.Middleware(h.AuditLogger, audit.ResourceQuery, audit.ActionExecute), h.ExecuteStoredAQLVersionPost)

	v1.Get("/definition/template/adl1.4", audit.Middleware(h.AuditLogger, audit.ResourceTemplate, audit.ActionRead), h.GetTemplatesADL14)
	v1.Post("/definition/template/adl1.4", audit.Middleware(h.AuditLogger, audit.ResourceTemplate, audit.ActionCreate), h.UploadTemplateADL14)
	v1.Get("/definition/template/adl1.4/:template_id", audit.Middleware(h.AuditLogger, audit.ResourceTemplate, audit.ActionRead), h.GetTemplateADL14ByID)

	v1.Get("/definition/template/adl2", audit.Middleware(h.AuditLogger, audit.ResourceTemplate, audit.ActionRead), h.GetTemplatesADL2)
	v1.Post("/definition/template/adl2", audit.Middleware(h.AuditLogger, audit.ResourceTemplate, audit.ActionCreate), h.UploadTemplateADL2)
	v1.Get("/definition/template/adl2/:template_id", audit.Middleware(h.AuditLogger, audit.ResourceTemplate, audit.ActionRead), h.GetTemplateADL2ByID)
	v1.Get("/definition/template/adl2/:template_id/:version", audit.Middleware(h.AuditLogger, audit.ResourceTemplate, audit.ActionRead), h.GetTemplateADL2AtVersion)

	v1.Get("/definition/query/:qualified_query_name", audit.Middleware(h.AuditLogger, audit.ResourceQuery, audit.ActionRead), h.ListStoredQueries)
	v1.Put("/definition/query/:qualified_query_name", audit.Middleware(h.AuditLogger, audit.ResourceQuery, audit.ActionCreate), h.StoreQuery)
	v1.Put("/definition/query/:qualified_query_name/:version", audit.Middleware(h.AuditLogger, audit.ResourceQuery, audit.ActionCreate), h.StoreQueryVersion)
	v1.Get("/definition/query/:qualified_query_name/:version", audit.Middleware(h.AuditLogger, audit.ResourceQuery, audit.ActionRead), h.GetStoredQueryAtVersion)

	v1.Delete("/admin/ehr/all", audit.Middleware(h.AuditLogger, audit.ResourceEHR, audit.ActionDelete), h.DeleteMultipleEHRs)
	v1.Delete("/admin/ehr/:ehr_id", audit.Middleware(h.AuditLogger, audit.ResourceEHR, audit.ActionDelete), h.DeleteEHRByID)
}

func (h *Handler) SystemInfo(c *fiber.Ctx) error {
	response := map[string]any{
		"solution":              h.Settings.Name,
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
	auditCtx := audit.From(c)

	subjectID := c.Query("subject_id")
	if subjectID == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "subject_id query parameters are required",
			Status:  "bad_request",
		})
	}

	subjectNamespace := c.Query("subject_namespace")
	if subjectNamespace == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "subject_namespace query parameters are required",
			Status:  "bad_request",
		})
	}

	ehr, err := h.OpenEHRService.GetEHRBySubject(ctx, subjectID, subjectNamespace)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given subject",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get EHR by subject", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get EHR by subject",
			Status:  "error",
		})
	}

	auditCtx.Success()

	return c.Status(fiber.StatusOK).JSON(ehr)
}

func (h *Handler) CreateEHR(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}
	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	ehrStatus := NewEHRStatus(uuid.New())
	if len(c.Body()) > 0 {
		if err := ParseBody(c, auditCtx, &ehrStatus); err != nil {
			return err
		}
	}

	ehrID := uuid.New()
	auditCtx.Event.Details["ehr_id"] = ehrID

	ehr, err := h.OpenEHRService.CreateEHR(ctx, ehrID, ehrStatus)
	if err != nil {
		if err == ErrEHRStatusAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "EHR status exists",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to create EHR", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create EHR",
			Status:  "error",
		})
	}

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeEHRCreated, map[string]any{
		"ehr_id": ehrID,
	})

	c.Set("ETag", `"`+ehrID.String()+`"`)
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String())

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(ehr)
	default:
		return c.JSON(ehr)
	}
}

func (h *Handler) GetEHR(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	ehr, err := h.OpenEHRService.GetEHR(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get EHR by ID", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	return c.Status(fiber.StatusOK).JSON(ehr)
}

func (h *Handler) CreateEHRWithID(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	ehrStatus := NewEHRStatus(uuid.New())
	if len(c.Body()) > 0 {
		if err := ParseBody(c, auditCtx, &ehrStatus); err != nil {
			return err
		}
	}

	ehr, err := h.OpenEHRService.CreateEHR(ctx, ehrID, ehrStatus)
	if err != nil {
		if err == ErrEHRAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "EHR with the given ID already exists",
				Status:  "conflict",
			})
		}
		if err == ErrEHRStatusAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "EHR Status with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to create EHR", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeEHRCreated, map[string]any{
		"ehr_id": ehrID,
	})

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
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(ehr)
	}
}

func (h *Handler) GetEHRStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	var filterAtTime time.Time
	if atTimeStr := c.Query("version_at_time"); atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid version_at_time format. Use RFC3339 format.",
				Status:  "bad_request",
			})
		}
		filterAtTime = parsedTime
	}

	ehrStatus, err := h.OpenEHRService.GetEHRStatus(ctx, ehrID, filterAtTime, "")
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR Status not found for the given EHR ID at the specified time",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get EHR Status at time", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	return c.Status(fiber.StatusOK).JSON(ehrStatus)
}

func (h *Handler) UpdateEhrStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	ifMatch, err := StringFromHeader(c, auditCtx, "If-Match")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	var ehrStatus model.EHR_STATUS
	if err := ParseBody(c, auditCtx, &ehrStatus); err != nil {
		return err
	}

	// Check collision using If-Match header
	currentEHRStatus, err := h.OpenEHRService.GetEHRStatus(ctx, ehrID, time.Time{}, "")
	if err != nil {
		if err == ErrEHRStatusNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR Status not found for the given EHR ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get current EHR Status", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}
	if currentEHRStatus.UID.V.ValueAsString() != ifMatch {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "EHR Status has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	updatedEHRStatus, err := h.OpenEHRService.UpdateEHRStatus(ctx, ehrID, ehrStatus)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR Status not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrEHRStatusAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "EHR Status with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err == ErrEHRStatusVersionLowerOrEqualToCurrent {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "EHR Status version in request body must be incremented",
				Status:  "bad_request",
			})
		}
		if err == ErrInvalidEHRStatusUIDMismatch {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "EHR Status UID HIER_OBJECT_ID in request body does not match current EHR Status UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to update EHR Status", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}
	updatedEHRStatusID := updatedEHRStatus.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeEHRStatusUpdated, map[string]any{
		"ehr_id":              ehrID,
		"prev_ehr_status_uid": currentEHRStatus.UID.V.ValueAsString(),
		"curr_ehr_status_uid": updatedEHRStatusID,
	})

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
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedEHRStatus)
	}
}

func (h *Handler) GetEHRStatusByVersionID(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version_uid parameter is required",
			Status:  "bad_request",
		})
	}

	ehrStatus, err := h.OpenEHRService.GetEHRStatus(ctx, ehrID, time.Time{}, versionUID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR Status not found for the given EHR ID and version UID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get EHR Status at time", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	return c.Status(fiber.StatusOK).JSON(ehrStatus)
}

func (h *Handler) GetVersionedEHRStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	versionedStatus, err := h.OpenEHRService.GetVersionedEHRStatus(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned EHR Status not found for the given EHR ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	return c.Status(fiber.StatusOK).JSON(versionedStatus)
}

func (h *Handler) GetVersionedEHRStatusRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	revisionHistory, err := h.OpenEHRService.GetVersionedEHRStatusRevisionHistory(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned EHR Status revision history not found for the given EHR ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status revision history", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	return c.Status(fiber.StatusOK).JSON(revisionHistory)
}

func (h *Handler) GetVersionedEHRStatusVersion(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	var filterAtTime time.Time
	atTimeStr := c.Query("version_at_time")
	if atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid version_at_time format. Use RFC3339 format.",
				Status:  "bad_request",
			})
		}
		filterAtTime = parsedTime
	}

	ehrStatusVersionJSON, err := h.OpenEHRService.GetVersionedEHRStatusVersionAsJSON(ctx, ehrID, filterAtTime, "")
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned EHR Status version not found for the given EHR ID at the specified time",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status version at time", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(ehrStatusVersionJSON)
}

func (h *Handler) GetVersionedEHRStatusVersionByID(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version_uid parameter is required",
			Status:  "bad_request",
		})
	}

	ehrStatusVersionJSON, err := h.OpenEHRService.GetVersionedEHRStatusVersionAsJSON(ctx, ehrID, time.Time{}, versionUID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned EHR Status version not found for the given EHR ID and version UID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status version by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(ehrStatusVersionJSON)
}

func (h *Handler) CreateComposition(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	var requestComposition model.COMPOSITION
	if err := ParseBody(c, auditCtx, &requestComposition); err != nil {
		return err
	}

	exists, err := h.OpenEHRService.ExistsEHR(ctx, ehrID)
	if err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to check if EHR exists", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}
	if !exists {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotFound,
			Message: "EHR not found for the given EHR ID",
			Status:  "not_found",
		})
	}

	composition, err := h.OpenEHRService.CreateComposition(ctx, ehrID, requestComposition)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Composition with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to create Composition", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeCompositionCreated, map[string]any{
		"ehr_id":          ehrID,
		"composition_uid": composition.UID.V.ValueAsString(),
	})

	compositionID := composition.UID.V.ValueAsString()
	c.Set("ETag", "\""+compositionID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String()+"/composition/"+compositionID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(composition)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + compositionID + `"}`)
	default:
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(composition)
	}
}

func (h *Handler) GetComposition(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	composition, err := h.OpenEHRService.GetComposition(ctx, ehrID, uidBasedID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Composition not found for the given UID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Composition by ID", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	return c.Status(fiber.StatusOK).JSON(composition)
}

func (h *Handler) UpdateComposition(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["uid_based_id"] = uidBasedID

	ifMatch, err := StringFromHeader(c, auditCtx, "If-Match")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	versionedCompositionID := strings.Split(uidBasedID, "::")[0]
	currentComposition, err := h.OpenEHRService.GetComposition(ctx, ehrID, versionedCompositionID)
	if err != nil {
		if err == ErrCompositionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Composition not found for the given EHR ID and UID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get current Composition", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}
	if currentComposition.UID.V.ValueAsString() != ifMatch {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Composition has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	var composition model.COMPOSITION
	if err := ParseBody(c, auditCtx, &composition); err != nil {
		return err
	}

	updatedComposition, err := h.OpenEHRService.UpdateComposition(ctx, ehrID, composition)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Composition not found for the given UID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Composition with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err == ErrCompositionVersionLowerOrEqualToCurrent {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Composition version in request body is lower or equal to the current version",
				Status:  "bad_request",
			})
		}
		if err == ErrInvalidCompositionUIDMismatch {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Composition UID HIER_OBJECT_ID in request body does not match current Composition UID",
				Status:  "bad_request",
			})
		}
		if err == ErrCompositionUIDNotProvided {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Composition UID must be provided in the request body for update",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to update Composition", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}
	updatedCompositionID := updatedComposition.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeCompositionUpdated, map[string]any{
		"ehr_id":               ehrID,
		"prev_composition_uid": currentComposition.UID.V.ValueAsString(),
		"curr_composition_uid": updatedCompositionID,
	})

	c.Set("ETag", "\""+updatedCompositionID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String()+"/composition/"+updatedCompositionID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedComposition)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + updatedCompositionID + `"}`)
	default:
		h.Telemetry.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedComposition)
	}
}

func (h *Handler) DeleteComposition(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["uid_based_id"] = uidBasedID

	versionedCompositionID := strings.Split(uidBasedID, "::")[0]
	currentComposition, err := h.OpenEHRService.GetComposition(ctx, ehrID, versionedCompositionID)
	if err != nil {
		if err == ErrCompositionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Composition not found for the given UID and EHR ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Composition by ID before deletion", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}
	if currentComposition.UID.V.ValueAsString() != uidBasedID {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Composition has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	if err := h.OpenEHRService.DeleteComposition(ctx, ehrID, currentComposition.UID.V.UUID()); err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Composition not found for the given UID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to delete Composition by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeCompositionDeleted, map[string]any{
		"ehr_id":          ehrID,
		"composition_uid": currentComposition.UID.V.ValueAsString(),
	})

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) GetVersionedCompositionByID(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	versionedCompositionID, err := StringFromPath(c, auditCtx, "versioned_object_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["versioned_composition_uid"] = versionedCompositionID

	versionedComposition, err := h.OpenEHRService.GetVersionedComposition(ctx, ehrID, versionedCompositionID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Composition not found for the given versioned object ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned Composition by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(versionedComposition)
}

func (h *Handler) GetVersionedCompositionRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	versionedCompositionID, err := StringFromPath(c, auditCtx, "versioned_object_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["versioned_composition_uid"] = versionedCompositionID

	revisionHistory, err := h.OpenEHRService.GetVersionedCompositionRevisionHistory(ctx, ehrID, versionedCompositionID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Composition revision history not found for the given versioned object ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned Composition revision history", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(revisionHistory)
}

func (h *Handler) GetVersionedCompositionVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	versionedCompositionID, err := StringFromPath(c, auditCtx, "versioned_object_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["versioned_composition_uid"] = versionedCompositionID

	var filterAtTime time.Time
	if c.Query("version_at_time") != "" {
		versionAtTime, err := TimeFromQuery(c, auditCtx, "version_at_time")
		if err != nil {
			return err
		}
		filterAtTime = versionAtTime
		auditCtx.Event.Details["version_at_time"] = filterAtTime.String()
	}

	versionJSON, err := h.OpenEHRService.GetVersionedCompositionVersionJSON(ctx, ehrID, versionedCompositionID, filterAtTime, "")
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Composition version not found for the given versioned object ID at the specified time",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned Composition version at time", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(versionJSON)
}

func (h *Handler) GetVersionedCompositionVersionByID(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	versionedCompositionID, err := StringFromPath(c, auditCtx, "versioned_object_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["versioned_composition_uid"] = versionedCompositionID

	versionID, err := StringFromPath(c, auditCtx, "version_uid")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["version_uid"] = versionID

	if versionedCompositionID != strings.Split(versionID, "::")[0] {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version_uid does not match the versioned_object_id",
			Status:  "bad_request",
		})
	}

	versionJSON, err := h.OpenEHRService.GetVersionedCompositionVersionJSON(ctx, ehrID, versionedCompositionID, time.Time{}, versionID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrCompositionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Composition version not found for the given versioned object ID and version UID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned Composition version by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}

	auditCtx.Success()
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(versionJSON)
}

func (h *Handler) CreateDirectory(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	var requestDirectory model.FOLDER
	if err := ParseBody(c, auditCtx, &requestDirectory); err != nil {
		return err
	}

	directory, err := h.OpenEHRService.CreateDirectory(ctx, ehrID, requestDirectory)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Directory with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to create Directory", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Status:  "error",
		})
	}
	directoryID := directory.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeDirectoryCreated, map[string]any{
		"directory_id": directoryID,
	})

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
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(directory)
	}
}

func (h *Handler) UpdateDirectory(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	ifMatch, err := StringFromHeader(c, auditCtx, "If-Match")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	var directory model.FOLDER
	if err := ParseBody(c, auditCtx, &directory); err != nil {
		return err
	}

	currentDirectory, err := h.OpenEHRService.GetDirectory(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get current Directory", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Directory",
			Status:  "error",
		})
	}
	if currentDirectory.UID.V.ValueAsString() != ifMatch {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Directory has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	updatedDirectory, err := h.OpenEHRService.UpdateDirectory(ctx, ehrID, directory)
	if err != nil {
		if err == ErrDirectoryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Directory with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err == ErrDirectoryVersionLowerOrEqualToCurrent {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Directory version in request body is lower or equal to the current version",
				Status:  "bad_request",
			})
		}
		if err == ErrInvalidDirectoryUIDMismatch {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Directory UID HIER_OBJECT_ID in request body does not match current Directory UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to update Directory", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Directory",
			Status:  "error",
		})
	}
	updatedDirectoryID := updatedDirectory.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeDirectoryUpdated, map[string]any{
		"prev_directory_id": currentDirectory.UID.V.ValueAsString(),
		"curr_directory_id": updatedDirectoryID,
	})

	c.Set("ETag", "\""+updatedDirectoryID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/ehr/"+ehrID.String()+"/directory/"+updatedDirectoryID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedDirectory)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusOK).JSON(`{"uid":"` + updatedDirectoryID + `"}`)
	default:
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedDirectory)
	}
}

func (h *Handler) DeleteDirectory(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	ifMatch, err := StringFromHeader(c, auditCtx, "If-Match")
	if err != nil {
		return err
	}

	currentDirectory, err := h.OpenEHRService.GetDirectory(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get current Directory", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Directory",
			Status:  "error",
		})
	}
	if currentDirectory.UID.V.ValueAsString() != ifMatch {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Directory has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	if err := h.OpenEHRService.DeleteDirectory(ctx, ehrID, uuid.MustParse(currentDirectory.UID.V.Value.(*model.OBJECT_VERSION_ID).Value)); err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to delete Directory", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Directory",
			Status:  "error",
		})
	}

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeDirectoryDeleted, map[string]any{
		"directory_id": currentDirectory.UID.V.ValueAsString(),
	})

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) GetFolderInDirectoryVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	var filterAtTime time.Time
	atTimeStr := c.Query("version_at_time")
	if atTimeStr != "" {
		parsedTime, err := TimeFromQuery(c, auditCtx, "version_at_time")
		if err != nil {
			return err
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

	folder, err := h.OpenEHRService.GetFolderInDirectoryVersion(ctx, ehrID, filterAtTime, "", pathParts)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID at the specified time",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory version not found at the specified time for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrFolderNotFoundInDirectory {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Folder not found in Directory version at the specified time for the given path",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Folder in Directory version at time", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Folder in Directory version at time",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(folder)
}

func (h *Handler) GetFolderInDirectoryVersion(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	versionUID, err := StringFromPath(c, auditCtx, "version_uid")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["version_uid"] = versionUID

	path := c.Query("path")
	var pathParts []string
	if path != "" {
		for part := range strings.SplitSeq(path, "/") {
			pathParts = append(pathParts, strings.TrimSpace(part))
		}
	}

	folder, err := h.OpenEHRService.GetFolderInDirectoryVersion(ctx, ehrID, time.Time{}, versionUID, pathParts)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory not found for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrDirectoryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Directory version not found at the specified time for the given EHR ID",
				Status:  "not_found",
			})
		}
		if err == ErrFolderNotFoundInDirectory {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Folder not found in Directory version at the specified time for the given path",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Folder in Directory version", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Folder in Directory version",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(folder)
}

func (h *Handler) CreateContribution(c *fiber.Ctx) error {
	c.Status(fiber.StatusNotImplemented)
	return nil
}

func (h *Handler) GetContribution(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	contributionUID, err := StringFromPath(c, auditCtx, "contribution_uid")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["contribution_uid"] = contributionUID

	contribution, err := h.OpenEHRService.GetContribution(ctx, contributionUID, utils.Some(ehrID))
	if err != nil {
		if err == ErrContributionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Contribution not found for the given EHR ID and Contribution UID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Contribution by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Contribution by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()
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
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	var agent model.AGENT
	if err := ParseBody(c, auditCtx, &agent); err != nil {
		return err
	}

	createdAgent, err := h.OpenEHRService.CreateAgent(ctx, agent)
	if err != nil {
		if err == ErrAgentAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Agent with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to create Agent", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create Agent",
			Status:  "error",
		})
	}
	createdAgentID := createdAgent.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeAgentCreated, map[string]any{
		"agent_uid": createdAgentID,
	})

	c.Set("ETag", "\""+createdAgentID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/agent/"+createdAgentID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(createdAgent)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + createdAgentID + `"}`)
	default:
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(createdAgent)
	}
}

func (h *Handler) GetAgent(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["uid_based_id"] = uidBasedID

	if strings.Count(uidBasedID, "::") == 2 {
		agent, err := h.OpenEHRService.GetAgentAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == ErrAgentNotFound {
				return SendErrorResponse(c, auditCtx, ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Agent not found for the given agent ID",
					Status:  "not_found",
				})
			}

			h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)

			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to get Agent by ID",
				Status:  "error",
			})
		}

		auditCtx.Success()
		return c.Status(fiber.StatusOK).JSON(agent)
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	agent, err := h.OpenEHRService.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrAgentNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Agent not found for the given agent ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Agent by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(agent)
}

func (h *Handler) UpdateAgent(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["uid_based_id"] = uidBasedID

	ifMatch, err := StringFromHeader(c, auditCtx, "If-Match")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	var agent model.AGENT
	if err := ParseBody(c, auditCtx, &agent); err != nil {
		return err
	}

	currentAgent, err := h.OpenEHRService.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrAgentNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Agent not found for the given agent ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get current Agent by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Agent by ID",
			Status:  "error",
		})
	}
	if currentAgent.UID.V.ValueAsString() != ifMatch {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Agent has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	updatedAgent, err := h.OpenEHRService.UpdateAgent(ctx, versionedPartyID, agent)
	if err != nil {
		if err == ErrAgentNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Agent not found for the given agent ID",
				Status:  "not_found",
			})
		}
		if err == ErrAgentVersionLowerOrEqualToCurrent {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Agent version is lower than or equal to the current version",
				Status:  "conflict",
			})
		}
		if err == ErrInvalidAgentUIDMismatch {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Agent UID HIER_OBJECT_ID in request body does not match current Agent UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to update Agent", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Agent",
			Status:  "error",
		})
	}
	updatedAgentID := updatedAgent.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeAgentUpdated, map[string]any{
		"prev_agent_uid": currentAgent.UID.V.ValueAsString(),
		"curr_agent_uid": updatedAgentID,
	})

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
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedAgent)
	}
}

func (h *Handler) DeleteAgent(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["uid_based_id"] = uidBasedID

	if strings.Count(uidBasedID, "::") != 2 {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Agent by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	currentAgent, err := h.OpenEHRService.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrAgentNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Agent not found for the given agent ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Agent by ID before deletion", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Agent by ID before deletion",
			Status:  "error",
		})
	}
	if currentAgent.UID.V.ValueAsString() != uidBasedID {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Agent has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	if err := h.OpenEHRService.DeleteAgent(ctx, uidBasedID); err != nil {
		if err == ErrAgentNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Agent not found for the given agent ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to delete Agent by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Agent by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeAgentDeleted, map[string]any{
		"versioned_party_id": versionedPartyID,
	})

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) CreateGroup(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	var group model.GROUP
	if err := ParseBody(c, auditCtx, &group); err != nil {
		return err
	}

	createdGroup, err := h.OpenEHRService.CreateGroup(ctx, group)
	if err != nil {
		if err == ErrGroupAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Group with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to create Group", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create Group",
			Status:  "error",
		})
	}
	createdGroupID := createdGroup.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeGroupCreated, map[string]any{
		"group_uid": createdGroupID,
	})

	c.Set("ETag", "\""+createdGroupID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/group/"+createdGroupID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(createdGroup)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + createdGroupID + `"}`)
	default:
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(createdGroup)
	}
}

func (h *Handler) GetGroup(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["uid_based_id"] = uidBasedID

	if strings.Count(uidBasedID, "::") == 2 {
		group, err := h.OpenEHRService.GetGroupAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == ErrGroupNotFound {
				return SendErrorResponse(c, auditCtx, ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Group not found for the given group ID",
					Status:  "not_found",
				})
			}

			h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Group by ID", "error", err)

			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to get Group by ID",
				Status:  "error",
			})
		}

		auditCtx.Success()
		return c.Status(fiber.StatusOK).JSON(group)
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	group, err := h.OpenEHRService.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrGroupNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Group not found for the given group ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Group by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Group by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(group)
}

func (h *Handler) UpdateGroup(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	ifMatch, err := StringFromHeader(c, auditCtx, "If-Match")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	var group model.GROUP
	if err := ParseBody(c, auditCtx, &group); err != nil {
		return err
	}

	currentGroup, err := h.OpenEHRService.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrGroupNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Group not found for the given group ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get current Group by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Group by ID",
			Status:  "error",
		})
	}
	if currentGroup.UID.V.ValueAsString() != ifMatch {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Group has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	updatedGroup, err := h.OpenEHRService.UpdateGroup(ctx, versionedPartyID, group)
	if err != nil {
		if err == ErrGroupNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Group not found for the given group ID",
				Status:  "not_found",
			})
		}
		if err == ErrGroupVersionLowerOrEqualToCurrent {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Group version is lower than or equal to the current version",
				Status:  "conflict",
			})
		}
		if err == ErrInvalidGroupUIDMismatch {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Group UID HIER_OBJECT_ID in request body does not match current Group UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to update Agent", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Agent",
			Status:  "error",
		})
	}
	updatedGroupID := updatedGroup.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeGroupUpdated, map[string]any{
		"prev_group_uid": currentGroup.UID.V.ValueAsString(),
		"curr_group_uid": updatedGroupID,
	})

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
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedGroup)
	}
}

func (h *Handler) DeleteGroup(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	versionUID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	if strings.Count(versionUID, "::") != 2 {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Group by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	versionedPartyID, err := uuid.Parse(strings.Split(versionUID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	currentGroup, err := h.OpenEHRService.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrGroupNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Group not found for the given group ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Group by ID before deletion", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Group by ID before deletion",
			Status:  "error",
		})
	}
	if currentGroup.UID.V.ValueAsString() != versionUID {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Group has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	if err := h.OpenEHRService.DeleteGroup(ctx, versionedPartyID); err != nil {
		if err == ErrGroupNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Group not found for the given group ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to delete Group by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Group by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeGroupDeleted, map[string]any{
		"versioned_party_id": versionedPartyID,
	})

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) CreatePerson(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	var person model.PERSON
	if err := ParseBody(c, auditCtx, &person); err != nil {
		return err
	}

	createdPerson, err := h.OpenEHRService.CreatePerson(ctx, person)
	if err != nil {
		if err == ErrPersonAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Person with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to create Person", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create Person",
			Status:  "error",
		})
	}
	createdPersonID := createdPerson.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypePersonCreated, map[string]any{
		"person_uid": createdPersonID,
	})

	c.Set("ETag", "\""+createdPersonID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/person/"+createdPersonID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(createdPerson)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + createdPersonID + `"}`)
	default:
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(createdPerson)
	}
}

func (h *Handler) GetPerson(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	if strings.Count(uidBasedID, "::") == 2 {
		person, err := h.OpenEHRService.GetPersonAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == ErrPersonNotFound {
				return SendErrorResponse(c, auditCtx, ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Person not found for the given person ID",
					Status:  "not_found",
				})
			}

			h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)

			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to get Person",
				Status:  "error",
			})
		}

		auditCtx.Success()
		return c.Status(fiber.StatusOK).JSON(person)
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	person, err := h.OpenEHRService.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrPersonNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Person not found for the given person ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Person",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(person)
}

func (h *Handler) UpdatePerson(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	ifMatch, err := StringFromHeader(c, auditCtx, "If-Match")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	var person model.PERSON
	if err := ParseBody(c, auditCtx, &person); err != nil {
		return err
	}

	currentPerson, err := h.OpenEHRService.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrPersonNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Person not found for the given person ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get current Person by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Person",
			Status:  "error",
		})
	}
	if currentPerson.UID.V.ValueAsString() != ifMatch {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Person has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	updatedPerson, err := h.OpenEHRService.UpdatePerson(ctx, versionedPartyID, person)
	if err != nil {
		if err == ErrPersonNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Person not found for the given person ID",
				Status:  "not_found",
			})
		}
		if err == ErrPersonVersionLowerOrEqualToCurrent {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Person version is lower than or equal to the current version",
				Status:  "conflict",
			})
		}
		if err == ErrInvalidPersonUIDMismatch {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Person UID HIER_OBJECT_ID in request body does not match current Person UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to update Person", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Person",
			Status:  "error",
		})
	}
	updatedPersonID := updatedPerson.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypePersonUpdated, map[string]any{
		"prev_person_uid": currentPerson.UID.V.ValueAsString(),
		"curr_person_uid": updatedPersonID,
	})

	c.Set("ETag", "\""+updatedPersonID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/person/"+updatedPersonID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedPerson)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusOK).JSON(`{"uid":"` + updatedPersonID + `"}`)
	default:
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedPerson)
	}
}

func (h *Handler) DeletePerson(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	versionID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	if strings.Count(versionID, "::") != 2 {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Person by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}
	currentPerson, err := h.OpenEHRService.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrPersonNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Person not found for the given person ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Person by ID before deletion", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Person by ID before deletion",
			Status:  "error",
		})
	}
	if currentPerson.UID.V.ValueAsString() != versionID {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Person has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	if err := h.OpenEHRService.DeletePerson(ctx, versionedPartyID); err != nil {
		if err == ErrPersonNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Person not found for the given person ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to delete Person by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Person by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypePersonDeleted, map[string]any{
		"versioned_party_id": versionedPartyID,
	})

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) CreateOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	var organisation model.ORGANISATION
	if err := ParseBody(c, auditCtx, &organisation); err != nil {
		return err
	}

	createdOrganisation, err := h.OpenEHRService.CreateOrganisation(ctx, organisation)
	if err != nil {
		if err == ErrOrganisationAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Organisation with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: err,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to create Organisation", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create Organisation",
			Status:  "error",
		})
	}
	createdOrganisationID := createdOrganisation.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeOrganisationCreated, map[string]any{
		"organisation_uid": createdOrganisationID,
	})

	c.Set("ETag", "\""+createdOrganisationID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/organisation/"+createdOrganisationID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(organisation)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + createdOrganisationID + `"}`)
	default:
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(organisation)
	}
}

func (h *Handler) GetOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	if strings.Count(uidBasedID, "::") == 2 {
		organisation, err := h.OpenEHRService.GetOrganisationAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == ErrOrganisationNotFound {
				return SendErrorResponse(c, auditCtx, ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Organisation not found for the given organisation ID",
					Status:  "not_found",
				})
			}

			h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Organisation by ID", "error", err)

			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to get Organisation by ID",
				Status:  "error",
			})
		}

		auditCtx.Success()
		return c.Status(fiber.StatusOK).JSON(organisation)
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	organisation, err := h.OpenEHRService.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrOrganisationNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Organisation not found for the given organisation ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Organisation by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Organisation by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()

	return c.Status(fiber.StatusOK).JSON(organisation)
}

func (h *Handler) UpdateOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	versionID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	ifMatch, err := StringFromHeader(c, auditCtx, "If-Match")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	var organisation model.ORGANISATION
	if err := ParseBody(c, auditCtx, &organisation); err != nil {
		return err
	}

	currentOrganisation, err := h.OpenEHRService.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrOrganisationNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Organisation not found for the given organisation ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get current Organisation by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Organisation by ID",
			Status:  "error",
		})
	}
	if currentOrganisation.UID.V.ValueAsString() != ifMatch {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Organisation has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	updatedOrganisation, err := h.OpenEHRService.UpdateOrganisation(ctx, versionedPartyID, organisation)
	if err != nil {
		if err == ErrOrganisationNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Organisation not found for the given organisation ID",
				Status:  "not_found",
			})
		}
		if err == ErrOrganisationVersionLowerOrEqualToCurrent {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Organisation version is lower than or equal to the current version",
				Status:  "conflict",
			})
		}
		if err == ErrInvalidOrganisationUIDMismatch {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Organisation with the given UID already exists",
				Status:  "conflict",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: validationErrs,
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to update Organisation", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Organisation",
			Status:  "error",
		})
	}
	updatedOrganisationID := updatedOrganisation.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeOrganisationUpdated, map[string]any{
		"prev_organisation_uid": currentOrganisation.UID.V.ValueAsString(),
		"curr_organisation_uid": updatedOrganisationID,
	})

	c.Set("ETag", "\""+updatedOrganisationID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/organisation/"+updatedOrganisationID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedOrganisation)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusOK).JSON(`{"uid":"` + updatedOrganisationID + `"}`)
	default:
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedOrganisation)
	}
}

func (h *Handler) DeleteOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	if strings.Count(uidBasedID, "::") != 2 {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Organisation by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	currentOrganisation, err := h.OpenEHRService.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrOrganisationNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Organisation not found for the given organisation ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Organisation by ID before deletion", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Organisation by ID before deletion",
			Status:  "error",
		})
	}
	if currentOrganisation.UID.V.ValueAsString() != uidBasedID {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Organisation has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	if err := h.OpenEHRService.DeleteOrganisation(ctx, versionedPartyID); err != nil {
		if err == ErrOrganisationNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Organisation not found for the given organisation ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to delete Organisation by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Organisation by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeOrganisationDeleted, map[string]any{
		"versioned_party_id": versionedPartyID,
	})

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) CreateRole(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	var role model.ROLE
	if err := ParseBody(c, auditCtx, &role); err != nil {
		return err
	}

	createdRole, err := h.OpenEHRService.CreateRole(ctx, role)
	if err != nil {
		if err == ErrRoleAlreadyExists {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Role with the given UID already exists",
				Status:  "conflict",
			})
		}
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to create Role", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create Role",
			Status:  "error",
		})
	}
	createdRoleID := createdRole.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeRoleCreated, map[string]any{
		"role_uid": createdRoleID,
	})

	c.Set("ETag", "\""+createdRoleID+"\"")
	c.Set("Location", c.Protocol()+"://"+c.Hostname()+"/openehr/v1/role/"+createdRoleID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusCreated)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusCreated).JSON(role)
	case ReturnTypeIdentifier:
		return c.Status(fiber.StatusCreated).JSON(`{"uid":"` + createdRoleID + `"}`)
	default:
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusCreated).JSON(role)
	}
}

func (h *Handler) GetRole(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	if strings.Count(uidBasedID, "::") == 2 {
		role, err := h.OpenEHRService.GetRoleAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == ErrRoleNotFound {
				return SendErrorResponse(c, auditCtx, ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Role not found for the given role ID",
					Status:  "not_found",
				})
			}

			h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Role by ID", "error", err)

			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to get Role by ID",
				Status:  "error",
			})
		}

		auditCtx.Success()
		return c.Status(fiber.StatusOK).JSON(role)
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	role, err := h.OpenEHRService.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrRoleNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Role not found for the given role ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Role by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Role by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(role)
}

func (h *Handler) UpdateRole(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	versionID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	ifMatch, err := StringFromHeader(c, auditCtx, "If-Match")
	if err != nil {
		return err
	}

	returnType, err := ReturnTypeFromHeader(c, auditCtx)
	if err != nil {
		return err
	}

	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	var role model.ROLE
	if err := ParseBody(c, auditCtx, &role); err != nil {
		return err
	}

	currentRole, err := h.OpenEHRService.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrRoleNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Role not found for the given role ID",
				Status:  "not_found",
			})
		}
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get current Role by ID", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get current Role by ID",
			Status:  "error",
		})
	}

	if currentRole.UID.V.ValueAsString() != ifMatch {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Role has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	updatedRole, err := h.OpenEHRService.UpdateRole(ctx, versionedPartyID, role)
	if err != nil {
		if err == ErrRoleNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Role not found for the given role ID",
				Status:  "not_found",
			})
		}
		if err == ErrRoleVersionLowerOrEqualToCurrent {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusConflict,
				Message: "Role version is lower than or equal to the current version",
				Status:  "conflict",
			})
		}
		if err == ErrInvalidRoleUIDMismatch {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Role UID HIER_OBJECT_ID in request body does not match current Role UID",
				Status:  "bad_request",
			})
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to update Role", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update Role",
			Status:  "error",
		})
	}
	updatedRoleID := updatedRole.UID.V.ValueAsString()

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeRoleUpdated, map[string]any{
		"prev_role_uid": currentRole.UID.V.ValueAsString(),
		"curr_role_uid": updatedRoleID,
	})

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
		h.Telemetry.Logger.WarnContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedRole)
	}
}

func (h *Handler) DeleteRole(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	uidBasedID, err := StringFromPath(c, auditCtx, "uid_based_id")
	if err != nil {
		return err
	}

	if !strings.Contains(uidBasedID, "::") {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Cannot delete Role by versioned object ID. Please provide the object version ID.",
			Status:  "bad_request",
		})
	}

	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid uid_based_id format",
			Status:  "bad_request",
		})
	}

	currentRole, err := h.OpenEHRService.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == ErrRoleNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Role not found for the given role ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Role by ID before deletion", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Role by ID before deletion",
			Status:  "error",
		})
	}
	if currentRole.UID.V.ValueAsString() != uidBasedID {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusPreconditionFailed,
			Message: "Role has been modified since the provided version",
			Status:  "precondition_failed",
		})
	}

	if err := h.OpenEHRService.DeleteRole(ctx, versionedPartyID); err != nil {
		if err == ErrRoleNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Role not found for the given role ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to delete Role by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete Role by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()

	h.WebhookSaver.Enqueue(webhook.EventTypeRoleDeleted, map[string]any{
		"versioned_party_id": versionedPartyID,
	})

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) GetVersionedParty(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	versionedObjectID, err := UUIDFromPath(c, auditCtx, "versioned_object_id")
	if err != nil {
		return err
	}

	party, err := h.OpenEHRService.GetVersionedParty(ctx, versionedObjectID)
	if err != nil {
		if err == ErrVersionedPartyNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party not found for the given ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned Party by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Versioned Party by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(party)
}

func (h *Handler) GetVersionedPartyRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	versionedObjectID, err := UUIDFromPath(c, auditCtx, "versioned_object_id")
	if err != nil {
		return err
	}

	revisionHistory, err := h.OpenEHRService.GetVersionedPartyRevisionHistory(ctx, versionedObjectID)
	if err != nil {
		if err == ErrRevisionHistoryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party not found for the given ID",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned Party Revision History by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Versioned Party Revision History by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(revisionHistory)
}

func (h *Handler) GetVersionedPartyVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	versionedObjectID, err := UUIDFromPath(c, auditCtx, "versioned_object_id")
	if err != nil {
		return err
	}

	var filterAtTime time.Time
	if c.Query("version_at_time") != "" {
		versionAtTime, err := TimeFromQuery(c, auditCtx, "version_at_time")
		if err != nil {
			return err
		}
		filterAtTime = versionAtTime
	}

	partyVersionJSON, err := h.OpenEHRService.GetVersionedPartyVersionJSON(ctx, versionedObjectID, filterAtTime, "")
	if err != nil {
		if err == ErrVersionedPartyNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party not found for the given ID",
				Status:  "not_found",
			})
		}
		if err == ErrVersionedPartyVersionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party version not found for the given time",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned Party Version at Time by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Versioned Party Version at Time by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(partyVersionJSON)
}

func (h *Handler) GetVersionedPartyVersion(c *fiber.Ctx) error {
	ctx := c.Context()
	auditCtx := audit.From(c)

	versionedObjectID, err := UUIDFromPath(c, auditCtx, "versioned_object_id")
	if err != nil {
		return err
	}

	versionID, err := StringFromPath(c, auditCtx, "version_id")
	if err != nil {
		return err
	}

	partyVersionJSON, err := h.OpenEHRService.GetVersionedPartyVersionJSON(ctx, versionedObjectID, time.Time{}, versionID)
	if err != nil {
		if err == ErrVersionedPartyNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party not found for the given ID",
				Status:  "not_found",
			})
		}
		if err == ErrVersionedPartyVersionNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Versioned Party version not found for the given version ID",
				Status:  "not_found",
			})
		}
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get Versioned Party Version by ID", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get Versioned Party Version by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(partyVersionJSON)
}

func (h *Handler) CreateDemographicContribution(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Create Demographic Contribution not implemented yet")
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

	auditCtx := audit.From(c)

	query := c.Query("q")
	if query == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "q query parameter is required",
			Status:  "bad_request",
		})
	}

	ehrID := c.Query("ehr_id")
	if ehrID != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	fetch := c.Query("fetch")
	if fetch != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	offset := c.Query("offset")
	if offset != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	queryParametersStr := c.Query("query_parameters")
	queryParametersURLValues, err := url.ParseQuery(queryParametersStr)
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
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

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), query, queryParameters); err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to execute Ad Hoc AQL", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Ad Hoc AQL",
			Status:  "error",
		})
	}

	auditCtx.Success()

	// Register event
	h.WebhookSaver.Enqueue(webhook.EventTypeQueryExecuted, map[string]any{
		"query":      query,
		"parameters": queryParameters,
	})

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

	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	var aqlRequest AdHocAQLRequest
	if err := ParseBody(c, auditCtx, &aqlRequest); err != nil {
		return err
	}

	if aqlRequest.Query == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "query field is required in the request body",
			Status:  "bad_request",
		})
	}

	if aqlRequest.EHRID != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL Post with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Fetch != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL Post with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Offset != 0 {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Ad Hoc AQL Post with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), aqlRequest.Query, aqlRequest.QueryParameters); err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to execute Ad Hoc AQL Post", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Ad Hoc AQL Post",
			Status:  "error",
		})
	}

	auditCtx.Success()

	// Register event
	h.WebhookSaver.Enqueue(webhook.EventTypeQueryExecuted, map[string]any{
		"query":      aqlRequest.Query,
		"parameters": aqlRequest.QueryParameters,
	})

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) ExecuteStoredAQL(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	ehrID := c.Query("ehr_id")
	if ehrID != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	fetch := c.Query("fetch")
	if fetch != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	offset := c.Query("offset")
	if offset != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	queryParametersStr := c.Query("query_parameters")
	queryParametersURLValues, err := url.ParseQuery(queryParametersStr)
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
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

	// Retrieve stored query by name
	storedQuery, err := h.OpenEHRService.GetQueryByName(ctx, name, "")
	if err != nil {
		if err == ErrQueryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Stored query not found for the given name",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get stored query by name", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get stored query by name",
			Status:  "error",
		})
	}

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), storedQuery.Query, nil); err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to execute Stored AQL", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Stored AQL",
			Status:  "error",
		})
	}

	auditCtx.Success()

	// Register event
	h.WebhookSaver.Enqueue(webhook.EventTypeQueryExecuted, map[string]any{
		"query_name": name,
		"parameters": queryParameters,
	})

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

	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	var aqlRequest StoredAQLRequest
	if err := ParseBody(c, auditCtx, &aqlRequest); err != nil {
		return err
	}

	if aqlRequest.EHRID != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Post with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Fetch != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Post with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Offset != 0 {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Post with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	// Retrieve stored query by name
	storedQuery, err := h.OpenEHRService.GetQueryByName(ctx, name, "")
	if err != nil {
		if err == ErrQueryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Stored query not found for the given name",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get stored query by name", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get stored query by name",
			Status:  "error",
		})
	}

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), storedQuery.Query, aqlRequest.QueryParameters); err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to execute Stored AQL Post", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Stored AQL Post",
			Status:  "error",
		})
	}

	auditCtx.Success()

	// Register event
	h.WebhookSaver.Enqueue(webhook.EventTypeQueryExecuted, map[string]any{
		"query_name": name,
		"parameters": aqlRequest.QueryParameters,
	})

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) ExecuteStoredAQLVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	version := c.Params("version")
	if version == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version path parameter is required",
			Status:  "bad_request",
		})
	}

	ehrID := c.Query("ehr_id")
	if ehrID != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	fetch := c.Query("fetch")
	if fetch != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	offset := c.Query("offset")
	if offset != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	queryParametersStr := c.Query("query_parameters")
	queryParametersURLValues, err := url.ParseQuery(queryParametersStr)
	if err != nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
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

	// Retrieve stored query by name and version
	storedQuery, err := h.OpenEHRService.GetQueryByName(ctx, name, version)
	if err != nil {
		if err == ErrQueryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Stored query not found for the given name and version",
				Status:  "not_found",
			})
		}
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get stored query by name and version", "error", err)
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get stored query by name and version",
			Status:  "error",
		})
	}

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), storedQuery.Query, queryParameters); err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to execute Stored AQL Version", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Stored AQL Version",
			Status:  "error",
		})
	}

	auditCtx.Success()

	// Register event
	h.WebhookSaver.Enqueue(webhook.EventTypeQueryExecuted, map[string]any{
		"query_name": name,
		"version":    version,
		"parameters": queryParameters,
	})

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

	auditCtx := audit.From(c)

	err := Accepts(c, auditCtx, "application/json")
	if err != nil {
		return err
	}

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	version := c.Params("version")
	if version == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version path parameter is required",
			Status:  "bad_request",
		})
	}

	var aqlRequest StoredAQLVersionRequest
	if err := ParseBody(c, auditCtx, &aqlRequest); err != nil {
		return err
	}

	if aqlRequest.EHRID != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version Post with ehr_id not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Fetch != "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version Post with fetch not implemented yet",
			Status:  "not_implemented",
		})
	}

	if aqlRequest.Offset != 0 {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotImplemented,
			Message: "Execute Stored AQL Version Post with offset not implemented yet",
			Status:  "not_implemented",
		})
	}

	// Retrieve stored query by name and version
	storedQuery, err := h.OpenEHRService.GetQueryByName(ctx, name, version)
	if err != nil {
		if err == ErrQueryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Stored query not found for the given name and version",
				Status:  "not_found",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get stored query by name and version", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get stored query by name and version",
			Status:  "error",
		})
	}

	// Execute AQL query
	if err := h.OpenEHRService.QueryWithStream(ctx, c.Response().BodyWriter(), storedQuery.Query, aqlRequest.QueryParameters); err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to execute Stored AQL Version Post", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to execute Stored AQL Version Post",
			Status:  "error",
		})
	}

	auditCtx.Success()

	// Register event
	h.WebhookSaver.Enqueue(webhook.EventTypeQueryExecuted, map[string]any{
		"query_name": name,
		"version":    version,
		"parameters": aqlRequest.QueryParameters,
	})

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

	auditCtx := audit.From(c)

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	queries, err := h.OpenEHRService.ListStoredQueries(ctx, name)
	if err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to list stored queries", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to list stored queries",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(queries)
}

func (h *Handler) StoreQuery(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	if c.Accepts("text/plain") == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotAcceptable,
			Message: "Accept header must include text/plain",
			Status:  "not_acceptable",
		})
	}

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	queryType := c.Query("query_type")
	if queryType != "" && queryType != "AQL" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Unsupported query_type. Only 'AQL' is supported.",
			Status:  "bad_request",
		})
	}

	query := string(c.Body())
	if query == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Query in request body is required",
			Status:  "bad_request",
		})
	}

	// Check if query with the same name already exists
	_, err := h.OpenEHRService.GetQueryByName(ctx, name, "")
	if err == nil {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusConflict,
			Message: "Query with the given name already exists, system cannot update without knowing the target version, please use Store Query Version endpoint instead",
			Status:  "conflict",
		})
	}
	if err != ErrQueryNotFound {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to check existing query by name", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to store query",
			Status:  "error",
		})
	}

	// Store the query
	err = h.OpenEHRService.StoreQuery(ctx, name, "1.0.0", query)
	if err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to store query", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to store query",
			Status:  "error",
		})
	}

	auditCtx.Success()

	// Register event
	h.WebhookSaver.Enqueue(webhook.EventTypeQueryStored, map[string]any{
		"query_name": name,
		"version":    "1.0.0",
		"query":      query,
	})

	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) StoreQueryVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	if c.Accepts("text/plain") == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotAcceptable,
			Message: "Accept header must include text/plain",
			Status:  "not_acceptable",
		})
	}

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	version := c.Params("version")
	if version == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version path parameter is required",
			Status:  "bad_request",
		})
	}

	queryType := c.Query("query_type")
	if queryType != "" && queryType != "AQL" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Unsupported query_type. Only 'AQL' is supported.",
			Status:  "bad_request",
		})
	}

	query := string(c.Body())
	if query == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Query in request body is required",
			Status:  "bad_request",
		})
	}

	// Store the new version of the query
	err := h.OpenEHRService.StoreQuery(ctx, name, version, query)
	if err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to store query version", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to store query version",
			Status:  "error",
		})
	}

	auditCtx.Success()

	// Register event
	h.WebhookSaver.Enqueue(webhook.EventTypeQueryStored, map[string]any{
		"query_name": name,
		"version":    version,
		"query":      query,
	})

	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) GetStoredQueryAtVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	name := c.Params("qualified_query_name")
	if name == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "qualified_query_name path parameter is required",
			Status:  "bad_request",
		})
	}

	version := c.Params("version")
	if version == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "version path parameter is required",
			Status:  "bad_request",
		})
	}

	storedQuery, err := h.OpenEHRService.GetQueryByName(ctx, name, version)
	if err != nil {
		if err == ErrQueryNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Stored query not found for the given name and version",
				Status:  "error",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to get stored query by name and version", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get stored query by name and version",
			Status:  "error",
		})
	}

	auditCtx.Success()
	return c.Status(fiber.StatusOK).JSON(storedQuery)
}

func (h *Handler) DeleteEHRByID(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	ehrID, err := UUIDFromPath(c, auditCtx, "ehr_id")
	if err != nil {
		return err
	}
	auditCtx.Event.Details["ehr_id"] = ehrID

	err = h.OpenEHRService.DeleteEHR(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "EHR with the given ID not found",
				Status:  "error",
			})
		}

		h.Telemetry.Logger.ErrorContext(ctx, "Failed to delete EHR by ID", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete EHR by ID",
			Status:  "error",
		})
	}

	auditCtx.Success()

	// Register event
	h.WebhookSaver.Enqueue(webhook.EventTypeEHRDeleted, map[string]any{
		"ehr_id": ehrID,
	})

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) DeleteMultipleEHRs(c *fiber.Ctx) error {
	ctx := c.Context()

	auditCtx := audit.From(c)

	// Parse multiple ehr_id query parameters
	// Example: ?ehr_id=7d44b88c-4199-4bad-97dc-d78268e01398&ehr_id=297c3e91-7c17-4497-85dd-01e05aaae44e
	var ehrIDList []string
	for key, value := range c.Context().QueryArgs().All() {
		if string(key) == "ehr_id" {
			ehrIDList = append(ehrIDList, string(value))
		}
	}

	if len(ehrIDList) == 0 {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "At least one ehr_id query parameter is required",
			Status:  "bad_request",
		})
	}

	err := h.OpenEHRService.DeleteEHRBulk(ctx, ehrIDList)
	if err != nil {
		h.Telemetry.Logger.ErrorContext(ctx, "Failed to delete multiple EHRs", "error", err)

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete multiple EHRs",
			Status:  "error",
		})
	}

	auditCtx.Success()

	// Register event
	h.WebhookSaver.Enqueue(webhook.EventTypeEHRDeleted, map[string]any{
		"ehr_id_list": ehrIDList,
	})

	c.Status(fiber.StatusNoContent)
	return nil
}

func Accepts(c *fiber.Ctx, auditCtx *audit.Context, contentType string) error {
	if c.Accepts(contentType) == "" {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusNotAcceptable,
			Message: "Accept header must include " + contentType,
			Status:  "not_acceptable",
		})
	}

	return nil
}

func StringFromHeader(c *fiber.Ctx, auditCtx *audit.Context, name string) (string, error) {
	value := c.Get(name)
	if value == "" {
		return "", SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: name + " header is required",
			Status:  "bad_request",
		})
	}
	return value, nil
}

func ReturnTypeFromHeader(c *fiber.Ctx, auditCtx *audit.Context) (ReturnType, error) {
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return "", SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid Prefer header",
		})
	}
	return returnType, nil
}

func UUIDFromPath(c *fiber.Ctx, auditCtx *audit.Context, paramName string) (uuid.UUID, error) {
	valueStr := c.Params(paramName)
	if valueStr == "" {
		return uuid.Nil, SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: paramName + " path parameter is required",
			Status:  "bad_request",
		})
	}
	value, err := uuid.Parse(valueStr)
	if err != nil {
		return uuid.Nil, SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid " + paramName + " format",
			Status:  "bad_request",
		})
	}
	return value, nil
}

func StringFromPath(c *fiber.Ctx, auditCtx *audit.Context, value string) (string, error) {
	pathParam := c.Params(value)
	if pathParam == "" {
		return "", SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: value + " parameter is required",
			Status:  "bad_request",
		})
	}
	return pathParam, nil
}

func TimeFromQuery(c *fiber.Ctx, auditCtx *audit.Context, paramName string) (time.Time, error) {
	valueStr := c.Query(paramName)
	if valueStr == "" {
		return time.Time{}, SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: paramName + " query parameter is required",
			Status:  "bad_request",
		})
	}
	value, err := time.Parse(time.RFC3339, valueStr)
	if err != nil {
		return time.Time{}, SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid " + paramName + " format, must be RFC3339",
			Status:  "bad_request",
		})
	}
	return value, nil
}

func ParseBody(c *fiber.Ctx, auditCtx *audit.Context, out any) error {
	if len(c.Body()) == 0 {
		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Request body is required",
			Status:  "bad_request",
		})
	}

	if err := c.BodyParser(&out); err != nil {
		verr, ok := err.(util.ValidateError)
		if ok {
			return SendErrorResponse(c, auditCtx, ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Validation error in request body",
				Status:  "bad_request",
				Details: verr,
			})
		}

		return SendErrorResponse(c, auditCtx, ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Status:  "bad_request",
		})
	}
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

func SendErrorResponse(c *fiber.Ctx, auditCtx *audit.Context, errorRes ErrorResponse) error {
	auditCtx.Fail(errorRes.Status, errorRes.Message)
	return c.Status(errorRes.Code).JSON(map[string]ErrorResponse{
		"error": errorRes,
	})
}
