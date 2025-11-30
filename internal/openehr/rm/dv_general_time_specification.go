package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_GENERAL_TIME_SPECIFICATION_TYPE string = "DV_GENERAL_TIME_SPECIFICATION"

type DV_GENERAL_TIME_SPECIFICATION struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Value DV_PARSABLE            `json:"value"`
}

func (d *DV_GENERAL_TIME_SPECIFICATION) SetModelName() {
	d.Type_ = utils.Some(DV_GENERAL_TIME_SPECIFICATION_TYPE)
	d.Value.SetModelName()
}

func (d *DV_GENERAL_TIME_SPECIFICATION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_GENERAL_TIME_SPECIFICATION_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_GENERAL_TIME_SPECIFICATION_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_GENERAL_TIME_SPECIFICATION",
		})
	}

	// Validate value
	attrPath = path + ".value"
	validateErr.Errs = append(validateErr.Errs, d.Value.Validate(attrPath).Errs...)

	return validateErr
}
