package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_INTERVAL_MODEL_NAME string = "DV_INTERVAL"

type DV_INTERVAL struct {
	Type_          util.Optional[string] `json:"_type,omitzero"`
	Lower          any                   `json:"lower"`
	Upper          any                   `json:"upper"`
	LowerUnbounded bool                  `json:"lower_unbounded"`
	UpperUnbounded bool                  `json:"upper_unbounded"`
	LowerIncluded  bool                  `json:"lower_included"`
	UpperIncluded  bool                  `json:"upper_included"`
}

func (d *DV_INTERVAL) isDataValueModel() {}

func (d *DV_INTERVAL) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_INTERVAL) SetModelName() {
	d.Type_ = util.Some(DV_INTERVAL_MODEL_NAME)
}

func (d *DV_INTERVAL) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_INTERVAL_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_INTERVAL_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_INTERVAL_MODEL_NAME, d.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_INTERVAL_MODEL_NAME),
		})
	}

	return validateErr
}
