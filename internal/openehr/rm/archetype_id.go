package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ARCHETYPE_ID_TYPE string = "ARCHETYPE_ID"

type ARCHETYPE_ID struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Value string                 `json:"value"`
}

func (a *ARCHETYPE_ID) SetModelName() {
	a.Type_ = utils.Some(ARCHETYPE_ID_TYPE)
}

func (a *ARCHETYPE_ID) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ARCHETYPE_ID_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ARCHETYPE_ID_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", ARCHETYPE_ID_TYPE, a.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", ARCHETYPE_ID_TYPE),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if a.Value == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ARCHETYPE_ID_TYPE,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else if !util.ArchetypeIDRegex.MatchString(a.Value) {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ARCHETYPE_ID_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid value format: %s", a.Value),
			Recommendation: "Ensure value field follows the lexical form: rm_originator '-' rm_name '-' rm_entity '.' concept_name { '-' specialisation }* '.v' number.",
		})
	}

	return validateErr
}
