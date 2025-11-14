package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const DV_IDENTIFIER_MODEL_NAME string = "DV_IDENTIFIER"

type DV_IDENTIFIER struct {
	Type_    util.Optional[string] `json:"_type,omitzero"`
	Issuer   util.Optional[string] `json:"issuer,omitzero"`
	Assigner util.Optional[string] `json:"assigner,omitzero"`
	ID       string                `json:"id"`
	Type     util.Optional[string] `json:"type,omitzero"`
}

func (d DV_IDENTIFIER) isDataValueModel() {}

func (d DV_IDENTIFIER) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_IDENTIFIER) SetModelName() {
	d.Type_ = util.Some(DV_IDENTIFIER_MODEL_NAME)
}

func (d DV_IDENTIFIER) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_IDENTIFIER_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          DV_IDENTIFIER_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_IDENTIFIER",
		})
	}

	// Validate id
	attrPath = path + ".id"
	if d.ID == "" {
		errors = append(errors, util.ValidationError{
			Model:          DV_IDENTIFIER_MODEL_NAME,
			Path:           attrPath,
			Message:        "id field cannot be empty",
			Recommendation: "Ensure id field is not empty",
		})
	}

	return errors
}
