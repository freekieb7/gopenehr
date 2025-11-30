package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_ORDINAL_TYPE string = "DV_ORDINAL"

type DV_ORDINAL struct {
	Type_                utils.Optional[string]          `json:"_type,omitzero"`
	NormalStatus         utils.Optional[CODE_PHRASE]     `json:"normal_status,omitzero"`
	NormalRange          utils.Optional[DV_INTERVAL]     `json:"normal_range,omitzero"`
	OtherReferenceRanges utils.Optional[REFERENCE_RANGE] `json:"other_reference_ranges,omitzero"`
	Symbol               DV_CODED_TEXT                   `json:"symbol"`
	Value                int64                           `json:"value"`
}

func (d *DV_ORDINAL) SetModelName() {
	d.Type_ = utils.Some(DV_ORDINAL_TYPE)
	if d.NormalStatus.E {
		d.NormalStatus.V.SetModelName()
	}
	if d.NormalRange.E {
		d.NormalRange.V.SetModelName()
	}
	if d.OtherReferenceRanges.E {
		d.OtherReferenceRanges.V.SetModelName()
	}
	d.Symbol.SetModelName()
}

func (d *DV_ORDINAL) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_ORDINAL_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_ORDINAL_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_ORDINAL",
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
		attrPath = path + ".other_reference_ranges"
		validateErr.Errs = append(validateErr.Errs, d.OtherReferenceRanges.V.Validate(attrPath).Errs...)
	}

	// Validate symbol
	attrPath = path + ".symbol"
	validateErr.Errs = append(validateErr.Errs, d.Symbol.Validate(attrPath).Errs...)

	return validateErr
}
