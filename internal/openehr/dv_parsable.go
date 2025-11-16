package openehr

import (
	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_PARSABLE_MODEL_NAME string = "DV_PARSABLE"

type DV_PARSABLE struct {
	Type_     util.Optional[string]      `json:"_type,omitzero"`
	Charset   util.Optional[CODE_PHRASE] `json:"charset,omitzero"`
	Language  util.Optional[CODE_PHRASE] `json:"language,omitzero"`
	Value     string                     `json:"value"`
	Formalism string                     `json:"formalism"`
}

func (d *DV_PARSABLE) isDataValueModel() {}

func (d *DV_PARSABLE) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_PARSABLE) SetModelName() {
	d.Type_ = util.Some(DV_PARSABLE_MODEL_NAME)
	if d.Charset.E {
		d.Charset.V.SetModelName()
	}
	if d.Language.E {
		d.Language.V.SetModelName()
	}
}

func (d *DV_PARSABLE) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_PARSABLE_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          DV_PARSABLE_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_PARSABLE",
		})
	}

	// Validate charset
	if d.Charset.E {
		attrPath = path + ".charset"

		if !terminology.IsValidCharsetTerminologyID(d.Charset.V.TerminologyID.Value) {
			errors = append(errors, util.ValidationError{
				Model:          DV_PARSABLE_MODEL_NAME,
				Path:           attrPath,
				Message:        "invalid charset terminology ID",
				Recommendation: "Ensure charset terminology ID is valid",
			})
		}

		if !terminology.IsValidCharset(d.Charset.V.CodeString) {
			errors = append(errors, util.ValidationError{
				Model:          DV_PARSABLE_MODEL_NAME,
				Path:           attrPath,
				Message:        "invalid charset code string",
				Recommendation: "Ensure charset code string is valid",
			})
		}

		errors = append(errors, d.Charset.V.Validate(attrPath)...)
	}

	// Validate language
	if d.Language.E {
		attrPath = path + ".language"
		if !terminology.IsValidLanguageTerminologyID(d.Language.V.TerminologyID.Value) {
			errors = append(errors, util.ValidationError{
				Model:          DV_PARSABLE_MODEL_NAME,
				Path:           attrPath,
				Message:        "invalid language terminology ID",
				Recommendation: "Ensure language terminology ID is valid",
			})
		}
		if !terminology.IsValidLanguageCode(d.Language.V.CodeString) {
			errors = append(errors, util.ValidationError{
				Model:          DV_PARSABLE_MODEL_NAME,
				Path:           attrPath,
				Message:        "invalid language code string",
				Recommendation: "Ensure language code string is valid",
			})
		}

		errors = append(errors, d.Language.V.Validate(attrPath)...)
	}

	// Validate value
	if d.Value == "" {
		attrPath = path + ".value"
		errors = append(errors, util.ValidationError{
			Model:          DV_PARSABLE_MODEL_NAME,
			Path:           attrPath,
			Message:        "value cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	}

	// Validate formalism
	if d.Formalism == "" {
		attrPath = path + ".formalism"
		errors = append(errors, util.ValidationError{
			Model:          DV_PARSABLE_MODEL_NAME,
			Path:           attrPath,
			Message:        "formalism cannot be empty",
			Recommendation: "Ensure formalism field is not empty",
		})
	}

	return errors
}
