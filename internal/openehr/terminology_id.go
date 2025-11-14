package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const TERMINOLOGY_ID_MODEL_NAME string = "TERMINOLOGY_ID"

var _ util.ReferenceModel = (*TERMINOLOGY_ID)(nil)

type TERMINOLOGY_ID struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
	Value string                `json:"value"`
}

func (t TERMINOLOGY_ID) HasModelName() bool {
	return t.Type_.IsSet()
}

func (t TERMINOLOGY_ID) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError

	// Validate _type
	if t.Type_.IsSet() && t.Type_.Unwrap() != TERMINOLOGY_ID_MODEL_NAME {
		errors = append(errors, util.ValidationError{
			Model:          TERMINOLOGY_ID_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", TERMINOLOGY_ID_MODEL_NAME, t.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", TERMINOLOGY_ID_MODEL_NAME),
		})
	}

	// Validate value
	attrPath := path + ".value"
	if t.Value == "" {
		errors = append(errors, util.ValidationError{
			Model:          TERMINOLOGY_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field is required",
			Recommendation: "Ensure value field is not empty",
		})
	}

	return errors
}
