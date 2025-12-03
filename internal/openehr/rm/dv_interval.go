package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_INTERVAL_TYPE string = "DV_INTERVAL"

type DV_INTERVAL[T any] struct {
	Type_          utils.Optional[string] `json:"_type,omitzero"`
	Lower          T                      `json:"lower"`
	Upper          T                      `json:"upper"`
	LowerUnbounded bool                   `json:"lower_unbounded"`
	UpperUnbounded bool                   `json:"upper_unbounded"`
	LowerIncluded  bool                   `json:"lower_included"`
	UpperIncluded  bool                   `json:"upper_included"`
}

func (d *DV_INTERVAL[T]) SetModelName() {
	d.Type_ = utils.Some(DV_INTERVAL_TYPE)
}

func (d *DV_INTERVAL[T]) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_INTERVAL_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_INTERVAL_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_INTERVAL_TYPE, d.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_INTERVAL_TYPE),
		})
	}

	return validateErr
}
