package rm

import (
	"fmt"
	"strings"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_EHR_URI_MODEL_NAME string = "DV_EHR_URI"

type DV_EHR_URI struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Value string                 `json:"value"`
}

func (d *DV_EHR_URI) isDataValueModel() {}

func (d *DV_EHR_URI) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_EHR_URI) SetModelName() {
	d.Type_ = utils.Some(DV_EHR_URI_MODEL_NAME)
}

func (d *DV_EHR_URI) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_EHR_URI_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_EHR_URI_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_EHR_URI_MODEL_NAME, d.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_EHR_URI_MODEL_NAME),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if d.Value == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_EHR_URI_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else {
		if !strings.HasPrefix(d.Value, "ehr://") {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_EHR_URI_MODEL_NAME,
				Path:           attrPath,
				Message:        "value field must start with 'ehr://'",
				Recommendation: "Ensure value field starts with 'ehr://'",
			})
		}
	}

	return validateErr
}
