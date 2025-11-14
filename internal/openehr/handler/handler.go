package handler

import (
	"log/slog"
	"net/http"

	"github.com/freekieb7/gopenehr/internal/openehr/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Version    string
	Logger     *slog.Logger
	EHRService *service.EHR
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/openehr/v1")

	v1.Options("", h.SystemInfo)

	v1.Get("/ehr", h.GetEHRBySubject)
	v1.Post("/ehr", h.CreateEHR)
	v1.Get("/ehr/:ehr_id", h.GetEHRByID)
	v1.Put("/ehr/:ehr_id", h.CreateEHRWithID)

	v1.Get("/ehr/:ehr_id/ehr_status", h.GetEHRStatusAtTime)
	v1.Put("/ehr/:ehr_id/ehr_status", h.UpdateEhrStatus)
	v1.Get("/ehr/:ehr_id/ehr_status/:version_uid", h.GetEHRStatusByVersionID)

	v1.Get("/ehr/:ehr_id/versioned_ehr_status", h.GetVersionedEHRStatus)
	v1.Get("/ehr/:ehr_id/versioned_ehr_status/revision_history", h.GetVersionedEHRStatusRevisionHistory)
	v1.Get("/ehr/:ehr_id/versioned_ehr_status/version", h.GetVersionedEHRStatusVersionAtTime)
	v1.Get("/ehr/:ehr_id/versioned_ehr_status/version/:version_uid", h.GetVersionedEHRStatusVersionByID)

	v1.Post("/ehr/:ehr_id/composition", h.CreateComposition)
	v1.Get("/ehr/:ehr_id/composition/:uid_based_id", h.GetCompositionByID)
	v1.Put("/ehr/:ehr_id/composition/:uid_based_id", h.UpdateCompositionByID)
	v1.Delete("/ehr/:ehr_id/composition/:uid_based_id", h.DeleteCompositionByID)

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
	v1.Get("/ehr/:ehr_id/contribution/:contribution_uid", h.GetContributionByID)

	v1.Get("/ehr/:ehr_id/tags", h.GetEHRTags)
	v1.Get("/ehr/:ehr_id/composition/:uid_based_id/tags", h.GetCompositionTags)
	v1.Put("/ehr/:ehr_id/composition/:uid_based_id/tags", h.UpdateCompositionTags)
	v1.Delete("/ehr/:ehr_id/composition/:uid_based_id/tags", h.DeleteCompositionTagByKey)
	v1.Get("/ehr/:ehr_id/ehr_status/tags", h.GetEHRStatusTags)
	v1.Get("/ehr/:ehr_id/ehr_status/:version_uid/tags", h.GetEHRStatusVersionTags)
	v1.Put("/ehr/:ehr_id/ehr_status/:version_uid/tags", h.UpdateEHRStatusVersionTags)
	v1.Delete("/ehr/:ehr_id/ehr_status/:version_uid/tags/:key", h.DeleteEHRStatusVersionTagByKey)

	v1.Post("/demographic/agent", h.CreateAgent)
	v1.Get("/demographic/agent/:uid_based_id", h.GetAgentByID)
	v1.Put("/demographic/agent/:uid_based_id", h.UpdateAgentByID)
	v1.Delete("/demographic/agent/:uid_based_id", h.DeleteAgentByID)

	v1.Post("/demographic/group", h.CreateGroup)
	v1.Get("/demographic/group/:uid_based_id", h.GetGroupByID)
	v1.Put("/demographic/group/:uid_based_id", h.UpdateGroupByID)
	v1.Delete("/demographic/group/:uid_based_id", h.DeleteGroupByID)

	v1.Post("/demographic/person", h.CreatePerson)
	v1.Get("/demographic/person/:uid_based_id", h.GetPersonByID)
	v1.Put("/demographic/person/:uid_based_id", h.UpdatePersonByID)
	v1.Delete("/demographic/person/:uid_based_id", h.DeletePersonByID)

	v1.Post("/demographic/organisation", h.CreateOrganisation)
	v1.Get("/demographic/organisation/:uid_based_id", h.GetOrganisationByID)
	v1.Put("/demographic/organisation/:uid_based_id", h.UpdateOrganisationByID)
	v1.Delete("/demographic/organisation/:uid_based_id", h.DeleteOrganisationByID)

	v1.Post("/demographic/role", h.CreateRole)
	v1.Get("/demographic/role/:uid_based_id", h.GetRoleByID)
	v1.Put("/demographic/role/:uid_based_id", h.UpdateRoleByID)
	v1.Delete("/demographic/role/:uid_based_id", h.DeleteRoleByID)

	v1.Get("/demographic/versioned_party/:versioned_object_uid", h.GetVersionedPartyByID)
	v1.Get("/demographic/versioned_party/:versioned_object_uid/revision_history", h.GetVersionedPartyRevisionHistory)
	v1.Get("/demographic/versioned_party/:versioned_object_uid/version", h.GetVersionedPartyVersionAtTime)
	v1.Get("/demographic/versioned_party/:versioned_object_uid/version/:version_uid", h.GetVersionedPartyVersionByID)

	v1.Post("/demographic/contribution", h.CreateDemographicContribution)
	v1.Get("/demographic/contribution/:contribution_uid", h.GetDemographicContributionByID)

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
	subjectID := c.Query("subject_id")
	subjectNamespace := c.Query("subject_namespace")

	if subjectID == "" || subjectNamespace == "" {
		return c.Status(fiber.StatusBadRequest).SendString("subject_id and subject_namespace query parameters are required")
	}

	ehr, err := h.EHRService.GetEHRBySubject(c.Context(), subjectID, subjectNamespace)
	if err != nil {
		if err == service.ErrEHRNotFound {
			return c.Status(fiber.StatusNotFound).SendString("EHR not found for the given subject")
		}
		h.Logger.ErrorContext(c.Context(), "Failed to get EHR by subject", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(ehr)
}

func (h *Handler) CreateEHR(c *fiber.Ctx) error {
	ehr, err := h.EHRService.CreateEHR(c.Context())
	if err != nil {
		h.Logger.ErrorContext(c.Context(), "Failed to create EHR", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(ehr)
}

func (h *Handler) GetEHRByID(c *fiber.Ctx) error {
	ehrID := c.Params("ehr_id")

	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	ehr, err := h.EHRService.GetEHRByID(c.Context(), ehrID)
	if err != nil {
		h.Logger.ErrorContext(c.Context(), "Failed to get EHR by ID", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(ehr)
}

func (h *Handler) CreateEHRWithID(c *fiber.Ctx) error {
	ctx := c.Context()
	ehrID := c.Params("ehr_id")

	if ehrID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ehr_id parameter is required")
	}

	ehr, err := h.EHRService.CreateEHRWithID(ctx, ehrID)
	if err != nil {
		if err == service.ErrEHRAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("EHR with the given ID already exists")
		}

		h.Logger.ErrorContext(ctx, "Failed to create EHR", "error", err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(ehr)
}

func (h *Handler) GetEHRStatusAtTime(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetEHRStatusAtTime not implemented yet")
}

func (h *Handler) UpdateEhrStatus(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateEhrStatus not implemented yet")
}

func (h *Handler) GetEHRStatusByVersionID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetEHRStatusByVersionID not implemented yet")
}

func (h *Handler) GetVersionedEHRStatus(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedEHRStatus not implemented yet")
}

func (h *Handler) GetVersionedEHRStatusRevisionHistory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedEHRStatusRevisionHistory not implemented yet")
}

func (h *Handler) GetVersionedEHRStatusVersionAtTime(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedEHRStatusVersionAtTime not implemented yet")
}

func (h *Handler) GetVersionedEHRStatusVersionByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedEHRStatusVersionByID not implemented yet")
}

func (h *Handler) CreateComposition(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("CreateComposition not implemented yet")
}

func (h *Handler) GetCompositionByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetCompositionByID not implemented yet")
}

func (h *Handler) UpdateCompositionByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateCompositionByID not implemented yet")
}

func (h *Handler) DeleteCompositionByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteCompositionByID not implemented yet")
}

