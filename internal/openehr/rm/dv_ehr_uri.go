package rm

import (
	"fmt"
	"strings"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_EHR_URI_TYPE string = "DV_EHR_URI"

type DV_EHR_URI struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Value string                 `json:"value"`
}

func (d *DV_EHR_URI) SetModelName() {
	d.Type_ = utils.Some(DV_EHR_URI_TYPE)
}

func (d *DV_EHR_URI) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_EHR_URI_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_EHR_URI_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_EHR_URI_TYPE, d.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_EHR_URI_TYPE),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if d.Value == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_EHR_URI_TYPE,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else {
		if !strings.HasPrefix(d.Value, "ehr://") {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_EHR_URI_TYPE,
				Path:           attrPath,
				Message:        "value field must start with 'ehr://'",
				Recommendation: "Ensure value field starts with 'ehr://'",
			})
		}
	}

	return validateErr
}
