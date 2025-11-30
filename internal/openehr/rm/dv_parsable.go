package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_PARSABLE_TYPE string = "DV_PARSABLE"

type DV_PARSABLE struct {
	Type_     utils.Optional[string]      `json:"_type,omitzero"`
	Charset   utils.Optional[CODE_PHRASE] `json:"charset,omitzero"`
	Language  utils.Optional[CODE_PHRASE] `json:"language,omitzero"`
	Value     string                      `json:"value"`
	Formalism string                      `json:"formalism"`
}

func (d *DV_PARSABLE) SetModelName() {
	d.Type_ = utils.Some(DV_PARSABLE_TYPE)
	if d.Charset.E {
		d.Charset.V.SetModelName()
	}
	if d.Language.E {
		d.Language.V.SetModelName()
	}
}

func (d *DV_PARSABLE) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_PARSABLE_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_PARSABLE_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_PARSABLE",
		})
	}

	// Validate charset
	if d.Charset.E {
		attrPath = path + ".charset"

		if !terminology.IsValidCharsetTerminologyID(d.Charset.V.TerminologyID.Value) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_PARSABLE_TYPE,
				Path:           attrPath,
				Message:        "invalid charset terminology ID",
				Recommendation: "Ensure charset terminology ID is valid",
			})
		}

		if !terminology.IsValidCharset(d.Charset.V.CodeString) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_PARSABLE_TYPE,
				Path:           attrPath,
				Message:        "invalid charset code string",
				Recommendation: "Ensure charset code string is valid",
			})
		}

		validateErr.Errs = append(validateErr.Errs, d.Charset.V.Validate(attrPath).Errs...)
	}

	// Validate language
	if d.Language.E {
		attrPath = path + ".language"
		if !terminology.IsValidLanguageTerminologyID(d.Language.V.TerminologyID.Value) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_PARSABLE_TYPE,
				Path:           attrPath,
				Message:        "invalid language terminology ID",
				Recommendation: "Ensure language terminology ID is valid",
			})
		}
		if !terminology.IsValidLanguageCode(d.Language.V.CodeString) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_PARSABLE_TYPE,
				Path:           attrPath,
				Message:        "invalid language code string",
				Recommendation: "Ensure language code string is valid",
			})
		}

		validateErr.Errs = append(validateErr.Errs, d.Language.V.Validate(attrPath).Errs...)
	}

	// Validate value
	if d.Value == "" {
		attrPath = path + ".value"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_PARSABLE_TYPE,
			Path:           attrPath,
			Message:        "value cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	}

	// Validate formalism
	if d.Formalism == "" {
		attrPath = path + ".formalism"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_PARSABLE_TYPE,
			Path:           attrPath,
			Message:        "formalism cannot be empty",
			Recommendation: "Ensure formalism field is not empty",
		})
	}

	return validateErr
}
