package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const REFERENCE_RANGE_MODEL_NAME string = "REFERENCE_RANGE"

var _ util.ReferenceModel = (*REFERENCE_RANGE)(nil)

type REFERENCE_RANGE struct {
	Type_   util.Optional[string] `json:"_type,omitzero"`
	Meaning X_DV_TEXT             `json:"meaning"`
	Range   DV_INTERVAL           `json:"range"`
}

func (r REFERENCE_RANGE) HasModelName() bool {
	return r.Type_.IsSet()
}

func (r REFERENCE_RANGE) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError

	// Validate _type
	if r.Type_.IsSet() && r.Type_.Unwrap() != REFERENCE_RANGE_MODEL_NAME {
		errors = append(errors, util.ValidationError{
			Model:          REFERENCE_RANGE_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", REFERENCE_RANGE_MODEL_NAME, r.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", REFERENCE_RANGE_MODEL_NAME),
		})
	}

	// Validate meaning
	attrPath := path + ".meaning"
	errors = append(errors, r.Meaning.Validate(attrPath)...)

	// Validate range
	attrPath = path + ".range"
	errors = append(errors, r.Range.Validate(attrPath)...)

	return errors
}
