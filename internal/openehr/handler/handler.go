package handler

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/openehr"
	"github.com/freekieb7/gopenehr/internal/openehr/service"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	Config             *config.Config
	Logger             *slog.Logger
	EHRService         *service.EHRService
	DemographicService *service.DemographicService
	QueryService       *service.QueryService
}

func NewHandler(cfg *config.Config, logger *slog.Logger, ehrService *service.EHRService, demographicService *service.DemographicService, queryService *service.QueryService) Handler {
	return Handler{
		Config:             cfg,
		Logger:             logger,
		EHRService:         ehrService,
		DemographicService: demographicService,
		QueryService:       queryService,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/openehr/v1")

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
		"version":               h.Config.Version,
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

	c.Accepts("application/json")

	subjectID := c.Query("subject_id")
	if subjectID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("subject_id query parameters are required")
	}

	subjectNamespace := c.Query("subject_namespace")
	if subjectNamespace == "" {
		return c.Status(fiber.StatusBadRequest).SendString("subject_namespace query parameters are required")
	}

	ehr, err := h.EHRService.GetEHRBySubject(ctx, subjectID, subjectNamespace)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given subject")
		}
		h.Logger.ErrorContext(ctx, "Failed to get EHR by subject", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(ehr)
}

func (h *Handler) CreateEHR(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	// response type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Check for optional EHR_STATUS in the request body
	ehrStatus := h.EHRService.NewEHRStatus()
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&ehrStatus); err != nil {
			if err, ok := err.(util.ValidateError); ok {
				return c.Status(fiber.StatusBadRequest).JSON(err)
			}
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}
	}

	// Create EHR
	ehr, err := h.EHRService.CreateEHR(ctx, uuid.New(), ehrStatus)
	if err != nil {
		if err == service.ErrEHRStatusAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("EHR Status with the given UID already exists")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to create EHR", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Set response headers
	c.Set("ETag", "\""+ehr.EHRID.Value+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/ehr/"+ehr.EHRID.Value)

	c.Status(fiber.StatusCreated)
	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(ehr)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + ehr.EHRID.Value + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(ehr)
	}
}

func (h *Handler) GetEHR(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	ehr, err := h.EHRService.GetEHR(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get EHR by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(ehr)
}

func (h *Handler) CreateEHRWithID(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id format")
	}

	// response type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Check for optional EHR_STATUS in the request body
	ehrStatus := h.EHRService.NewEHRStatus()
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&ehrStatus); err != nil {
			if err, ok := err.(util.ValidateError); ok {
				return c.Status(fiber.StatusBadRequest).JSON(err)
			}

			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}
	}

	// Create EHR with specified ID and EHR_STATUS
	ehr, err := h.EHRService.CreateEHR(ctx, ehrID, ehrStatus)
	if err != nil {
		if err == service.ErrEHRAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("EHR with the given ID already exists")
		}
		if err == service.ErrEHRStatusAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("EHR Status with the given UID already exists")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to create EHR", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	c.Set("ETag", "\""+ehrID.String()+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/ehr/"+ehrID.String())

	c.Status(fiber.StatusCreated)
	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(ehr)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + ehrID.String() + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return nil
	}
}

func (h *Handler) GetEHRStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	var filterAtTime time.Time
	if atTimeStr := c.Query("version_at_time"); atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid version_at_time format. Use RFC3339 format.")
		}
		filterAtTime = parsedTime
	}

	ehrStatus, err := h.EHRService.GetEHRStatus(ctx, ehrID, filterAtTime, "")
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR Status not found for the given EHR ID at the specified time")
		}
		h.Logger.ErrorContext(ctx, "Failed to get EHR Status at time", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(ehrStatus)
}

func (h *Handler) UpdateEhrStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// response type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	var requestEhrStatus openehr.EHR_STATUS
	if err := c.BodyParser(&requestEhrStatus); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Check collision using If-Match header
	currentEHRStatus, err := h.EHRService.GetEHRStatus(ctx, ehrID, time.Time{}, "")
	if err != nil {
		if err == service.ErrEHRStatusNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR Status not found for the given EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get current EHR Status", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	currentEHRStatusID, ok := currentEHRStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Current EHR Status UID is not of type OBJECT_VERSION_ID")
	}

	// Check collision using If-Match header
	if currentEHRStatusID.Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("EHR Status has been modified since the provided version")
	}

	// Proceed to update EHR Status
	updatedEHRStatus, err := h.EHRService.UpdateEHRStatus(ctx, ehrID, requestEhrStatus)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR Status not found for the given EHR ID")
		}
		if err == service.ErrEHRStatusAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("EHR Status with the given UID already exists")
		}
		if err == service.ErrEHRStatusVersionLowerOrEqualToCurrent {
			return c.Status(fiber.StatusBadRequest).SendString("EHR Status version in request body must be incremented")
		}
		if err == service.ErrInvalidEHRStatusUIDMismatch {
			return c.Status(fiber.StatusBadRequest).SendString("EHR Status UID HIER_OBJECT_ID in request body does not match current EHR Status UID")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to update EHR Status", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	updatedEHRStatusID := updatedEHRStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedEHRStatusID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/ehr/"+ehrID+"/ehr_status/"+updatedEHRStatusID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedEHRStatus)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + updatedEHRStatusID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedEHRStatus)
	}
}

func (h *Handler) GetEHRStatusByVersionID(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version_uid parameter is required")
	}

	ehrStatus, err := h.EHRService.GetEHRStatus(ctx, ehrID, time.Time{}, versionUID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR Status not found for the given EHR ID and version UID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get EHR Status at time", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(ehrStatus)
}

