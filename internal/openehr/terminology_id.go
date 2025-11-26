package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const TERMINOLOGY_ID_MODEL_NAME string = "TERMINOLOGY_ID"

type TERMINOLOGY_ID struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
	Value string                `json:"value"`
}

func (t *TERMINOLOGY_ID) SetModelName() {
	t.Type_ = util.Some(TERMINOLOGY_ID_MODEL_NAME)
}

func (t *TERMINOLOGY_ID) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if t.Type_.E && t.Type_.V != TERMINOLOGY_ID_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          TERMINOLOGY_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", TERMINOLOGY_ID_MODEL_NAME, t.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", TERMINOLOGY_ID_MODEL_NAME),
		})
	}

	// Validate value
	if t.Value == "" {
		attrPath = path + ".value"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          TERMINOLOGY_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field is required",
			Recommendation: "Ensure value field is not empty",
		})
	}

	return validateErr
}
