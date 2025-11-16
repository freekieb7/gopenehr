package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/openehr"
	"github.com/freekieb7/gopenehr/internal/openehr/service"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Version            string
	Logger             *slog.Logger
	EHRService         *service.EHRService
	DemographicService *service.DemographicService
}

func NewHandler(version string, logger *slog.Logger, ehrService *service.EHRService, demographicService *service.DemographicService) Handler {
	return Handler{
		Version:            version,
		Logger:             logger,
		EHRService:         ehrService,
		DemographicService: demographicService,
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
	v1.Get("/query/:qualified_query_name/version", h.ExecuteStoredAQLVersion)
	v1.Post("/query/:qualified_query_name/version", h.ExecuteStoredAQLVersionPost)

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
		"version":               h.Version,
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
		return c.Status(fiber.StatusBadRequest).SendString("subject_id query parameters are required")
	}

	subjectNamespace := c.Query("subject_namespace")
	if subjectNamespace == "" {
		return c.Status(fiber.StatusBadRequest).SendString("subject_namespace query parameters are required")
	}

	ehrJSON, err := h.EHRService.GetEHRBySubjectAsJSON(ctx, subjectID, subjectNamespace)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given subject")
		}
		h.Logger.ErrorContext(ctx, "Failed to get EHR by subject", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(ehrJSON)
}

func (h *Handler) CreateEHR(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	// Check for optional EHR_STATUS in the request body
	var newEhrStatus util.Optional[openehr.EHR_STATUS]
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&newEhrStatus.V); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}
		newEhrStatus.E = true
	}

	// Create EHR
	ehr, err := h.EHRService.CreateEHR(ctx, newEhrStatus)
	if err != nil {
		if newEhrStatus.E && err == service.ErrEHRStatusAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("EHR Status with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create EHR", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(ehr)
}

func (h *Handler) GetEHR(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	ehrJSON, err := h.EHRService.GetEHRAsJSON(ctx, ehrID)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to get EHR by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(ehrJSON)
}

func (h *Handler) CreateEHRWithID(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	// Check for optional EHR_STATUS in the request body
	var newEhrStatus util.Optional[openehr.EHR_STATUS]
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&newEhrStatus.V); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}
		newEhrStatus.E = true
	}

	// Create EHR with specified ID and EHR_STATUS
	ehr, err := h.EHRService.CreateEHRWithID(ctx, newEhrStatus, ehrID)
	if err != nil {
		if err == service.ErrEHRAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("EHR with the given ID already exists")
		}
		if newEhrStatus.E && err == service.ErrEHRStatusAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("EHR Status with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create EHR", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(ehr)
}

func (h *Handler) GetEHRStatus(c *fiber.Ctx) error {
	ctx := c.Context()

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

	ehrStatusJSON, err := h.EHRService.GetEHRStatusAsJSON(ctx, ehrID, filterAtTime, "")
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR Status not found for the given EHR ID at the specified time")
		}
		h.Logger.ErrorContext(ctx, "Failed to get EHR Status at time", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(ehrStatusJSON)
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

	var newEhrStatus openehr.EHR_STATUS
	if err := c.BodyParser(&newEhrStatus); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Check collision using If-Match header
	currentEHRStatus, err := h.EHRService.GetEHRStatus(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRStatusNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR Status not found for the given EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get current EHR Status", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing EHR Status
	if currentEHRStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("EHR Status has been modified since the provided version")
	}

	// Proceed to update EHR Status
	if err := h.EHRService.UpdateEHRStatus(ctx, ehrID, newEhrStatus); err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR Status not found for the given EHR ID")
		}
		if err == service.ErrEHRStatusAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("EHR Status with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to update EHR Status", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetEHRStatusByVersionID(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version_uid parameter is required")
	}

	ehrStatusJSON, err := h.EHRService.GetEHRStatusAsJSON(ctx, ehrID, time.Time{}, versionUID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR Status not found for the given EHR ID and version UID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get EHR Status at time", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(ehrStatusJSON)
}

func (h *Handler) GetVersionedEHRStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionedStatusJSON, err := h.EHRService.GetVersionedEHRStatusAsJSON(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned EHR Status not found for the given EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(versionedStatusJSON)
}

func (h *Handler) GetVersionedEHRStatusRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	revisionHistoryJSON, err := h.EHRService.GetVersionedEHRStatusRevisionHistoryAsJSON(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned EHR Status revision history not found for the given EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned EHR Status revision history", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(revisionHistoryJSON)
}

func (h *Handler) GetVersionedEHRStatusVersion(c *fiber.Ctx) error {
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

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
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

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	var newComposition openehr.COMPOSITION
	if err := c.BodyParser(&newComposition); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	composition, err := h.EHRService.CreateComposition(ctx, ehrID, newComposition)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Composition with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create Composition", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(composition)
}

func (h *Handler) GetComposition(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	compositionJSON, err := h.EHRService.GetCompositionAsJSON(ctx, ehrID, uidBasedID)
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

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(compositionJSON)
}

func (h *Handler) UpdateComposition(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// Check collision using If-Match header
	currentComposition, err := h.EHRService.GetComposition(ctx, uidBasedID, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Composition not found for the given UID")
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
	var updatedComposition openehr.COMPOSITION
	if err := c.BodyParser(&updatedComposition); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Composition
	if err := h.EHRService.UpdateCompositionByID(ctx, ehrID, updatedComposition); err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrCompositionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Composition not found for the given UID")
		}
		if err == service.ErrCompositionAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Composition with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to update Composition", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) DeleteComposition(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if !strings.Contains(uidBasedID, "::") {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete Composition by versioned object ID. Please provide the object version ID.")
	}

	// Check existence before deletion
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentComposition, err := h.EHRService.GetComposition(ctx, versionedObjectID, ehrID)
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
	if err := h.EHRService.DeleteComposition(ctx, versionedObjectID); err != nil {
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

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetVersionedCompositionByID(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}

	versionedCompositionJSON, err := h.EHRService.GetVersionedCompositionByIDAsJSON(ctx, ehrID, versionedObjectID)
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

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(versionedCompositionJSON)
}