func (h *Handler) GetVersionedEHRStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionedStatus, err := h.EHRService.GetVersionedEHRStatus(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned EHR Status not found for the given EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(versionedStatus)
}

func (h *Handler) GetVersionedEHRStatusRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	revisionHistory, err := h.EHRService.GetVersionedEHRStatusRevisionHistory(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned EHR Status revision history not found for the given EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status revision history", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(revisionHistory)
}

func (h *Handler) GetVersionedEHRStatusVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	}

	var filterAtTime time.Time
	atTimeStr := c.Query("version_at_time")
	if atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid version_at_time format. Use RFC3339 format.")
		}
		filterAtTime = parsedTime
	}

	ehrStatusVersionJSON, err := h.EHRService.GetVersionedEHRStatusVersionAsJSON(ctx, ehrID, filterAtTime, "")
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned EHR Status version not found for the given EHR ID at the specified time")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status version at time", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(ehrStatusVersionJSON)
}

func (h *Handler) GetVersionedEHRStatusVersionByID(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	}

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version_uid parameter is required")
	}

	ehrStatusVersionJSON, err := h.EHRService.GetVersionedEHRStatusVersionAsJSON(ctx, ehrID, time.Time{}, versionUID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned EHR Status version not found for the given EHR ID and version UID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status version by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(ehrStatusVersionJSON)
}

func (h *Handler) CreateComposition(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	}

	// response type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Parse new composition from request body
	var requestComposition openehr.COMPOSITION
	if err := c.BodyParser(&requestComposition); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Check if EHR exists
	_, err = h.EHRService.GetEHR(ctx, ehrID.String())
	if err != nil && err != service.ErrEHRNotFound {
		h.Logger.ErrorContext(ctx, "Failed to check if EHR exists", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	if err == service.ErrEHRNotFound {
		return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
	}

	// Create Composition
	newComposition, err := h.EHRService.CreateComposition(ctx, ehrID, requestComposition)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Composition with the given UID already exists")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to create Composition", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	compositionID := newComposition.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+compositionID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/ehr/"+ehrID.String()+"/composition/"+compositionID)
	c.Status(fiber.StatusCreated)

	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(newComposition)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + compositionID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(newComposition)
	}
}

func (h *Handler) GetComposition(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	}

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	composition, err := h.EHRService.GetComposition(ctx, ehrID, uidBasedID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Composition not found for the given UID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Composition by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(composition)
}

func (h *Handler) UpdateComposition(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	}

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// response type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Check collision using If-Match header
	currentComposition, err := h.EHRService.GetComposition(ctx, ehrID, strings.Split(uidBasedID, "::")[0])
	if err != nil {
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Composition not found for the given EHR ID and UID")
		}

		h.Logger.ErrorContext(ctx, "Failed to get current Composition", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Composition
	if currentComposition.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Composition has been modified since the provided version")
	}

	// Parse updated composition from request body
	var requestComposition openehr.COMPOSITION
	if err := c.BodyParser(&requestComposition); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Composition
	composition, err := h.EHRService.UpdateComposition(ctx, ehrID, requestComposition)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Composition not found for the given UID")
		}
		if err == service.ErrCompositionAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Composition with the given UID already exists")
		}
		if err == service.ErrCompositionVersionLowerOrEqualToCurrent {
			return c.Status(fiber.StatusBadRequest).SendString("Composition version in request body is lower or equal to the current version")
		}
		if err == service.ErrInvalidCompositionUIDMismatch {
			return c.Status(fiber.StatusBadRequest).SendString("Composition UID HIER_OBJECT_ID in request body does not match current Composition UID")
		}
		if err == service.ErrCompositionUIDNotProvided {
			return c.Status(fiber.StatusBadRequest).SendString("Composition UID must be provided in the request body for update")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to update Composition", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	compositionID := composition.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+compositionID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/ehr/"+ehrID.String()+"/composition/"+compositionID)

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
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	}

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if strings.Count(uidBasedID, "::") != 2 {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete Composition by versioned object ID. Please provide the object version ID.")
	}

	// Check existence before deletion
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentComposition, err := h.EHRService.GetComposition(ctx, ehrID, versionedObjectID)
	if err != nil {
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Composition not found for the given UID and EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Composition by ID before deletion", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Composition
	if currentComposition.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != uidBasedID {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Composition has been modified since the provided version")
	}

	// Proceed to delete Composition
	if err := h.EHRService.DeleteComposition(ctx, ehrID, currentComposition.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value); err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Composition not found for the given UID")
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Composition by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) GetVersionedCompositionByID(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}

	versionedComposition, err := h.EHRService.GetVersionedComposition(ctx, ehrID, versionedObjectID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Composition not found for the given versioned object ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Composition by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(versionedComposition)
}

func (h *Handler) GetVersionedCompositionRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}

	revisionHistory, err := h.EHRService.GetVersionedCompositionRevisionHistory(ctx, ehrID, versionedObjectID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Composition revision history not found for the given versioned object ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Composition revision history", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(revisionHistory)
}

func (h *Handler) GetVersionedCompositionVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}

	var filterAtTime time.Time
	atTimeStr := c.Query("version_at_time")
	if atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid version_at_time format. Use RFC3339 format.")
		}
		filterAtTime = parsedTime
	}

	versionJSON, err := h.EHRService.GetVersionedCompositionVersionJSON(ctx, ehrID, versionedObjectID, filterAtTime, "")
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Composition version not found for the given versioned object ID at the specified time")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Composition version at time", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(versionJSON)
}

