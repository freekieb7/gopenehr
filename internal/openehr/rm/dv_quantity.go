package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_QUANTITY_TYPE string = "DV_QUANTITY"

type DV_QUANTITY struct {
	Type_                utils.Optional[string]                          `json:"_type,omitzero"`
	NormalStatus         utils.Optional[CODE_PHRASE]                     `json:"normal_status,omitzero"`
	NormalRange          utils.Optional[DV_INTERVAL[*DV_QUANTITY]]       `json:"normal_range,omitzero"`
	OtherReferenceRanges utils.Optional[[]REFERENCE_RANGE[*DV_QUANTITY]] `json:"other_reference_ranges,omitzero"`
	MagnitudeStatus      utils.Optional[string]                          `json:"magnitude_status,omitzero"`
	AccuracyIsPercent    utils.Optional[bool]                            `json:"accuracy_is_percent,omitzero"`
	Accuracy             utils.Optional[float64]                         `json:"accuracy,omitzero"`
	Magnitude            float64                                         `json:"magnitude"`
	Precision            utils.Optional[int]                             `json:"precision,omitzero"`
	Units                string                                          `json:"units"`
	UnitsSystem          utils.Optional[string]                          `json:"units_system,omitzero"`
	UnitsDisplayName     utils.Optional[string]                          `json:"units_display_name,omitzero"`
}

func (d *DV_QUANTITY) SetModelName() {
	d.Type_ = utils.Some(DV_QUANTITY_TYPE)
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

func (d *DV_QUANTITY) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_QUANTITY_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_QUANTITY_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_QUANTITY",
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

	// // Validate magnitude_status
	// if d.MagnitudeStatus.E {
	// 	attrPath = path + ".magnitude_status"
	// 	// Add any specific validation for magnitude_status if needed
	// }

	// // Validate accuracy_is_percent
	// if d.AccuracyIsPercent.E {
	// 	attrPath = path + ".accuracy_is_percent"
	// 	// Add any specific validation for accuracy_is_percent if needed
	// }

	// // Validate accuracy
	// if d.Accuracy.E {
	// 	attrPath = path + ".accuracy"
	// 	// Add any specific validation for accuracy if needed
	// }

	// // Validate precision
	// if d.Precision.E {
	// 	attrPath = path + ".precision"
	// 	// Add any specific validation for precision if needed
	// }

	// // Validate units_system
	// if d.UnitsSystem.E {
	// 	attrPath = path + ".units_system"
	// 	// Add any specific validation for units_system if needed
	// }

	// // Validate units_display_name
	// if d.UnitsDisplayName.E {
	// 	attrPath = path + ".units_display_name"
	// 	// Add any specific validation for units_display_name if needed
	// }

	return validateErr
}
