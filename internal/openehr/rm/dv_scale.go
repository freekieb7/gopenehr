package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_SCALE_TYPE string = "DV_SCALE"

type DV_SCALE struct {
	Type_                utils.Optional[string]            `json:"_type,omitzero"`
	NormalStatus         utils.Optional[CODE_PHRASE]       `json:"normal_status,omitzero"`
	NormalRange          utils.Optional[DV_INTERVAL]       `json:"normal_range,omitzero"`
	OtherReferenceRanges utils.Optional[[]REFERENCE_RANGE] `json:"other_reference_ranges,omitzero"`
	Symbol               DV_CODED_TEXT                     `json:"symbol"`
	Value                float64                           `json:"value"`
}

func (d *DV_SCALE) SetModelName() {
	d.Type_ = utils.Some(DV_SCALE_TYPE)
}

func (d *DV_SCALE) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_SCALE_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_SCALE_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_SCALE",
		})
	}

	// Validate normal_status
	if d.NormalStatus.E {
		attrPath = path + ".normal_status"
		validateErr.Errs = append(validateErr.Errs, d.NormalStatus.V.Validate(attrPath).Errs...)
	}

	// Validate normal_range
	if d.NormalRange.E {
		attrPath = path + ".normal_range"
		validateErr.Errs = append(validateErr.Errs, d.NormalRange.V.Validate(attrPath).Errs...)
	}

	// Validate other_reference_ranges
	if d.OtherReferenceRanges.E {
		for i := range d.OtherReferenceRanges.V {
			itemPath := fmt.Sprintf("%s.other_reference_ranges[%d]", attrPath, i)
			validateErr.Errs = append(validateErr.Errs, d.OtherReferenceRanges.V[i].Validate(itemPath).Errs...)
		}
	}

	// Validate symbol
	attrPath = path + ".symbol"
	validateErr.Errs = append(validateErr.Errs, d.Symbol.Validate(attrPath).Errs...)

	return validateErr
}
