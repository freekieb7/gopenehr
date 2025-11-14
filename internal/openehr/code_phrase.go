package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const CODE_PHRASE_MODEL_NAME string = "CODE_PHRASE"

var _ util.ReferenceModel = (*CODE_PHRASE)(nil)

type CODE_PHRASE struct {
	Type_         util.Optional[string] `json:"_type,omitzero"`
	TerminologyId TERMINOLOGY_ID        `json:"terminology_id"`
	CodeString    string                `json:"code_string"`
	PreferredTerm util.Optional[string] `json:"preferred_term,omitzero"`
}

func (c CODE_PHRASE) HasModelName() bool {
	return c.Type_.IsSet()
}

func (c CODE_PHRASE) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError

	// Validate _type
	if c.Type_.IsSet() && c.Type_.Unwrap() != CODE_PHRASE_MODEL_NAME {
		errors = append(errors, util.ValidationError{
			Model:          CODE_PHRASE_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", CODE_PHRASE_MODEL_NAME, c.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", CODE_PHRASE_MODEL_NAME),
		})
	}

	// Validate terminology_id
	attrPath := path + ".terminology_id"
	errors = append(errors, c.TerminologyId.Validate(attrPath)...)

	// Validate code_string
	attrPath = path + ".code_string"
	if c.CodeString == "" {
		errors = append(errors, util.ValidationError{
			Model:          CODE_PHRASE_MODEL_NAME,
			Path:           attrPath,
			Message:        "code_string field is required",
			Recommendation: "Ensure code_string field is not empty",
		})
	}

	return errors
}
