package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_SCALE_MODEL_NAME string = "DV_SCALE"

type DV_SCALE struct {
	Type_                util.Optional[string]            `json:"_type,omitzero"`
	NormalStatus         util.Optional[CODE_PHRASE]       `json:"normal_status,omitzero"`
	NormalRange          util.Optional[DV_INTERVAL]       `json:"normal_range,omitzero"`
	OtherReferenceRanges util.Optional[[]REFERENCE_RANGE] `json:"other_reference_ranges,omitzero"`
	Symbol               DV_CODED_TEXT                    `json:"symbol"`
	Value                float64                          `json:"value"`
}

func (d *DV_SCALE) isDataValueModel() {}

func (d *DV_SCALE) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_SCALE) SetModelName() {
	d.Type_ = util.Some(DV_SCALE_MODEL_NAME)
}

func (d *DV_SCALE) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_SCALE_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          DV_SCALE_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_SCALE",
		})
	}

	// Validate normal_status
	if d.NormalStatus.E {
		attrPath = path + ".normal_status"
		errors = append(errors, d.NormalStatus.V.Validate(attrPath)...)
	}

	// Validate normal_range
	if d.NormalRange.E {
		attrPath = path + ".normal_range"
		errors = append(errors, d.NormalRange.V.Validate(attrPath)...)
	}

	// Validate other_reference_ranges
	if d.OtherReferenceRanges.E {
		for i := range d.OtherReferenceRanges.V {
			itemPath := fmt.Sprintf("%s.other_reference_ranges[%d]", attrPath, i)
			errors = append(errors, d.OtherReferenceRanges.V[i].Validate(itemPath)...)
		}
	}

	// Validate symbol
	attrPath = path + ".symbol"
	errors = append(errors, d.Symbol.Validate(attrPath)...)

	return errors
}
