package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const CODE_PHRASE_MODEL_NAME string = "CODE_PHRASE"

type CODE_PHRASE struct {
	Type_         util.Optional[string] `json:"_type,omitzero"`
	TerminologyID TERMINOLOGY_ID        `json:"terminology_id"`
	CodeString    string                `json:"code_string"`
	PreferredTerm util.Optional[string] `json:"preferred_term,omitzero"`
}

func (c *CODE_PHRASE) SetModelName() {
	c.Type_ = util.Some(CODE_PHRASE_MODEL_NAME)
	c.TerminologyID.SetModelName()
}

func (c *CODE_PHRASE) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if c.Type_.E && c.Type_.V != CODE_PHRASE_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          CODE_PHRASE_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", CODE_PHRASE_MODEL_NAME, c.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", CODE_PHRASE_MODEL_NAME),
		})
	}

	// Validate terminology_id
	attrPath = path + ".terminology_id"
	validateErr.Errs = append(validateErr.Errs, c.TerminologyID.Validate(attrPath).Errs...)

	// Validate code_string
	if c.CodeString == "" {
		attrPath = path + ".code_string"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          CODE_PHRASE_MODEL_NAME,
			Path:           attrPath,
			Message:        "code_string field is required",
			Recommendation: "Ensure code_string field is not empty",
		})
	}

	return validateErr
}
