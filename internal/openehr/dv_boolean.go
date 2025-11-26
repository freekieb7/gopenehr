package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const DV_BOOLEAN_MODEL_NAME string = "DV_BOOLEAN"

type DV_BOOLEAN struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
	Value bool                  `json:"value"`
}

func (d *DV_BOOLEAN) isDataValueModel() {}

func (d *DV_BOOLEAN) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_BOOLEAN) SetModelName() {
	d.Type_ = util.Some(DV_BOOLEAN_MODEL_NAME)
}

func (d *DV_BOOLEAN) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_BOOLEAN_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_BOOLEAN_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_BOOLEAN",
		})
	}

	return validateErr
}
