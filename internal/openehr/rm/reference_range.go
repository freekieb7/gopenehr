package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const REFERENCE_RANGE_TYPE string = "REFERENCE_RANGE"

type REFERENCE_RANGE struct {
	Type_   utils.Optional[string] `json:"_type,omitzero"`
	Meaning DvTextUnion            `json:"meaning"`
	Range   DV_INTERVAL            `json:"range"`
}

func (r *REFERENCE_RANGE) SetModelName() {
	r.Type_ = utils.Some(REFERENCE_RANGE_TYPE)
	r.Meaning.SetModelName()
	r.Range.SetModelName()
}

func (r *REFERENCE_RANGE) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError

	// Validate _type
	if r.Type_.E && r.Type_.V != REFERENCE_RANGE_TYPE {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          REFERENCE_RANGE_TYPE,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", REFERENCE_RANGE_TYPE, r.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", REFERENCE_RANGE_TYPE),
		})
	}

	// Validate meaning
	attrPath := path + ".meaning"
	validateErr.Errs = append(validateErr.Errs, r.Meaning.Validate(attrPath).Errs...)

	// Validate range
	attrPath = path + ".range"
	validateErr.Errs = append(validateErr.Errs, r.Range.Validate(attrPath).Errs...)

	return validateErr
}
