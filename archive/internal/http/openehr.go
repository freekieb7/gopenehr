package http

import (
	"encoding/json"
	"net/http"

	"github.com/freekieb7/gopenehr/internal/storage"
)

type OpenEHRHandler struct {
	db *storage.Database
}

func NewOpenEHRHandler(db *storage.Database) OpenEHRHandler {
	return OpenEHRHandler{db: db}
}

func (h *OpenEHRHandler) ServerInfo(w http.ResponseWriter, r *http.Request) error {
	payload := map[string]any{
		"solution":              "openEHRSys",
		"solution_version":      "v1.0",
		"vendor":                "GOpenEHR",
		"restapi_specs_version": "1.0.3",
		"conformance_profile":   "CUSTOM",
		"endpoints": []string{
			// "/ehr",
			// "/demographics",
			// "/definition",
			"/query",
			// "/admin",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(payload)
}

func (h *OpenEHRHandler) ListEhr(w http.ResponseWriter, r *http.Request) error {
	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) CreateEhr(w http.ResponseWriter, r *http.Request) error {
	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetEhrById(w http.ResponseWriter, r *http.Request) error {
	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) DeleteEhrById(w http.ResponseWriter, r *http.Request) error {
	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetEhrStatusById(w http.ResponseWriter, r *http.Request) error {
	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) UpdateEhrStatusById(w http.ResponseWriter, r *http.Request) error {
	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) ListComposition(w http.ResponseWriter, r *http.Request) error {
	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) CreateComposition(w http.ResponseWriter, r *http.Request) error {
	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetCompositionById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) UpdateCompositionById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) DeleteCompositionById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) ListFolder(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) CreateFolder(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetFolderById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) UpdateFolderById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) DeleteFolderById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) CreateContribution(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetContributionById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

// Definition

func (h *OpenEHRHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) ListTemplates(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetTemplateById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) UpdateTemplateById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) DeleteTemplateById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

// Demographics

func (h *OpenEHRHandler) ListAgent(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) CreateAgent(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetAgentById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) UpdateAgentById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) DeleteAgentById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) ListGroup(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) CreateGroup(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetGroupById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) UpdateGroupById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) DeleteGroupById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) ListOrganisation(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) CreateOrganisation(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetOrganisationById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) UpdateOrganisationById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) DeleteOrganisationById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) ListPerson(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) CreatePerson(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetPersonById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) UpdatePersonById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) DeletePersonById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) ListRole(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) CreateRole(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) GetRoleById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) UpdateRoleById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}

func (h *OpenEHRHandler) DeleteRoleById(w http.ResponseWriter, r *http.Request) error {

	// todo
	w.WriteHeader(http.StatusNotImplemented)
	return nil
}
