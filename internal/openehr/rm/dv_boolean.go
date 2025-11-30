package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_BOOLEAN_TYPE string = "DV_BOOLEAN"

type DV_BOOLEAN struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Value bool                   `json:"value"`
}

func (d *DV_BOOLEAN) SetModelName() {
	d.Type_ = utils.Some(DV_BOOLEAN_TYPE)
}

func (d *DV_BOOLEAN) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_BOOLEAN_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_BOOLEAN_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_BOOLEAN",
		})
	}

	return validateErr
}