func (h *Handler) GetVersionedCompositionRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}

	revisionHistoryJSON, err := h.EHRService.GetVersionedCompositionRevisionHistoryAsJSON(ctx, ehrID, versionedObjectID)
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

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(revisionHistoryJSON)
}

func (h *Handler) GetVersionedCompositionVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()

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

	versionAsJSON, err := h.EHRService.GetVersionedCompositionVersionAsJSON(ctx, ehrID, versionedObjectID, filterAtTime, "")
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
	return c.Status(fiber.StatusOK).Send(versionAsJSON)
}

func (h *Handler) GetVersionedCompositionVersionByID(c *fiber.Ctx) error {
	ctx := c.Context()

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

	versionAsJSON, err := h.EHRService.GetVersionedCompositionVersionAsJSON(ctx, ehrID, versionedObjectID, time.Time{}, versionUID)
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
	return c.Status(fiber.StatusOK).Send(versionAsJSON)
}

func (h *Handler) CreateDirectory(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	directory, err := h.EHRService.CreateDirectory(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrDirectoryAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Directory with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create Directory", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(directory)
}

func (h *Handler) UpdateDirectory(c *fiber.Ctx) error {
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

	var currentDirectory openehr.FOLDER
	rawCurrentDirectory, err := h.EHRService.GetRawDirectory(ctx, ehrID)
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
	if err := json.Unmarshal(rawCurrentDirectory, &currentDirectory); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to unmarshal current Directory", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	// Safe to assume that UID and OBJECT_VERSION_ID are always set for existing Directory
	if currentDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value != ifMatch {
		return c.Status(fiber.StatusPreconditionFailed).SendString("Directory has been modified since the provided version")
	}

	// Parse updated directory from request body
	var updatedDirectory openehr.FOLDER
	if err := c.BodyParser(&updatedDirectory); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Check collision using If-Match header

	if err := h.EHRService.UpdateDirectory(ctx, ehrID, updatedDirectory); err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given EHR ID")
		}
		if err == service.ErrDirectoryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory not found for the given EHR ID")
		}
		if err == service.ErrDirectoryAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Directory with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to update Directory", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) DeleteDirectory(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	var currentDirectory openehr.FOLDER
	rawCurrentDirectory, err := h.EHRService.GetRawDirectory(ctx, ehrID)
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
	if err := json.Unmarshal(rawCurrentDirectory, &currentDirectory); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to unmarshal current Directory", "error", err)
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

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetFolderInDirectoryVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()

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

	folderAsJSON, err := h.EHRService.GetFolderInDirectoryVersionAsJSON(ctx, ehrID, filterAtTime, "")
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory not found for the given EHR ID at the specified time")
		}
		if err == service.ErrDirectoryNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory version not found at the specified time for the given EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Folder in Directory version at time", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(folderAsJSON)
}

func (h *Handler) GetFolderInDirectoryVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	versionUID := c.Params("version_uid")
	if versionUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version_uid parameter is required")
	}

	path := c.Query("path")
	if path != "" {
		return c.Status(fiber.StatusBadRequest).SendString("path query parameter is not supported")
	}

	folderAsJSON, err := h.EHRService.GetFolderInDirectoryVersionAsJSON(ctx, ehrID, time.Time{}, versionUID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Directory not found for the given EHR ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Folder in Directory version", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(folderAsJSON)
}

func (h *Handler) CreateContribution(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	var newContribution openehr.CONTRIBUTION
	if err := c.BodyParser(&newContribution); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	contribution, err := h.EHRService.CreateContribution(ctx, ehrID, newContribution)
	if err != nil {
		if err == service.ErrContributionAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Contribution with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create Contribution", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(contribution)
}

func (h *Handler) GetContribution(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	ehrID := c.Params("ehr_id")
	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	contributionUID := c.Params("contribution_uid")
	if contributionUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("contribution_uid parameter is required")
	}

	contributionAsJSON, err := h.EHRService.GetContributionAsJSON(ctx, ehrID, contributionUID)
	if err != nil {
		if err == service.ErrContributionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Contribution not found for the given EHR ID and Contribution UID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Contribution by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(contributionAsJSON)
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

	var newAgent openehr.AGENT
	if err := c.BodyParser(&newAgent); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Agent", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if errs := newAgent.Validate("$"); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	agent, err := h.DemographicService.CreateAgent(ctx, newAgent)
	if err != nil {
		if err == service.ErrAgentAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Agent with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create Agent", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(agent)
}

func (h *Handler) GetAgent(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	agentAsJSON, err := h.DemographicService.GetAgentAsJSON(ctx, uidBasedID)
	if err != nil {
		if err == service.ErrAgentNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Agent not found for the given agent ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Agent by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(agentAsJSON)
}

func (h *Handler) UpdateAgent(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// Ensure Agent exists before update
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentAgent, err := h.DemographicService.GetAgent(ctx, versionedObjectID)
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
	var updatedAgent openehr.AGENT
	if err := c.BodyParser(&updatedAgent); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Agent", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Agent
	if err := h.DemographicService.UpdateAgent(ctx, updatedAgent); err != nil {
		if err == service.ErrAgentNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Agent not found for the given agent ID")
		}
		if err == service.ErrAgentAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Agent with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to update Agent", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) DeleteAgent(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if !strings.Contains(uidBasedID, "::") {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete Agent by versioned object ID. Please provide the object version ID.")
	}

	// Check existence before deletion
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentAgent, err := h.DemographicService.GetAgent(ctx, versionedObjectID)
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

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) CreateGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	var newGroup openehr.GROUP
	if err := c.BodyParser(&newGroup); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Group", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if errs := newGroup.Validate("$"); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	group, err := h.DemographicService.CreateGroup(ctx, newGroup)
	if err != nil {
		if err == service.ErrGroupAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Group with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create Group", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}

func (h *Handler) GetGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	groupAsJSON, err := h.DemographicService.GetGroupAsJSON(ctx, uidBasedID)
	if err != nil {
		if err == service.ErrGroupNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Group not found for the given group ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Group by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(groupAsJSON)
}

func (h *Handler) UpdateGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// Ensure Group exists before update
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentGroup, err := h.DemographicService.GetGroup(ctx, versionedObjectID)
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
	var updatedGroup openehr.GROUP
	if err := c.BodyParser(&updatedGroup); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Group", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Group
	if err := h.DemographicService.UpdateGroup(ctx, updatedGroup); err != nil {
		if err == service.ErrGroupNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Group not found for the given group ID")
		}
		if err == service.ErrGroupAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Group with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to update Group", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.SendStatus(fiber.StatusNoContent)
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
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentGroup, err := h.DemographicService.GetGroup(ctx, versionedObjectID)
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

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) CreatePerson(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	var newPerson openehr.PERSON
	if err := c.BodyParser(&newPerson); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Person", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if errs := newPerson.Validate("$"); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	person, err := h.DemographicService.CreatePerson(ctx, newPerson)
	if err != nil {
		if err == service.ErrPersonAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Person with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create Person", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(person)
}

func (h *Handler) GetPerson(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	personAsJSON, err := h.DemographicService.GetPersonAsJSON(ctx, uidBasedID)
	if err != nil {
		if err == service.ErrPersonNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Person not found for the given person ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Person by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(personAsJSON)
}

func (h *Handler) UpdatePerson(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// Ensure Person exists before update
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentPerson, err := h.DemographicService.GetPerson(ctx, versionedObjectID)
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
	var updatedPerson openehr.PERSON
	if err := c.BodyParser(&updatedPerson); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Person", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Person
	if err := h.DemographicService.UpdatePerson(ctx, updatedPerson); err != nil {
		if err == service.ErrPersonNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Person not found for the given person ID")
		}
		if err == service.ErrPersonAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Person with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to update Person", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) DeletePerson(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if !strings.Contains(uidBasedID, "::") {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete Person by versioned object ID. Please provide the object version ID.")
	}

	// Check existence before deletion
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentPerson, err := h.DemographicService.GetPerson(ctx, versionedObjectID)
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

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) CreateOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	var newOrganisation openehr.ORGANISATION
	if err := c.BodyParser(&newOrganisation); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Organisation", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if errs := newOrganisation.Validate("$"); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	organisation, err := h.DemographicService.CreateOrganisation(ctx, newOrganisation)
	if err != nil {
		if err == service.ErrOrganisationAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Organisation with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create Organisation", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(organisation)
}

func (h *Handler) GetOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	organisationAsJSON, err := h.DemographicService.GetOrganisationAsJSON(ctx, uidBasedID)
	if err != nil {
		if err == service.ErrOrganisationNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Organisation not found for the given organisation ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Organisation by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(organisationAsJSON)
}

func (h *Handler) UpdateOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// Ensure Organisation exists before update
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentOrganisation, err := h.DemographicService.GetOrganisation(ctx, versionedObjectID)
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
	var updatedOrganisation openehr.ORGANISATION
	if err := c.BodyParser(&updatedOrganisation); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Organisation", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Organisation
	if err := h.DemographicService.UpdateOrganisation(ctx, updatedOrganisation); err != nil {
		if err == service.ErrOrganisationNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Organisation not found for the given organisation ID")
		}
		if err == service.ErrOrganisationAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Organisation with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to update Organisation", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) DeleteOrganisation(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	if !strings.Contains(uidBasedID, "::") {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete Organisation by versioned object ID. Please provide the object version ID.")
	}

	// Check existence before deletion
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentOrganisation, err := h.DemographicService.GetOrganisation(ctx, versionedObjectID)
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

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) CreateRole(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	var newRole openehr.ROLE
	if err := c.BodyParser(&newRole); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Role", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if errs := newRole.Validate("$"); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	role, err := h.DemographicService.CreateRole(ctx, newRole)
	if err != nil {
		if err == service.ErrRoleAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Role with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create Role", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(role)
}

func (h *Handler) GetRole(c *fiber.Ctx) error {
	ctx := c.Context()

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	roleAsJSON, err := h.DemographicService.GetRoleAsJSON(ctx, uidBasedID)
	if err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Role not found for the given role ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Role by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(roleAsJSON)
}

func (h *Handler) UpdateRole(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	uidBasedID := c.Params("uid_based_id")
	if uidBasedID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("uid_based_id parameter is required")
	}

	ifMatch := c.Get("If-Match")
	if ifMatch == "" {
		return c.Status(fiber.StatusBadRequest).SendString("If-Match header is required")
	}

	// Ensure Role exists before update
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentRole, err := h.DemographicService.GetRole(ctx, versionedObjectID)
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
	var updatedRole openehr.ROLE
	if err := c.BodyParser(&updatedRole); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for updated Role", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Proceed to update Role
	if err := h.DemographicService.UpdateRole(ctx, updatedRole); err != nil {
		if err == service.ErrRoleNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Role not found for the given role ID")
		}
		if err == service.ErrRoleAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Role with the given UID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to update Role", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.SendStatus(fiber.StatusNoContent)
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
	versionedObjectID := strings.Split(uidBasedID, "::")[0]
	currentRole, err := h.DemographicService.GetRole(ctx, versionedObjectID)
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

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetVersionedParty(c *fiber.Ctx) error {
	ctx := c.Context()

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}

	partyAsJSON, err := h.DemographicService.GetVersionedPartyAsJSON(ctx, versionedObjectID)
	if err != nil {
		if err == service.ErrVersionedPartyNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Party not found for the given ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(partyAsJSON)
}

func (h *Handler) GetVersionedPartyRevisionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}

	historyAsJSON, err := h.DemographicService.GetVersionedPartyRevisionHistoryAsJSON(ctx, versionedObjectID)
	if err != nil {
		if err == service.ErrVersionedPartyNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Party not found for the given ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party Revision History by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(historyAsJSON)
}

