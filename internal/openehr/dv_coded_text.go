package openehr

import (
	"fmt"
	"slices"
	"strings"

	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_CODED_TEXT_MODEL_NAME string = "DV_CODED_TEXT"

type DV_CODED_TEXT struct {
	Type_        util.Optional[string]         `json:"_type,omitzero"`
	Value        string                        `json:"value"`
	Hyperlink    util.Optional[DV_URI]         `json:"hyperlink,omitzero"`
	Formatting   util.Optional[string]         `json:"formatting,omitzero"`
	Mappings     util.Optional[[]TERM_MAPPING] `json:"mappings,omitzero"`
	Language     util.Optional[CODE_PHRASE]    `json:"language,omitzero"`
	Encoding     util.Optional[CODE_PHRASE]    `json:"encoding,omitzero"`
	DefiningCode CODE_PHRASE                   `json:"defining_code"`
}

func (d *DV_CODED_TEXT) isDataValueModel() {}

func (d *DV_CODED_TEXT) isDvTextModel() {}

func (d *DV_CODED_TEXT) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_CODED_TEXT) SetModelName() {
	d.Type_ = util.Some(DV_CODED_TEXT_MODEL_NAME)
	d.DefiningCode.SetModelName()
	if d.Hyperlink.E {
		d.Hyperlink.V.SetModelName()
	}
	if d.Language.E {
		d.Language.V.SetModelName()
	}
	if d.Encoding.E {
		d.Encoding.V.SetModelName()
	}
	if d.Mappings.E {
		for i := range d.Mappings.V {
			d.Mappings.V[i].SetModelName()
		}
	}
}

func (d *DV_CODED_TEXT) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_CODED_TEXT_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          DV_CODED_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_CODED_TEXT_MODEL_NAME, d.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_CODED_TEXT_MODEL_NAME),
		})
	}

	// Validate defining_code
	attrPath = path + ".defining_code"
	errors = append(errors, d.DefiningCode.Validate(attrPath)...)

	// Validate formatting
	if d.Formatting.E {
		attrPath = path + ".formatting"
		validFormats := []string{"plain", "plain_no_newlines", "markdown"}
		if !slices.Contains(validFormats, d.Formatting.V) {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid formatting field: %s", d.Formatting.V),
				Recommendation: "Ensure formatting field is one of 'plain', 'plain_no_newlines', 'markdown'",
			})
		}
	}

	// Validate value
	attrPath = path + ".value"
	if d.Value == "" {
		errors = append(errors, util.ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field is required",
			Recommendation: "Ensure value field is not empty",
		})
	} else if len(d.Value) > 10000 {
		errors = append(errors, util.ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field exceeds maximum length of 10000 characters",
			Recommendation: "Ensure value field does not exceed 10000 characters",
		})
	}

	if d.Formatting.E && d.Formatting.V == "plain_no_newlines" {
		if strings.ContainsAny(d.Value, "\n\r") {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        "value field contains newlines but formatting is 'plain_no_newlines'",
				Recommendation: "Ensure value field does not contain newlines when formatting is 'plain_no_newlines'",
			})
		}
	}

	// Validate mappings
	if d.Mappings.E {
		attrPath = path + ".mappings"
		for i, v := range d.Mappings.V {
			errors = append(errors, v.Validate(fmt.Sprintf("%s[%d]", attrPath, i))...)
		}
	}

	// Validate language
	if d.Language.E {
		attrPath = path + ".language"
		if !terminology.IsValidLanguageTerminologyID(d.Language.V.TerminologyID.Value) {
			attrPath = path + ".terminology_id.value"
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language terminology ID: %s", d.Language.V.TerminologyID.Value),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}

		if !terminology.IsValidLanguageCode(d.Language.V.CodeString) {
			attrPath = path + ".code_string"
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language code: %s", d.Language.V.CodeString),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}
		errors = append(errors, d.Language.V.Validate(attrPath)...)
	}

	// Validate encoding
	if d.Encoding.E {
		attrPath = path + ".encoding"
		if !terminology.IsValidCharsetTerminologyID(d.Encoding.V.TerminologyID.Value) {
			attrPath = path + ".terminology_id.value"
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding terminology ID: %s", d.Encoding.V.TerminologyID.Value),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}

		if !terminology.IsValidCharset(d.Encoding.V.CodeString) {
			attrPath = path + ".code_string"
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding charset: %s", d.Encoding.V.CodeString),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}
		errors = append(errors, d.Encoding.V.Validate(attrPath)...)
	}

	return errors
}
