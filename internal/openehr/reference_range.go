package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const REFERENCE_RANGE_MODEL_NAME string = "REFERENCE_RANGE"

type REFERENCE_RANGE struct {
	Type_   util.Optional[string] `json:"_type,omitzero"`
	Meaning X_DV_TEXT             `json:"meaning"`
	Range   DV_INTERVAL           `json:"range"`
}

func (r REFERENCE_RANGE) HasModelName() bool {
	return r.Type_.E
}

func (r *REFERENCE_RANGE) SetModelName() {
	r.Type_ = util.Some(REFERENCE_RANGE_MODEL_NAME)
	r.Meaning.SetModelName()
	r.Range.SetModelName()
}

func (r REFERENCE_RANGE) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError

	// Validate _type
	if r.Type_.E && r.Type_.V != REFERENCE_RANGE_MODEL_NAME {
		errors = append(errors, util.ValidationError{
			Model:          REFERENCE_RANGE_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", REFERENCE_RANGE_MODEL_NAME, r.Type_.V),
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
