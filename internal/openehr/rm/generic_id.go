package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const GENERIC_ID_TYPE string = "GENERIC_ID"

type GENERIC_ID struct {
	Type_  utils.Optional[string] `json:"_type,omitzero"`
	Value  string                 `json:"value"`
	Scheme string                 `json:"scheme"`
}

func (g *GENERIC_ID) SetModelName() {
	g.Type_ = utils.Some(GENERIC_ID_TYPE)
}

func (g *GENERIC_ID) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if g.Type_.E && g.Type_.V != GENERIC_ID_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          GENERIC_ID_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", GENERIC_ID_TYPE, g.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", GENERIC_ID_TYPE),
		})
	}

	// Validate value
	if g.Value == "" {
		attrPath = path + ".value"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          GENERIC_ID_TYPE,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	}

	// Validate scheme
	if g.Scheme == "" {
		attrPath = path + ".scheme"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          GENERIC_ID_TYPE,
			Path:           attrPath,
			Message:        "scheme field cannot be empty",
			Recommendation: "Ensure scheme field is not empty",
		})
	}

	return validateErr
}
