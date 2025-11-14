package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const GENERIC_ID_MODEL_NAME string = "GENERIC_ID"

var _ util.ReferenceModel = (*GENERIC_ID)(nil)

type GENERIC_ID struct {
	Type_  util.Optional[string] `json:"_type,omitzero"`
	Value  string                `json:"value"`
	Scheme string                `json:"scheme"`
}

func (g GENERIC_ID) HasModelName() bool {
	return g.Type_.IsSet()
}

func (g GENERIC_ID) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if g.Type_.IsSet() && g.Type_.Unwrap() != GENERIC_ID_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          GENERIC_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", GENERIC_ID_MODEL_NAME, g.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", GENERIC_ID_MODEL_NAME),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if g.Value == "" {
		errors = append(errors, util.ValidationError{
			Model:          GENERIC_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	}

	// Validate scheme
	attrPath = path + ".scheme"
	if g.Scheme == "" {
		errors = append(errors, util.ValidationError{
			Model:          GENERIC_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "scheme field cannot be empty",
			Recommendation: "Ensure scheme field is not empty",
		})
	}

	return errors
}
