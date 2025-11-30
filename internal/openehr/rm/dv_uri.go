package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_URI_TYPE string = "DV_URI"

type DV_URI struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Value string                 `json:"value"`
}

func (d *DV_URI) SetModelName() {
	d.Type_ = utils.Some(DV_URI_TYPE)
}

func (d *DV_URI) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_URI_TYPE {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_URI_TYPE,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_URI_TYPE, d.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_URI_TYPE),
		})
	}

	// Validate value
	attrPath := path + ".value"
	if d.Value == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_URI_TYPE,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else if !util.URIRegex.MatchString(d.Value) {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_URI_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid URI value: %s", d.Value),
			Recommendation: "Ensure value field is a valid URI according to RFC 3986",
		})
	}

	return validateErr
}