func (h *Handler) GetVersionedCompositionByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedCompositionByID not implemented yet")
}

func (h *Handler) GetVersionedCompositionRevisionHistory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedCompositionRevisionHistory not implemented yet")
}

func (h *Handler) GetVersionedCompositionVersionAtTime(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedCompositionVersionAtTime not implemented yet")
}

func (h *Handler) GetVersionedCompositionVersionByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedCompositionVersionByID not implemented yet")
}

func (h *Handler) CreateDirectory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("CreateDirectory not implemented yet")
}

func (h *Handler) UpdateDirectory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateDirectory not implemented yet")
}

func (h *Handler) DeleteDirectory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteDirectory not implemented yet")
}

func (h *Handler) GetFolderInDirectoryVersionAtTime(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetFolderInDirectoryVersionAtTime not implemented yet")
}

func (h *Handler) GetFolderInDirectoryVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetFolderInDirectoryVersion not implemented yet")
}

func (h *Handler) CreateContribution(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("CreateContribution not implemented yet")
}

func (h *Handler) GetContributionByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetContributionByID not implemented yet")
}

func (h *Handler) GetEHRTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetEHRTags not implemented yet")
}

func (h *Handler) GetCompositionTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetCompositionTags not implemented yet")
}

func (h *Handler) UpdateCompositionTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateCompositionTags not implemented yet")
}

