package rm

import (
	"fmt"
	"slices"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const TERM_MAPPING_TYPE string = "TERM_MAPPING"

type TERM_MAPPING struct {
	Type_   utils.Optional[string]        `json:"_type,omitzero"`
	Match   byte                          `json:"match"`
	Purpose utils.Optional[DV_CODED_TEXT] `json:"purpose,omitzero"`
	Target  CODE_PHRASE                   `json:"target"`
}

func (t *TERM_MAPPING) SetModelName() {
	t.Type_ = utils.Some(TERM_MAPPING_TYPE)
	if t.Purpose.E {
		t.Purpose.V.SetModelName()
	}
	t.Target.SetModelName()
}

func (t *TERM_MAPPING) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if t.Type_.E && t.Type_.V != TERM_MAPPING_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          TERM_MAPPING_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", TERM_MAPPING_TYPE, t.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", TERM_MAPPING_TYPE),
		})
	}

	// Validate purpose
	if t.Purpose.E {
		attrPath = path + ".purpose"
		validateErr.Errs = append(validateErr.Errs, t.Purpose.V.Validate(attrPath).Errs...)
	}

	// Validate target
	attrPath = path + ".target"
	validateErr.Errs = append(validateErr.Errs, t.Target.Validate(attrPath).Errs...)

	// Validate match
	validMatches := []byte{'=', '>', '<', '?'}
	if !slices.Contains(validMatches, t.Match) {
		attrPath = path + ".match"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          TERM_MAPPING_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid match value: %c", t.Match),
			Recommendation: "Ensure match field is one of '=', '>', '<', '?'",
		})
	}

	return validateErr
}
