package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_URI_MODEL_NAME string = "DV_URI"

type DV_URI struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
	Value string                `json:"value"`
}

func (d *DV_URI) isDataValueModel() {}

func (d *DV_URI) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_URI) SetModelName() {
	d.Type_ = util.Some(DV_URI_MODEL_NAME)
}

func (d *DV_URI) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_URI_MODEL_NAME {
		errors = append(errors, util.ValidationError{
			Model:          DV_URI_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_URI_MODEL_NAME, d.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_URI_MODEL_NAME),
		})
	}

	// Validate value
	attrPath := path + ".value"
	if d.Value == "" {
		errors = append(errors, util.ValidationError{
			Model:          DV_URI_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else if !util.URIRegex.MatchString(d.Value) {
		errors = append(errors, util.ValidationError{
			Model:          DV_URI_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid URI value: %s", d.Value),
			Recommendation: "Ensure value field is a valid URI according to RFC 3986",
		})
	}

	return errors
}