func (h *Handler) GetVersionedCompositionVersionByID(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version_uid parameter is required")
	}

	if versionedObjectID != strings.Split(versionUID, "::")[0] {
		return c.Status(fiber.StatusBadRequest).SendString("version_uid does not match the versioned_object_id")
	}

	versionJSON, err := h.EHRService.GetVersionedCompositionVersionJSON(ctx, ehrID, versionedObjectID, time.Time{}, versionUID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Composition version not found for the given versioned object ID and version UID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Composition version by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(versionJSON)
}

func (h *Handler) CreateDirectory(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	}

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	if len(c.Body()) <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Request body is required")
	}

	var requestDirectory openehr.FOLDER
	if err := c.BodyParser(&requestDirectory); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	directory, err := h.EHRService.CreateDirectory(ctx, ehrID, requestDirectory)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrDirectoryAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Directory with the given UID already exists")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to create Directory", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	directoryID := directory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+directoryID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/ehr/"+ehrID.String()+"/directory/"+directoryID)

	c.Status(fiber.StatusCreated)
	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(directory)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + directoryID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(directory)
	}
}

func (h *Handler) UpdateDirectory(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Check collision using If-Match header
	currentDirectory, err := h.EHRService.GetDirectory(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrDirectoryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory not found for the given EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Directory", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Directory
	if currentDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Directory has been modified since the provided version")
	}

	// Parse updated directory from request body
	var requestDirectory openehr.FOLDER
	if err := c.BodyParser(&requestDirectory); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	directory, err := h.EHRService.UpdateDirectory(ctx, ehrID, requestDirectory)
	if err != nil {
		if err == service.ErrDirectoryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory not found for the given EHR ID")
		}
		if err == service.ErrDirectoryAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Directory with the given UID already exists")
		}
		if err == service.ErrDirectoryVersionLowerOrEqualToCurrent {
			return c.Status(fiber.StatusBadRequest).SendString("Directory version in request body is lower or equal to the current version")
		}
		if err == service.ErrInvalidDirectoryUIDMismatch {
			return c.Status(fiber.StatusBadRequest).SendString("Directory UID HIER_OBJECT_ID in request body does not match current Directory UID")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to update Directory", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	updatedDirectoryID := directory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedDirectoryID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/ehr/"+ehrID.String()+"/directory/"+updatedDirectoryID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(directory)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + updatedDirectoryID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(directory)
	}
}

func (h *Handler) DeleteDirectory(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrIDStr := c.Params("ehr_id")
	if ehrIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	currentDirectory, err := h.EHRService.GetDirectory(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrDirectoryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory not found for the given EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Directory", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Directory
	if currentDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Directory has been modified since the provided version")
	}

	if err := h.EHRService.DeleteDirectory(ctx, ehrID); err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrDirectoryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory not found for the given EHR ID")
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Directory", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) GetFolderInDirectoryVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	var filterAtTime time.Time
	atTimeStr := c.Query("version_at_time")
	if atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid version_at_time format. Use RFC3339 format.")
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

	folder, err := h.EHRService.GetFolderInDirectoryVersion(ctx, ehrID, filterAtTime, "", pathParts)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory not found for the given EHR ID at the specified time")
		}
		if err == service.ErrDirectoryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory version not found at the specified time for the given EHR ID")
		}
		if err == service.ErrFolderNotFoundInDirectory {
			return c.Status(fiber.StatusNotFound).SendString("Folder not found in Directory version at the specified time for the given path")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Folder in Directory version at time", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(folder)
}

func (h *Handler) GetFolderInDirectoryVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version_uid parameter is required")
	}

	path := c.Query("path")
	var pathParts []string
	if path != "" {
		for part := range strings.SplitSeq(path, "/") {
			pathParts = append(pathParts, strings.TrimSpace(part))
		}
	}

	folder, err := h.EHRService.GetFolderInDirectoryVersion(ctx, ehrID, time.Time{}, versionUID, pathParts)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory not found for the given EHR ID")
		}
		if err == service.ErrDirectoryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory version not found at the specified time for the given EHR ID")
		}
		if err == service.ErrFolderNotFoundInDirectory {
			return c.Status(fiber.StatusNotFound).SendString("Folder not found in Directory version for the given path")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Folder in Directory version", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

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

	// contribution, err := h.EHRService.CreateContribution(ctx, ehrID, newContribution)
	// if err != nil {
	// 	if err == service.ErrContributionAlreadyExists {
	// 		return c.Status(fiber.StatusConflict).SendString("Contribution with the given UID already exists")
	// 	}

	// 	h.Logger.ErrorContext(ctx, "Failed to create Contribution", "error", err)
	// 	c.Status(http.StatusInternalServerError)
	// 	return nil
	// }

	// // Determine response
	// contributionID := contribution.UID.Value
	// c.Set("ETag", "\""+contributionID+"\"")
	// c.Set("Location", h.Config.Host+"/openehr/v1/ehr/"+ehrID.String()+"/contribution/"+contributionID)

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
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}
	ehrID, err := uuid.Parse(ehrIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ehr_id parameter format")
	}

	contributionUID := c.Params("contribution_uid")
	if contributionUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("contribution_uid parameter is required")
	}

	contribution, err := h.EHRService.GetContribution(ctx, ehrID, contributionUID)
	if err != nil {
		if err == service.ErrContributionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Contribution not found for the given EHR ID and Contribution UID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Contribution by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
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
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Parse new agent from request body
	var newAgent openehr.AGENT
	if err := c.BodyParser(&newAgent); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Agent", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	agent, err := h.DemographicService.CreateAgent(ctx, newAgent)
	if err != nil {
		if err == service.ErrAgentAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Agent with the given UID already exists")
		}
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to create Agent", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	agentID := agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+agentID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/agent/"+agentID)

	c.Status(fiber.StatusCreated)
	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(agent)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + agentID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(agent)
	}
}