func (h *Handler) GetVersionedPartyVersionAtTime(c *fiber.Ctx) error {
	ctx := c.Context()

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

	partyAsJSON, err := h.DemographicService.GetVersionedPartyVersionAsJSON(ctx, versionedObjectID, filterAtTime, "")
	if err != nil {
		if err == service.ErrVersionedPartyNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Party not found for the given ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party Version at Time by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(partyAsJSON)
}

func (h *Handler) GetVersionedPartyVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	versionedObjectID := c.Params("versioned_object_id")
	if versionedObjectID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("versioned_object_id parameter is required")
	}

	versionID := c.Params("version_id")
	if versionID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("version_id parameter is required")
	}

	partyAsJSON, err := h.DemographicService.GetVersionedPartyVersionAsJSON(ctx, versionedObjectID, time.Time{}, versionID)
	if err != nil {
		if err == service.ErrVersionedPartyNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Versioned Party not found for the given ID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Versioned Party Version by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(partyAsJSON)
}

func (h *Handler) CreateDemographicContribution(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	var newContribution openehr.CONTRIBUTION
	if err := c.BodyParser(&newContribution); err != nil {
		h.Logger.ErrorContext(ctx, "Failed to parse request body for new Demographic Contribution", "error", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if errs := newContribution.Validate("$"); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	contribution, err := h.DemographicService.CreateContribution(ctx, newContribution)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to create Demographic Contribution", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(contribution)
}

func (h *Handler) GetDemographicContribution(c *fiber.Ctx) error {
	ctx := c.Context()

	c.Accepts("application/json")

	contributionUID := c.Params("contribution_uid")
	if contributionUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("contribution_uid parameter is required")
	}

	contributionAsJSON, err := h.DemographicService.GetContributionAsJSON(ctx, contributionUID)
	if err != nil {
		if err == service.ErrContributionNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Contribution not found for the given EHR ID and Contribution UID")
		}
		h.Logger.ErrorContext(ctx, "Failed to get Contribution by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(contributionAsJSON)
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
	return c.Status(fiber.StatusNotImplemented).SendString("Execute Ad Hoc AQL not implemented yet")
}

func (h *Handler) ExecuteAdHocAQLPost(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Execute Ad Hoc AQL Post not implemented yet")
}

func (h *Handler) ExecuteStoredAQL(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL not implemented yet")
}

func (h *Handler) ExecuteStoredAQLPost(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Post not implemented yet")
}

func (h *Handler) ExecuteStoredAQLVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Version not implemented yet")
}

func (h *Handler) ExecuteStoredAQLVersionPost(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Execute Stored AQL Version Post not implemented yet")
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
	return c.Status(fiber.StatusNotImplemented).SendString("List Stored Queries not implemented yet")
}

func (h *Handler) StoreQuery(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Store Query not implemented yet")
}

func (h *Handler) StoreQueryVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Store Query Version not implemented yet")
}

func (h *Handler) GetStoredQueryAtVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Get Stored Query At Version not implemented yet")
}

func (h *Handler) DeleteEHRByID(c *fiber.Ctx) error {
	ctx := c.Context()
	ehrID := c.Params("ehr_id")

	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	err := h.EHRService.DeleteEHRByID(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR with the given ID not found")
		}

		h.Logger.ErrorContext(ctx, "Failed to delete EHR by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) DeleteMultipleEHRs(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse multiple ehr_id query parameters
	// Example: ?ehr_id=7d44b88c-4199-4bad-97dc-d78268e01398&ehr_id=297c3e91-7c17-4497-85dd-01e05aaae44e
	var ehrIDList []string
	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		if string(key) == "ehr_id" {
			ehrIDList = append(ehrIDList, string(value))
		}
	})

	err := h.EHRService.DeleteMultipleEHRs(ctx, ehrIDList)
	if err != nil {
		h.Logger.ErrorContext(ctx, "Failed to delete multiple EHRs", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	return c.SendStatus(fiber.StatusNoContent)
}

type ValidationErrorResponse struct {
	Message          string   `json:"message"`
	ValidationErrors []string `json:"validationErrors"`
}