func (h *Handler) DeleteCompositionTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteCompositionTagByKey not implemented yet")
}

func (h *Handler) GetEHRStatusTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetEHRStatusTags not implemented yet")
}

func (h *Handler) GetEHRStatusVersionTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetEHRStatusVersionTags not implemented yet")
}

func (h *Handler) UpdateEHRStatusVersionTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateEHRStatusVersionTags not implemented yet")
}

func (h *Handler) DeleteEHRStatusVersionTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteEHRStatusVersionTagByKey not implemented yet")
}

func (h *Handler) CreateAgent(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("CreateAgent not implemented yet")
}

func (h *Handler) GetAgentByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetAgentByID not implemented yet")
}

func (h *Handler) UpdateAgentByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateAgentByID not implemented yet")
}

func (h *Handler) DeleteAgentByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteAgentByID not implemented yet")
}

func (h *Handler) CreateGroup(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("CreateGroup not implemented yet")
}

func (h *Handler) GetGroupByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetGroupByID not implemented yet")
}

func (h *Handler) UpdateGroupByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateGroupByID not implemented yet")
}

func (h *Handler) DeleteGroupByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteGroupByID not implemented yet")
}

func (h *Handler) CreatePerson(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("CreatePerson not implemented yet")
}

func (h *Handler) GetPersonByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetPersonByID not implemented yet")
}

func (h *Handler) UpdatePersonByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdatePersonByID not implemented yet")
}

func (h *Handler) DeletePersonByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeletePersonByID not implemented yet")
}

func (h *Handler) CreateOrganisation(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("CreateOrganisation not implemented yet")
}

func (h *Handler) GetOrganisationByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetOrganisationByID not implemented yet")
}

func (h *Handler) UpdateOrganisationByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateOrganisationByID not implemented yet")
}

func (h *Handler) DeleteOrganisationByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteOrganisationByID not implemented yet")
}
func (h *Handler) CreateRole(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("CreateRole not implemented yet")
}

func (h *Handler) GetRoleByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetRoleByID not implemented yet")
}

func (h *Handler) UpdateRoleByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateRoleByID not implemented yet")
}

