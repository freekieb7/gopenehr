package rm

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_DATE_TIME_TYPE string = "DV_DATE_TIME"

type DV_DATE_TIME struct {
	Type_                utils.Optional[string]            `json:"_type,omitzero"`
	NormalStatus         utils.Optional[CODE_PHRASE]       `json:"normal_status,omitzero"`
	NormalRange          utils.Optional[DV_INTERVAL]       `json:"normal_range,omitzero"`
	OtherReferenceRanges utils.Optional[[]REFERENCE_RANGE] `json:"other_reference_ranges,omitzero"`
	MagnitudeStatus      utils.Optional[string]            `json:"magnitude_status,omitzero"`
	Accuracy             utils.Optional[DV_DURATION]       `json:"accuracy,omitzero"`
	Value                string                            `json:"value"`
}

func (d *DV_DATE_TIME) SetModelName() {
	d.Type_ = utils.Some(DV_DATE_TIME_TYPE)
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
	if d.Accuracy.E {
		d.Accuracy.V.SetModelName()
	}
}

func (d *DV_DATE_TIME) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_DATE_TIME_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_DATE_TIME_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_DATE_TIME_TYPE, d.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_DATE_TIME_TYPE),
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
		if !slices.Contains([]string{"<", ">", "<=", ">=", "=", "~"}, d.MagnitudeStatus.V) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_DATE_TIME_TYPE,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s magnitude_status field: %s", DV_DATE_TIME_TYPE, d.MagnitudeStatus.V),
				Recommendation: "Ensure magnitude_status field is one of '<', '>', '<=', '>=', '=', '~'",
			})
		}
	}

	// Validate accuracy
	if d.Accuracy.E {
		attrPath = path + ".accuracy"
		validateErr.Errs = append(validateErr.Errs, d.Accuracy.V.Validate(attrPath).Errs...)
	}

	// Validate value
	attrPath = path + ".value"
	if d.Value == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:   DV_DATE_TIME_TYPE,
			Path:    attrPath,
			Message: fmt.Sprintf("%s value field is required", DV_DATE_TIME_TYPE),
		})
	} else if !strings.HasSuffix(d.Value, "Z") {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_DATE_TIME_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s value field: %s", DV_DATE_TIME_TYPE, d.Value),
			Recommendation: "Ensure value field is of format YYYY-MM-DDTHH:MM:SSZ",
		})
	} else {
		if _, err := time.Parse(time.RFC3339Nano, d.Value); err != nil {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_DATE_TIME_TYPE,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s value field: %s", DV_DATE_TIME_TYPE, d.Value),
				Recommendation: "Ensure value field is of format YYYY-MM-DDTHH:MM:SSZ",
			})
		}
	}

	return validateErr
}