func (h *Handler) GetAgent(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if strings.Count(uidBasedID, "::") == 3 {
		// Is version id
		agent, err := h.DemographicService.GetAgentAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == service.ErrAgentNotFound {
				return c.Status(fiber.StatusNotFound).SendString("Agent not found for the given agent ID")
			}
			h.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)
			c.Status(http.StatusInternalServerError)
			return nil
		}

		return c.Status(fiber.StatusOK).JSON(agent)
	}

	// Is versioned party id
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	agent, err := h.DemographicService.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrAgentNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Agent not found for the given agent ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(agent)
}

func (h *Handler) UpdateAgent(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionID := c.Params("uid_based_id")
	if versionID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Ensure Agent exists before update
	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	currentAgent, err := h.DemographicService.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrAgentNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Agent not found for the given agent ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Agent by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Agent
	if currentAgent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Agent has been modified since the provided version")
	}

	// Parse updated agent from request body
	var requestAgent openehr.AGENT
	if err := c.BodyParser(&requestAgent); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Agent", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Agent
	updatedAgent, err := h.DemographicService.UpdateAgent(ctx, versionedPartyID, requestAgent)
	if err != nil {
		if err == service.ErrAgentNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Agent not found for the given agent ID")
		}
		if err == service.ErrAgentVersionLowerOrEqualToCurrent {
			return c.Status(fiber.StatusConflict).SendString("Agent version is lower than or equal to the current version")
		}
		if err == service.ErrInvalidAgentUIDMismatch {
			return c.Status(fiber.StatusBadRequest).SendString("Agent UID HIER_OBJECT_ID in request body does not match current Agent UID")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to update Agent", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	updatedAgentID := updatedAgent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedAgentID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/agent/"+updatedAgentID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedAgent)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + updatedAgentID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedAgent)
	}
}

func (h *Handler) DeleteAgent(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if strings.Count(uidBasedID, "::") != 2 {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete Agent by versioned object ID. Please provide the object version ID.")
	}

	// Check existence before deletion
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	currentAgent, err := h.DemographicService.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrAgentNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Agent not found for the given agent ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Agent by ID before deletion", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Agent
	if currentAgent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != uidBasedID {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Agent has been modified since the provided version")
	}

	// Proceed to delete Agent
	if err := h.DemographicService.DeleteAgent(ctx, uidBasedID); err != nil {
		if err == service.ErrAgentNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Agent not found for the given agent ID")
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Agent by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
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
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Parse new group from request body
	var newGroup openehr.GROUP
	if err := c.BodyParser(&newGroup); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Group", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	group, err := h.DemographicService.CreateGroup(ctx, newGroup)
	if err != nil {
		if err == service.ErrGroupAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Group with the given UID already exists")
		}
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to create Group", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	groupID := group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+groupID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/group/"+groupID)

	c.Status(fiber.StatusCreated)
	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(group)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + groupID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(group)
	}
}

func (h *Handler) GetGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if strings.Count(uidBasedID, "::") == 3 {
		// Is version id
		group, err := h.DemographicService.GetGroupAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == service.ErrGroupNotFound {
				return c.Status(fiber.StatusNotFound).SendString("Group not found for the given group ID")
			}
			h.Logger.ErrorContext(ctx, "Failed to get Group by ID", "error", err)
			c.Status(http.StatusInternalServerError)
			return nil
		}

		return c.Status(fiber.StatusOK).JSON(group)

	}
	// Is versioned party id
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	group, err := h.DemographicService.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrGroupNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Group not found for the given group ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Group by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(group)
}

func (h *Handler) UpdateGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionID := c.Params("uid_based_id")
	if versionID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// return type
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Ensure Agent exists before update
	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	currentGroup, err := h.DemographicService.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrGroupNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Group not found for the given group ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Group by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Group
	if currentGroup.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Group has been modified since the provided version")
	}

	// Parse updated group from request body
	var requestGroup openehr.GROUP
	if err := c.BodyParser(&requestGroup); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Group", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Group
	updatedGroup, err := h.DemographicService.UpdateGroup(ctx, versionedPartyID, requestGroup)
	if err != nil {
		if err == service.ErrGroupNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Group not found for the given group ID")
		}
		if err == service.ErrGroupVersionLowerOrEqualToCurrent {
			return c.Status(fiber.StatusConflict).SendString("Group version is lower than or equal to the current version")
		}
		if err == service.ErrInvalidGroupUIDMismatch {
			return c.Status(fiber.StatusBadRequest).SendString("Group UID HIER_OBJECT_ID in request body does not match current Group UID")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to update Agent", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	updatedGroupID := updatedGroup.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedGroupID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/group/"+updatedGroupID)

	switch returnType {
	case ReturnTypeMinimal:
		c.Status(fiber.StatusNoContent)
		return nil
	case ReturnTypeRepresentation:
		return c.Status(fiber.StatusOK).JSON(updatedGroup)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + updatedGroupID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.Status(fiber.StatusOK).JSON(updatedGroup)
	}
}

