package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ARCHETYPED_TYPE string = "ARCHETYPED"

type ARCHETYPED struct {
	Type_       utils.Optional[string]      `json:"_type,omitzero"`
	ArchetypeID ARCHETYPE_ID                `json:"archetype_id"`
	TemplateID  utils.Optional[TEMPLATE_ID] `json:"template_id,omitzero"`
	RMVersion   string                      `json:"rm_version"`
}

func (a *ARCHETYPED) SetModelName() {
	a.Type_ = utils.Some(ARCHETYPED_TYPE)
	a.ArchetypeID.SetModelName()
	if a.TemplateID.E {
		a.TemplateID.V.SetModelName()
	}
}

func (a *ARCHETYPED) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ARCHETYPED_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ARCHETYPED_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to ARCHETYPED",
		})
	}

	// Validate archetype_id
	attrPath = path + ".archetype_id"
	validateErr.Errs = append(validateErr.Errs, a.ArchetypeID.Validate(attrPath).Errs...)

	// Validate template_id
	if a.TemplateID.E {
		attrPath = path + ".template_id"
		validateErr.Errs = append(validateErr.Errs, a.TemplateID.V.Validate(attrPath).Errs...)
	}

	// Validate rm_version
	if a.RMVersion == "" {
		attrPath = path + ".rm_version"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ARCHETYPED_TYPE,
			Path:           attrPath,
			Message:        "rm_version field cannot be empty",
			Recommendation: "Ensure rm_version field is not empty",
		})
	}

	return validateErr
}
