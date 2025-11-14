package openehr

import (
	"fmt"
	"slices"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_DURATION_MODEL_NAME string = "DV_DURATION"

var _ util.ReferenceModel = (*DV_DURATION)(nil)

type DV_DURATION struct {
	Type_                util.Optional[string]            `json:"_type,omitzero"`
	NormalStatus         util.Optional[CODE_PHRASE]       `json:"normal_status,omitzero"`
	NormalRange          util.Optional[DV_INTERVAL]       `json:"normal_range,omitzero"`
	OtherReferenceRanges util.Optional[[]REFERENCE_RANGE] `json:"other_reference_ranges,omitzero"`
	MagnitudeStatus      util.Optional[string]            `json:"magnitude_status,omitzero"`
	AccuracyIsPercent    util.Optional[bool]              `json:"accuracy_is_percent,omitzero"`
	Accuracy             util.Optional[float64]           `json:"accuracy,omitzero"`
	Value                string                           `json:"value"`
}

func (d DV_DURATION) HasModelName() bool {
	return d.Type_.IsSet()
}

func (d DV_DURATION) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.IsSet() && d.Type_.Unwrap() != DV_DURATION_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          DV_DURATION_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_DURATION_MODEL_NAME, d.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_DURATION_MODEL_NAME),
		})
	}

	// Validate normal_status
	if d.NormalStatus.IsSet() {
		attrPath = path + ".normal_status"
		errors = append(errors, d.NormalStatus.Unwrap().Validate(attrPath)...)
	}

	// Validate normal_range
	if d.NormalRange.IsSet() {
		attrPath = path + ".normal_range"
		errors = append(errors, d.NormalRange.Unwrap().Validate(attrPath)...)
	}

	// Validate other_reference_ranges
	if d.OtherReferenceRanges.IsSet() {
		attrPath = path + ".other_reference_ranges"
		for i, v := range d.OtherReferenceRanges.Unwrap() {
			errors = append(errors, v.Validate(fmt.Sprintf("%s[%d]", attrPath, i))...)
		}
	}

	// Validate magnitude_status
	if d.MagnitudeStatus.IsSet() {
		attrPath = path + ".magnitude_status"
		validValues := []string{"<", ">", "<=", ">=", "=", "~"}
		value := d.MagnitudeStatus.Unwrap()
		isValid := slices.Contains(validValues, d.MagnitudeStatus.Unwrap())
		if !isValid {
			errors = append(errors, util.ValidationError{
				Model:          DV_DURATION_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s magnitude_status field: %s", DV_DURATION_MODEL_NAME, value),
				Recommendation: "Ensure magnitude_status field is one of '<', '>', '<=', '>=', '=', '~'",
			})
		}
	}

	// Validate accuracy
	if d.Accuracy.IsSet() {
		attrPath = path + ".accuracy"
		value := d.Accuracy.Unwrap()
		if value < 0 {
			errors = append(errors, util.ValidationError{
				Model:          DV_DURATION_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s accuracy field: %f", DV_DURATION_MODEL_NAME, value),
				Recommendation: "Ensure accuracy field is a non-negative number",
			})
		}
	}

	return errors
}