func (h *Handler) DeleteGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if !strings.Contains(uidBasedID, "::") {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete Group by versioned object ID. Please provide the object version ID.")
	}

	// Check existence before deletion
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	currentGroup, err := h.DemographicService.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrGroupNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Group not found for the given group ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Group by ID before deletion", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Group
	if currentGroup.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != uidBasedID {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Group has been modified since the provided version")
	}

	// Proceed to delete Group
	if err := h.DemographicService.DeleteGroup(ctx, uidBasedID); err != nil {
		if err == service.ErrGroupNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Group not found for the given group ID")
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Group by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
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
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Parse new person from request body
	var newPerson openehr.PERSON
	if err := c.BodyParser(&newPerson); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Person", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	person, err := h.DemographicService.CreatePerson(ctx, newPerson)
	if err != nil {
		if err == service.ErrPersonAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Person with the given UID already exists")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to create Person", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	personID := person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+personID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/person/"+personID)

	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(person)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + personID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(person)
	}
}

func (h *Handler) GetPerson(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if strings.Count(uidBasedID, "::") == 3 {
		// Is version id
		person, err := h.DemographicService.GetPersonAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == service.ErrPersonNotFound {
				return c.Status(fiber.StatusNotFound).SendString("Person not found for the given person ID")
			}
			h.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)
			c.Status(http.StatusInternalServerError)
			return nil
		}
		return c.Status(fiber.StatusOK).JSON(person)
	}

	// Is versioned party id
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	person, err := h.DemographicService.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrPersonNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Person not found for the given person ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(person)
}

func (h *Handler) UpdatePerson(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionID := c.Params("uid_based_id")
	if versionID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// Parse return type from Prefer header
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Ensure Person exists before update
	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	currentPerson, err := h.DemographicService.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrPersonNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Person not found for the given person ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Person by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Person
	if currentPerson.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Person has been modified since the provided version")
	}

	// Parse updated person from request body
	var requestPerson openehr.PERSON
	if err := c.BodyParser(&requestPerson); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Person", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Person
	updatePerson, err := h.DemographicService.UpdatePerson(ctx, versionedPartyID, requestPerson)
	if err != nil {
		if err == service.ErrPersonNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Person not found for the given person ID")
		}
		if err == service.ErrPersonVersionLowerOrEqualToCurrent {
			return c.Status(fiber.StatusConflict).SendString("Person version is lower than or equal to the current version")
		}
		if err == service.ErrInvalidPersonUIDMismatch {
			return c.Status(fiber.StatusBadRequest).SendString("Person UID HIER_OBJECT_ID in request body does not match current Person UID")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to update Person", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	updatedPersonID := updatePerson.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedPersonID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/person/"+updatedPersonID)

	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(updatePerson)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + updatedPersonID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(updatePerson)
	}
}

func (h *Handler) DeletePerson(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if strings.Count(uidBasedID, "::") != 3 {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete Person by versioned object ID. Please provide the object version ID.")
	}

	// Check existence before deletion
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	currentPerson, err := h.DemographicService.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrPersonNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Person not found for the given person ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Person by ID before deletion", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Person
	if currentPerson.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != uidBasedID {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Person has been modified since the provided version")
	}

	// Proceed to delete Person
	if err := h.DemographicService.DeletePerson(ctx, uidBasedID); err != nil {
		if err == service.ErrPersonNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Person not found for the given person ID")
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Person by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
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
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Parse new organisation from request body
	var newOrganisation openehr.ORGANISATION
	if err := c.BodyParser(&newOrganisation); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Organisation", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	organisation, err := h.DemographicService.CreateOrganisation(ctx, newOrganisation)
	if err != nil {
		if err == service.ErrOrganisationAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Organisation with the given UID already exists")
		}
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to create Organisation", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	organisationID := organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+organisationID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/organisation/"+organisationID)

	c.Status(fiber.StatusCreated)
	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(organisation)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + organisationID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(organisation)
	}
}

func (h *Handler) GetOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if strings.Count(uidBasedID, "::") == 3 {
		// Is version id
		organisation, err := h.DemographicService.GetOrganisationAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == service.ErrOrganisationNotFound {
				return c.Status(fiber.StatusNotFound).SendString("Organisation not found for the given organisation ID")
			}
			h.Logger.ErrorContext(ctx, "Failed to get Organisation by ID", "error", err)
			c.Status(http.StatusInternalServerError)
			return nil
		}

		return c.Status(fiber.StatusOK).JSON(organisation)
	}

	// Is versioned party id
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	organisation, err := h.DemographicService.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrOrganisationNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Organisation not found for the given organisation ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Organisation by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(organisation)
}

func (h *Handler) UpdateOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionID := c.Params("uid_based_id")
	if versionID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// Parse return type from Prefer header
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Ensure Organisation exists before update
	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	currentOrganisation, err := h.DemographicService.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrOrganisationNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Organisation not found for the given organisation ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Organisation by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Organisation
	if currentOrganisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Organisation has been modified since the provided version")
	}

	// Parse updated organisation from request body
	var requestOrganisation openehr.ORGANISATION
	if err := c.BodyParser(&requestOrganisation); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Organisation", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Organisation
	organisation, err := h.DemographicService.UpdateOrganisation(ctx, versionedPartyID, requestOrganisation)
	if err != nil {
		if err == service.ErrOrganisationNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Organisation not found for the given organisation ID")
		}
		if err == service.ErrOrganisationVersionLowerOrEqualToCurrent {
			return c.Status(fiber.StatusConflict).SendString("Organisation version is lower than or equal to the current version")
		}
		if err == service.ErrInvalidOrganisationUIDMismatch {
			return c.Status(fiber.StatusConflict).SendString("Organisation with the given UID already exists")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to update Organisation", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	organisationID := organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+organisationID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/organisation/"+organisationID)

	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(organisation)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + organisationID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(organisation)
	}
}

