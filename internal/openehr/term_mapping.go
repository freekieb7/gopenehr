package openehr

import (
	"fmt"
	"slices"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const TERM_MAPPING_MODEL_NAME string = "TERM_MAPPING"

var _ util.ReferenceModel = (*TERM_MAPPING)(nil)

type TERM_MAPPING struct {
	Type_   util.Optional[string]        `json:"_type,omitzero"`
	Match   byte                         `json:"match"`
	Purpose util.Optional[DV_CODED_TEXT] `json:"purpose,omitzero"`
	Target  CODE_PHRASE                  `json:"target"`
}

func (t TERM_MAPPING) HasModelName() bool {
	return t.Type_.IsSet()
}

func (t TERM_MAPPING) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError

	// Validate _type
	if t.Type_.IsSet() && t.Type_.Unwrap() != TERM_MAPPING_MODEL_NAME {
		errors = append(errors, util.ValidationError{
			Model:          TERM_MAPPING_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", TERM_MAPPING_MODEL_NAME, t.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", TERM_MAPPING_MODEL_NAME),
		})
	}

	// Validate purpose
	if t.Purpose.IsSet() {
		attrPath := path + ".purpose"
		errors = append(errors, t.Purpose.Unwrap().Validate(attrPath)...)
	}

	// Validate target
	attrPath := path + ".target"
	errors = append(errors, t.Target.Validate(attrPath)...)

	// Validate match
	validMatches := []byte{'=', '>', '<', '?'}
	if !slices.Contains(validMatches, t.Match) {
		errors = append(errors, util.ValidationError{
			Model:          TERM_MAPPING_MODEL_NAME,
			Path:           path + ".match",
			Message:        fmt.Sprintf("invalid match value: %c", t.Match),
			Recommendation: "Ensure match field is one of '=', '>', '<', '?'",
		})
	}

	return errors
}
