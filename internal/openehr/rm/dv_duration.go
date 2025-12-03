package rm

import (
	"fmt"
	"slices"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_DURATION_TYPE string = "DV_DURATION"

type DV_DURATION struct {
	Type_                utils.Optional[string]                          `json:"_type,omitzero"`
	NormalStatus         utils.Optional[CODE_PHRASE]                     `json:"normal_status,omitzero"`
	NormalRange          utils.Optional[DV_INTERVAL[*DV_DURATION]]       `json:"normal_range,omitzero"`
	OtherReferenceRanges utils.Optional[[]REFERENCE_RANGE[*DV_DURATION]] `json:"other_reference_ranges,omitzero"`
	MagnitudeStatus      utils.Optional[string]                          `json:"magnitude_status,omitzero"`
	AccuracyIsPercent    utils.Optional[bool]                            `json:"accuracy_is_percent,omitzero"`
	Accuracy             utils.Optional[float64]                         `json:"accuracy,omitzero"`
	Value                string                                          `json:"value"`
}

func (d *DV_DURATION) SetModelName() {
	d.Type_ = utils.Some(DV_DURATION_TYPE)
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

func (d *DV_DURATION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_DURATION_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_DURATION_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_DURATION_TYPE, d.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_DURATION_TYPE),
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
			itemPath := fmt.Sprintf("%s.other_reference_ranges[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, d.OtherReferenceRanges.V[i].Validate(itemPath).Errs...)
		}
	}

	// Validate magnitude_status
	if d.MagnitudeStatus.E {
		attrPath = path + ".magnitude_status"
		validValues := []string{"<", ">", "<=", ">=", "=", "~"}
		isValid := slices.Contains(validValues, d.MagnitudeStatus.V)
		if !isValid {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_DURATION_TYPE,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s magnitude_status field: %s", DV_DURATION_TYPE, d.MagnitudeStatus.V),
				Recommendation: "Ensure magnitude_status field is one of '<', '>', '<=', '>=', '=', '~'",
			})
		}
	}

	// Validate accuracy
	if d.Accuracy.E {
		attrPath = path + ".accuracy"
		value := d.Accuracy.V
		if value < 0 {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_DURATION_TYPE,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s accuracy field: %f", DV_DURATION_TYPE, value),
				Recommendation: "Ensure accuracy field is a non-negative number",
			})
		}
	}

	return validateErr
}