func (h *Handler) DeleteOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if strings.Count(uidBasedID, "::") != 2 {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete Organisation by versioned object ID. Please provide the object version ID.")
	}

	// Check existence before deletion
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	currentOrganisation, err := h.DemographicService.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrOrganisationNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Organisation not found for the given organisation ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Organisation by ID before deletion", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Organisation
	if currentOrganisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != uidBasedID {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Organisation has been modified since the provided version")
	}

	// Proceed to delete Organisation
	if err := h.DemographicService.DeleteOrganisation(ctx, uidBasedID); err != nil {
		if err == service.ErrOrganisationNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Organisation not found for the given organisation ID")
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Organisation by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
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
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Parse new role from request body
	var newRole openehr.ROLE
	if err := c.BodyParser(&newRole); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Role", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	role, err := h.DemographicService.CreateRole(ctx, newRole)
	if err != nil {
		if err == service.ErrRoleAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Role with the given UID already exists")
		}
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to create Role", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	roleID := role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+roleID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/role/"+roleID)

	c.Status(fiber.StatusCreated)
	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(role)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + roleID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(role)
	}
}

func (h *Handler) GetRole(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if strings.Count(uidBasedID, "::") == 3 {
		// Is version id
		role, err := h.DemographicService.GetRoleAtVersion(ctx, uidBasedID)
		if err != nil {
			if err == service.ErrRoleNotFound {
				return c.Status(fiber.StatusNotFound).SendString("Role not found for the given role ID")
			}
			h.Logger.ErrorContext(ctx, "Failed to get Role by ID", "error", err)
			c.Status(http.StatusInternalServerError)
			return nil
		}

		return c.Status(fiber.StatusOK).JSON(role)
	}

	// Is versioned party id
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	role, err := h.DemographicService.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Role not found for the given role ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Role by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(role)
}

func (h *Handler) UpdateRole(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionID := c.Params("uid_based_id")
	if versionID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// Parse return type from Prefer header
	returnType := ReturnType(c.Get("Prefer", string(ReturnTypeMinimal)))
	if !returnType.IsValid() {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Prefer header value")
	}

	// Ensure Role exists before update
	versionedPartyID, err := uuid.Parse(strings.Split(versionID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	currentRole, err := h.DemographicService.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Role not found for the given role ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get current Role by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Role
	if currentRole.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Role has been modified since the provided version")
	}

	// Parse updated role from request body
	var requestRole openehr.ROLE
	if err := c.BodyParser(&requestRole); err != nil {
		if err, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Role", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Role
	updatedRole, err := h.DemographicService.UpdateRole(ctx, versionedPartyID, requestRole)
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Role not found for the given role ID")
		}
		if err == service.ErrRoleVersionLowerOrEqualToCurrent {
			return c.Status(fiber.StatusConflict).SendString("Role version is lower than or equal to the current version")
		}
		if err == service.ErrInvalidRoleUIDMismatch {
			return c.Status(fiber.StatusBadRequest).SendString("Role UID HIER_OBJECT_ID in request body does not match current Role UID")
		}
		if validationErrs, ok := err.(util.ValidateError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErrs)
		}

		h.Logger.ErrorContext(ctx, "Failed to update Role", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Determine response
	updatedRoleID := updatedRole.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value
	c.Set("ETag", "\""+updatedRoleID+"\"")
	c.Set("Location", h.Config.Host+"/openehr/v1/role/"+updatedRoleID)

	switch returnType {
	case ReturnTypeMinimal:
		return nil
	case ReturnTypeRepresentation:
		return c.JSON(updatedRole)
	case ReturnTypeIdentifier:
		return c.JSON(`{"uid":"` + updatedRoleID + `"}`)
	default:
		h.Logger.ErrorContext(ctx, "Unhandled Prefer header value", "value", returnType)
		return c.JSON(updatedRole)
	}
}

func (h *Handler) DeleteRole(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if !strings.Contains(uidBasedID, "::") {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete Role by versioned object ID. Please provide the object version ID.")
	}

	// Check existence before deletion
	versionedPartyID, err := uuid.Parse(strings.Split(uidBasedID, "::")[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid uid_based_id format")
	}

	currentRole, err := h.DemographicService.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Role not found for the given role ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Role by ID before deletion", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Check if provided version matches current version
	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Role
	if currentRole.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != uidBasedID {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Role has been modified since the provided version")
	}

	// Proceed to delete Role
	if err := h.DemographicService.DeleteRole(ctx, uidBasedID); err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Role not found for the given role ID")
		}

		h.Logger.ErrorContext(ctx, "Failed to delete Role by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *Handler) GetVersionedParty(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionedObjectIDStr := c.Params("versioned_object_id")
	if versionedObjectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}
	versionedObjectID, err := uuid.Parse(versionedObjectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid versioned_object_id format")
	}

	party, err := h.DemographicService.GetVersionedParty(ctx, versionedObjectID)
	if err != nil {
		if err == service.ErrVersionedPartyNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Party not found for the given ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(party)
}

func (h *Handler) GetVersionedPartyRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionedObjectIDStr := c.Params("versioned_object_id")
	if versionedObjectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}
	versionedObjectID, err := uuid.Parse(versionedObjectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid versioned_object_id format")
	}

	revisionHistory, err := h.DemographicService.GetVersionedPartyRevisionHistory(ctx, versionedObjectID)
	if err != nil {
		if err == service.ErrRevisionHistoryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Party not found for the given ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party Revision History by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(revisionHistory)
}

func (h *Handler) GetVersionedPartyVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionedObjectIDStr := c.Params("versioned_object_id")
	if versionedObjectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}
	versionedObjectID, err := uuid.Parse(versionedObjectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid versioned_object_id format")
	}

	// Parse version_at_time query parameter
	var filterAtTime time.Time
	atTimeStr := c.Query("version_at_time")
	if atTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, atTimeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid version_at_time format. Use RFC3339 format.")
		}
		filterAtTime = parsedTime
	}

	partyVersionJSON, err := h.DemographicService.GetVersionedPartyVersionJSON(ctx, versionedObjectID, filterAtTime, "")
	if err != nil {
		if err == service.ErrVersionedPartyNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Party not found for the given ID")
		}
		if err == service.ErrVersionedPartyVersionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Party version not found for the given time")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party Version at Time by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(partyVersionJSON)
}

