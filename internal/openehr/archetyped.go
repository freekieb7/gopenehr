package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const ARCHETYPED_MODEL_NAME string = "ARCHETYPED"

type ARCHETYPED struct {
	Type_       util.Optional[string]      `json:"_type,omitzero"`
	ArchetypeID ARCHETYPE_ID               `json:"archetype_id"`
	TemplateID  util.Optional[TEMPLATE_ID] `json:"template_id,omitzero"`
	RMVersion   string                     `json:"rm_version"`
}

func (a *ARCHETYPED) SetModelName() {
	a.Type_ = util.Some(ARCHETYPED_MODEL_NAME)
	a.ArchetypeID.SetModelName()
	if a.TemplateID.E {
		a.TemplateID.V.SetModelName()
	}
}

func (a *ARCHETYPED) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ARCHETYPED_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          ARCHETYPED_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to ARCHETYPED",
		})
	}

	// Validate archetype_id
	attrPath = path + ".archetype_id"
	errors = append(errors, a.ArchetypeID.Validate(attrPath)...)

	// Validate template_id
	if a.TemplateID.E {
		attrPath = path + ".template_id"
		errors = append(errors, a.TemplateID.V.Validate(attrPath)...)
	}

	// Validate rm_version
	if a.RMVersion == "" {
		attrPath = path + ".rm_version"
		errors = append(errors, util.ValidationError{
			Model:          ARCHETYPED_MODEL_NAME,
			Path:           attrPath,
			Message:        "rm_version field cannot be empty",
			Recommendation: "Ensure rm_version field is not empty",
		})
	}

	return errors
}
