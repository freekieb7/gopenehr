package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ARCHETYPE_ID_MODEL_NAME string = "ARCHETYPE_ID"

type ARCHETYPE_ID struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
	Value string                `json:"value"`
}

func (a ARCHETYPE_ID) isObjectIDModel() {}

func (a ARCHETYPE_ID) HasModelName() bool {
	return a.Type_.E
}

func (a *ARCHETYPE_ID) SetModelName() {
	a.Type_ = util.Some(ARCHETYPE_ID_MODEL_NAME)
}

func (a ARCHETYPE_ID) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ARCHETYPE_ID_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          ARCHETYPE_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", ARCHETYPE_ID_MODEL_NAME, a.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", ARCHETYPE_ID_MODEL_NAME),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if a.Value == "" {
		errors = append(errors, util.ValidationError{
			Model:          ARCHETYPE_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else if !util.ArchetypeIDRegex.MatchString(a.Value) {
		errors = append(errors, util.ValidationError{
			Model:          ARCHETYPE_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid value format: %s", a.Value),
			Recommendation: "Ensure value field follows the lexical form: rm_originator '-' rm_name '-' rm_entity '.' concept_name { '-' specialisation }* '.v' number.",
		})
	}

	return errors
}