func (h *Handler) GetVersionedPartyVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	versionedObjectIDStr := c.Params("versioned_object_id")
	if versionedObjectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}
	versionedObjectID, err := uuid.Parse(versionedObjectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid versioned_object_id format")
	}

	versionID := c.Params("version_id")
	if versionID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version_id parameter is required")
	}

	partyVersionJSON, err := h.DemographicService.GetVersionedPartyVersionJSON(ctx, versionedObjectID, time.Time{}, versionID)
	if err != nil {
		if err == service.ErrVersionedPartyNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Party not found for the given ID")
		}
		if err == service.ErrVersionedPartyVersionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Party version not found for the given version ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party Version by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

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

	// contribution, err := h.DemographicService.CreateContribution(ctx, newContribution)
	// if err != nil {
	// 	h.Logger.ErrorContext(ctx, "Failed to create Demographic Contribution", "error", err)
	// 	c.Status(http.StatusInternalServerError)
	// 	return nil
	// }

	// // Determine response
	// contributionID := contribution.UID.Value
	// c.Set("ETag", "\""+contributionID+"\"")
	// c.Set("Location", h.Config.Host+"/openehr/v1/demographic/contribution/"+contributionID)

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
	ctx := c.Context()

	c.Accepts("application/json")

	contributionUID := c.Params("contribution_uid")
	if contributionUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("contribution_uid parameter is required")
	}

	contribution, err := h.DemographicService.GetContribution(ctx, contributionUID)
	if err != nil {
		if err == service.ErrContributionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Contribution not found for the given EHR ID and Contribution UID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Contribution by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(contribution)
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
		return c.Status(fiber.StatusBadRequest).SendString("q query parameter is required")
	}

	ehrID := c.Query("ehr_id")
	if ehrID != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Ad Hoc AQL with ehr_id not implemented yet")
	}

	fetch := c.Query("fetch")
	if fetch != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Ad Hoc AQL with fetch not implemented yet")
	}

	offset := c.Query("offset")
	if offset != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Ad Hoc AQL with offset not implemented yet")
	}

	queryParametersStr := c.Query("query_parameters")
	queryParametersURLValues, err := url.ParseQuery(queryParametersStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid query_parameters format")
	}

	queryParameters := make(map[string]any)
	for key, values := range queryParametersURLValues {
		if len(values) > 0 {
			queryParameters[key] = values[0]
		}
	}

	// Execute AQL query
	if err := h.QueryService.QueryAndCopyTo(ctx, c.Response().BodyWriter(), query, queryParameters); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Ad Hoc AQL", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
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

	var aqlRequest AdHocAQLRequest
	if err := c.BodyParser(&aqlRequest); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse AQL request body", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if aqlRequest.Query == "" {
		return c.Status(fiber.StatusBadRequest).SendString("query field is required in the request body")
	}

	if aqlRequest.EHRID != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Ad Hoc AQL Post with ehr_id not implemented yet")
	}

	if aqlRequest.Fetch != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Ad Hoc AQL Post with fetch not implemented yet")
	}

	if aqlRequest.Offset != 0 {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Ad Hoc AQL Post with offset not implemented yet")
	}

	// Execute AQL query
	if err := h.QueryService.QueryAndCopyTo(ctx, c.Response().BodyWriter(), aqlRequest.Query, aqlRequest.QueryParameters); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Ad Hoc AQL Post", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) ExecuteStoredAQL(c *fiber.Ctx) error {
	ctx := c.Context()

	name := c.Params("qualified_query_name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("qualified_query_name path parameter is required")
	}

	ehrID := c.Query("ehr_id")
	if ehrID != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL with ehr_id not implemented yet")
	}

	fetch := c.Query("fetch")
	if fetch != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL with fetch not implemented yet")
	}

	offset := c.Query("offset")
	if offset != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL with offset not implemented yet")
	}

	queryParametersStr := c.Query("query_parameters")
	queryParametersURLValues, err := url.ParseQuery(queryParametersStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid query_parameters format")
	}

	queryParameters := make(map[string]any)
	for key, values := range queryParametersURLValues {
		if len(values) > 0 {
			queryParameters[key] = values[0]
		}
	}

	// Retrieve stored query by name
	storedQuery, err := h.QueryService.GetQueryByName(ctx, name, "")
	if err != nil {
		if err == service.ErrQueryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Stored query not found for the given name")
		}
		h.Logger.ErrorContext(ctx, "Failed to get stored query by name", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Execute AQL query
	if err := h.QueryService.QueryAndCopyTo(ctx, c.Response().BodyWriter(), storedQuery.Query, nil); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Stored AQL", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
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

	name := c.Params("qualified_query_name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("qualified_query_name path parameter is required")
	}

	var aqlRequest StoredAQLRequest
	if err := c.BodyParser(&aqlRequest); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse Stored AQL request body", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if aqlRequest.EHRID != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Post with ehr_id not implemented yet")
	}

	if aqlRequest.Fetch != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Post with fetch not implemented yet")
	}

	if aqlRequest.Offset != 0 {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Post with offset not implemented yet")
	}

	// Retrieve stored query by name
	storedQuery, err := h.QueryService.GetQueryByName(ctx, name, "")
	if err != nil {
		if err == service.ErrQueryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Stored query not found for the given name")
		}
		h.Logger.ErrorContext(ctx, "Failed to get stored query by name", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Execute AQL query
	if err := h.QueryService.QueryAndCopyTo(ctx, c.Response().BodyWriter(), storedQuery.Query, aqlRequest.QueryParameters); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Stored AQL Post", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) ExecuteStoredAQLVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	name := c.Params("qualified_query_name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("qualified_query_name path parameter is required")
	}

	version := c.Params("version")
	if version == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version path parameter is required")
	}

	ehrID := c.Query("ehr_id")
	if ehrID != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Version with ehr_id not implemented yet")
	}

	fetch := c.Query("fetch")
	if fetch != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Version with fetch not implemented yet")
	}

	offset := c.Query("offset")
	if offset != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Version with offset not implemented yet")
	}

	queryParametersStr := c.Query("query_parameters")
	queryParametersURLValues, err := url.ParseQuery(queryParametersStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid query_parameters format")
	}

	queryParameters := make(map[string]any)
	for key, values := range queryParametersURLValues {
		if len(values) > 0 {
			queryParameters[key] = values[0]
		}
	}

	// Retrieve stored query by name and version
	storedQuery, err := h.QueryService.GetQueryByName(ctx, name, version)
	if err != nil {
		if err == service.ErrQueryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Stored query not found for the given name and version")
		}
		h.Logger.ErrorContext(ctx, "Failed to get stored query by name and version", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Execute AQL query
	if err := h.QueryService.QueryAndCopyTo(ctx, c.Response().BodyWriter(), storedQuery.Query, queryParameters); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Stored AQL Version", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
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

	name := c.Params("qualified_query_name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("qualified_query_name path parameter is required")
	}

	version := c.Params("version")
	if version == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version path parameter is required")
	}

	var aqlRequest StoredAQLVersionRequest
	if err := c.BodyParser(&aqlRequest); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse Stored AQL Version request body", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if aqlRequest.EHRID != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Version Post with ehr_id not implemented yet")
	}

	if aqlRequest.Fetch != "" {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Version Post with fetch not implemented yet")
	}

	if aqlRequest.Offset != 0 {
		return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Version Post with offset not implemented yet")
	}

	// Retrieve stored query by name and version
	storedQuery, err := h.QueryService.GetQueryByName(ctx, name, version)
	if err != nil {
		if err == service.ErrQueryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Stored query not found for the given name and version")
		}
		h.Logger.ErrorContext(ctx, "Failed to get stored query by name and version", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Execute AQL query
	if err := h.QueryService.QueryAndCopyTo(ctx, c.Response().BodyWriter(), storedQuery.Query, aqlRequest.QueryParameters); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to execute Stored AQL Version Post", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
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

	c.Accepts("application/json")

	name := c.Params("qualified_query_name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("qualified_query_name path parameter is required")
	}

	queries, err := h.QueryService.ListStoredQueries(ctx, name)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to list stored queries", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(queries)
}

