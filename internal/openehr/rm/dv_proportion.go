package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_PROPORTION_TYPE string = "DV_PROPORTION"

type DV_PROPORTION struct {
	Type_                utils.Optional[string]            `json:"_type,omitzero"`
	NormalStatus         utils.Optional[CODE_PHRASE]       `json:"normal_status,omitzero"`
	MagnitudeStatus      utils.Optional[string]            `json:"magnitude_status,omitzero"`
	AccuracyIsPercent    utils.Optional[bool]              `json:"accuracy_is_percent,omitzero"`
	Accuracy             utils.Optional[float64]           `json:"accuracy,omitzero"`
	Numerator            float64                           `json:"numerator"`
	Denominator          float64                           `json:"denominator"`
	Type                 int64                             `json:"type"`
	Precision            utils.Optional[int64]             `json:"precision,omitzero"`
	NormalRange          utils.Optional[DV_INTERVAL]       `json:"normal_range,omitzero"`
	OtherReferenceRanges utils.Optional[[]REFERENCE_RANGE] `json:"other_reference_ranges,omitzero"`
}

func (d *DV_PROPORTION) SetModelName() {
	d.Type_ = utils.Some(DV_PROPORTION_TYPE)
	if d.NormalStatus.E {
		d.NormalStatus.V.SetModelName()
	}
	if d.NormalRange.E {
		d.NormalRange.V.SetModelName()
	}
	if d.OtherReferenceRanges.E {
		for i := range d.OtherReferenceRanges.V {
			d.OtherReferenceRanges.V[i].SetModelName()
		}
	}
}

func (d *DV_PROPORTION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_PROPORTION_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_PROPORTION_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_PROPORTION",
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

	return validateErr
}
