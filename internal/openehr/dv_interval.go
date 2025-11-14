package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_INTERVAL_MODEL_NAME string = "DV_INTERVAL"

var _ util.ReferenceModel = (*DV_INTERVAL)(nil)

type DV_INTERVAL struct {
	Type_          util.Optional[string] `json:"_type,omitzero"`
	Lower          any                   `json:"lower"`
	Upper          any                   `json:"upper"`
	LowerUnbounded bool                  `json:"lower_unbounded"`
	UpperUnbounded bool                  `json:"upper_unbounded"`
	LowerIncluded  bool                  `json:"lower_included"`
	UpperIncluded  bool                  `json:"upper_included"`
}

func (d DV_INTERVAL) HasModelName() bool {
	return d.Type_.IsSet()
}

func (d DV_INTERVAL) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.IsSet() && d.Type_.Unwrap() != DV_INTERVAL_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          DV_INTERVAL_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_INTERVAL_MODEL_NAME, d.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_INTERVAL_MODEL_NAME),
		})
	}

	return errors
}