func (h *Handler) StoreQuery(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("text/plain")

	name := c.Params("qualified_query_name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("qualified_query_name path parameter is required")
	}

	queryType := c.Query("query_type")
	if queryType != "" && queryType != "AQL" {
		return c.Status(fiber.StatusBadRequest).SendString("Unsupported query_type. Only 'AQL' is supported.")
	}

	query := string(c.Body())
	if query == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Query in request body is required")
	}

	// Check if query with the same name already exists
	_, err := h.QueryService.GetQueryByName(ctx, name, "")
	if err == nil {
		return c.Status(fiber.StatusConflict).SendString("Query with the given name already exists, system cannot update without knowing the target version, please use Store Query Version endpoint instead")
	}
	if err != service.ErrQueryNotFound {
		h.Logger.ErrorContext(ctx, "Failed to check existing query by name", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to store query")
	}

	// Store the query
	err = h.QueryService.StoreQuery(ctx, name, "1.0.0", query)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to store query", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to store query")
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) StoreQueryVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("text/plain")

	name := c.Params("qualified_query_name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("qualified_query_name path parameter is required")
	}

	version := c.Params("version")
	if version == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version path parameter is required")
	}

	queryType := c.Query("query_type")
	if queryType != "" && queryType != "AQL" {
		return c.Status(fiber.StatusBadRequest).SendString("Unsupported query_type. Only 'AQL' is supported.")
	}

	query := string(c.Body())
	if query == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Query in request body is required")
	}

	// Store the new version of the query
	err := h.QueryService.StoreQuery(ctx, name, version, query)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to store query version", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to store query version")
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (h *Handler) GetStoredQueryAtVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("text/plain")

	name := c.Params("qualified_query_name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("qualified_query_name path parameter is required")
	}

	version := c.Params("version")
	if version == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version path parameter is required")
	}

	storedQuery, err := h.QueryService.GetQueryByName(ctx, name, version)
	if err != nil {
		if err == service.ErrQueryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Stored query not found for the given name and version")
		}
		h.Logger.ErrorContext(ctx, "Failed to get stored query by name and version", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(storedQuery)
}

func (h *Handler) DeleteEHRByID(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	err := h.EHRService.DeleteEHR(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR with the given ID not found")
		}

		h.Logger.ErrorContext(ctx, "Failed to delete EHR by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
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

	err := h.EHRService.DeleteEHRBulk(ctx, ehrIDList)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to delete multiple EHRs", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

type ValidationErrorResponse struct {
	Message          string   `json:"message"`
	ValidationErrors []string `json:"validationErrors"`
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