func (h *Handler) DeleteRoleByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteRoleByID not implemented yet")
}

func (h *Handler) GetVersionedPartyByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedPartyByID not implemented yet")
}

func (h *Handler) GetVersionedPartyRevisionHistory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedPartyRevisionHistory not implemented yet")
}

func (h *Handler) GetVersionedPartyVersionAtTime(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedPartyVersionAtTime not implemented yet")
}

func (h *Handler) GetVersionedPartyVersionByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetVersionedPartyVersionByID not implemented yet")
}

func (h *Handler) CreateDemographicContribution(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("CreateDemographicContribution not implemented yet")
}

func (h *Handler) GetDemographicContributionByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetDemographicContributionByID not implemented yet")
}

func (h *Handler) GetDemographicTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetDemographicTags not implemented yet")
}

func (h *Handler) GetAgentTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetAgentTags not implemented yet")
}

func (h *Handler) UpdateAgentTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateAgentTags not implemented yet")
}

func (h *Handler) DeleteAgentTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteAgentTagByKey not implemented yet")
}

func (h *Handler) GetGroupTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetGroupTags not implemented yet")
}

func (h *Handler) UpdateGroupTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateGroupTags not implemented yet")
}

func (h *Handler) DeleteGroupTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteGroupTagByKey not implemented yet")
}

func (h *Handler) GetPersonTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetPersonTags not implemented yet")
}

func (h *Handler) UpdatePersonTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdatePersonTags not implemented yet")
}

func (h *Handler) DeletePersonTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeletePersonTagByKey not implemented yet")
}

func (h *Handler) GetOrganisationTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetOrganisationTags not implemented yet")
}

func (h *Handler) UpdateOrganisationTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateOrganisationTags not implemented yet")
}

func (h *Handler) DeleteOrganisationTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteOrganisationTagByKey not implemented yet")
}

func (h *Handler) GetRoleTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetRoleTags not implemented yet")
}

func (h *Handler) UpdateRoleTags(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateRoleTags not implemented yet")
}

func (h *Handler) DeleteRoleTagByKey(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteRoleTagByKey not implemented yet")
}

func (h *Handler) ExecuteAdHocAQL(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("ExecuteAdHocAQL not implemented yet")
}

func (h *Handler) ExecuteAdHocAQLPost(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("ExecuteAdHocAQLPost not implemented yet")
}

func (h *Handler) ExecuteStoredAQL(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("ExecuteStoredAQL not implemented yet")
}

func (h *Handler) ExecuteStoredAQLPost(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("ExecuteStoredAQLPost not implemented yet")
}

func (h *Handler) ExecuteStoredAQLVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("ExecuteStoredAQLVersion not implemented yet")
}

func (h *Handler) ExecuteStoredAQLVersionPost(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("ExecuteStoredAQLVersionPost not implemented yet")
}

func (h *Handler) GetTemplatesADL14(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetTemplatesADL14 not implemented yet")
}

func (h *Handler) UploadTemplateADL14(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UploadTemplateADL14 not implemented yet")
}

func (h *Handler) GetTemplateADL14ByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetTemplateADL14ByID not implemented yet")
}

func (h *Handler) GetTemplatesADL2(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetTemplatesADL2 not implemented yet")
}

func (h *Handler) UploadTemplateADL2(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("UploadTemplateADL2 not implemented yet")
}

func (h *Handler) GetTemplateADL2ByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetTemplateADL2ByID not implemented yet")
}

func (h *Handler) GetTemplateADL2AtVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetTemplateADL2AtVersion not implemented yet")
}

func (h *Handler) ListStoredQueries(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("ListStoredQueries not implemented yet")
}

func (h *Handler) StoreQuery(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("StoreQuery not implemented yet")
}

func (h *Handler) StoreQueryVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("StoreQueryVersion not implemented yet")
}

func (h *Handler) GetStoredQueryAtVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("GetStoredQueryAtVersion not implemented yet")
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
