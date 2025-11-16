package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const GENERIC_ID_MODEL_NAME string = "GENERIC_ID"

type GENERIC_ID struct {
	Type_  util.Optional[string] `json:"_type,omitzero"`
	Value  string                `json:"value"`
	Scheme string                `json:"scheme"`
}

func (g *GENERIC_ID) isObjectIDModel() {}

func (g *GENERIC_ID) HasModelName() bool {
	return g.Type_.E
}

func (g *GENERIC_ID) SetModelName() {
	g.Type_ = util.Some(GENERIC_ID_MODEL_NAME)
}

func (g *GENERIC_ID) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if g.Type_.E && g.Type_.V != GENERIC_ID_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          GENERIC_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", GENERIC_ID_MODEL_NAME, g.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", GENERIC_ID_MODEL_NAME),
		})
	}

	// Validate value
	if g.Value == "" {
		attrPath = path + ".value"
		errors = append(errors, util.ValidationError{
			Model:          GENERIC_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	}

	// Validate scheme
	if g.Scheme == "" {
		attrPath = path + ".scheme"
		errors = append(errors, util.ValidationError{
			Model:          GENERIC_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "scheme field cannot be empty",
			Recommendation: "Ensure scheme field is not empty",
		})
	}

	return errors
}
