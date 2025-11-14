package openehr

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_DATE_TIME_MODEL_NAME string = "DV_DATE_TIME"

var _ util.ReferenceModel = (*DV_DATE_TIME)(nil)

type DV_DATE_TIME struct {
	Type_                util.Optional[string]            `json:"_type"`
	NormalStatus         util.Optional[CODE_PHRASE]       `json:"normal_status"`
	NormalRange          util.Optional[DV_INTERVAL]       `json:"normal_range"`
	OtherReferenceRanges util.Optional[[]REFERENCE_RANGE] `json:"other_reference_ranges"`
	MagnitudeStatus      util.Optional[string]            `json:"magnitude_status"`
	Accuracy             util.Optional[DV_DURATION]       `json:"accuracy"`
	Value                string                           `json:"value"`
}

// HasModelName implements util.ReferenceModel.
func (d DV_DATE_TIME) HasModelName() bool {
	return d.Type_.IsSet()
}

// Validate implements util.ReferenceModel.
func (d DV_DATE_TIME) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.IsSet() && d.Type_.Unwrap() != DV_DATE_TIME_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          DV_DATE_TIME_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_DATE_TIME_MODEL_NAME, d.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_DATE_TIME_MODEL_NAME),
		})
	}

	// Validate normal_status
	if d.NormalStatus.IsSet() {
		attrPath = path + ".normal_status"
		errs = append(errs, d.NormalStatus.Unwrap().Validate(attrPath)...)
	}

	// Validate normal_range
	if d.NormalRange.IsSet() {
		attrPath = path + ".normal_range"
		errs = append(errs, d.NormalRange.Unwrap().Validate(attrPath)...)
	}

	// Validate other_reference_ranges
	if d.OtherReferenceRanges.IsSet() {
		attrPath = path + ".other_reference_ranges"
		for i, v := range d.OtherReferenceRanges.Unwrap() {
			errs = append(errs, v.Validate(fmt.Sprintf("%s[%d]", attrPath, i))...)
		}
	}

	// Validate magnitude_status
	if d.MagnitudeStatus.IsSet() {
		attrPath = path + ".magnitude_status"
		if !slices.Contains([]string{"<", ">", "<=", ">=", "=", "~"}, d.MagnitudeStatus.Unwrap()) {
			errs = append(errs, util.ValidationError{
				Model:          DV_DATE_TIME_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s magnitude_status field: %s", DV_DATE_TIME_MODEL_NAME, d.MagnitudeStatus.Unwrap()),
				Recommendation: "Ensure magnitude_status field is one of '<', '>', '<=', '>=', '=', '~'",
			})
		}
	}

	// Validate accuracy
	if d.Accuracy.IsSet() {
		attrPath = path + ".accuracy"
		errs = append(errs, d.Accuracy.Unwrap().Validate(attrPath)...)
	}

	// Validate value
	attrPath = path + ".value"
	if d.Value == "" {
		errs = append(errs, util.ValidationError{
			Model:   DV_DATE_TIME_MODEL_NAME,
			Path:    attrPath,
			Message: fmt.Sprintf("%s value field is required", DV_DATE_TIME_MODEL_NAME),
		})
	} else if !strings.HasSuffix(d.Value, "Z") {
		errs = append(errs, util.ValidationError{
			Model:          DV_DATE_TIME_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s value field: %s", DV_DATE_TIME_MODEL_NAME, d.Value),
			Recommendation: "Ensure value field is of format YYYY-MM-DDTHH:MM:SSZ",
		})
	} else {
		if _, err := time.Parse(time.RFC3339Nano, d.Value); err != nil {
			errs = append(errs, util.ValidationError{
				Model:          DV_DATE_TIME_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s value field: %s", DV_DATE_TIME_MODEL_NAME, d.Value),
				Recommendation: "Ensure value field is of format YYYY-MM-DDTHH:MM:SSZ",
			})
		}
	}

	return errs
}
