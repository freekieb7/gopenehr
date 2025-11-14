package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const TEMPLATE_ID_MODEL_NAME string = "TEMPLATE_ID"

type TEMPLATE_ID struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
	Value string                `json:"value"`
}

func (t TEMPLATE_ID) isObjectIDModel() {}

func (t TEMPLATE_ID) HasModelName() bool {
	return t.Type_.E
}

func (t *TEMPLATE_ID) SetModelName() {
	t.Type_ = util.Some(TEMPLATE_ID_MODEL_NAME)
}

func (t TEMPLATE_ID) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if t.Type_.E && t.Type_.V != TEMPLATE_ID_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          TEMPLATE_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", TEMPLATE_ID_MODEL_NAME, t.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", TEMPLATE_ID_MODEL_NAME),
		})
	}

	// Validate value
	if t.Value == "" {
		attrPath = path + ".value"
		errors = append(errors, util.ValidationError{
			Model:          TEMPLATE_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	}

	return errors
}
