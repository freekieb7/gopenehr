package openehr

import (
	"fmt"
	"slices"
	"strings"

	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_CODED_TEXT_MODEL_NAME string = "DV_CODED_TEXT"

var _ util.ReferenceModel = (*DV_CODED_TEXT)(nil)

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

func (d DV_CODED_TEXT) HasModelName() bool {
	return d.Type_.IsSet()
}

func (d DV_CODED_TEXT) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError

	// Validate _type
	if d.Type_.IsSet() && d.Type_.Unwrap() != DV_CODED_TEXT_MODEL_NAME {
		errors = append(errors, util.ValidationError{
			Model:          DV_CODED_TEXT_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_CODED_TEXT_MODEL_NAME, d.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_CODED_TEXT_MODEL_NAME),
		})
	}

	// Validate defining_code
	attrPath := path + ".defining_code"
	errors = append(errors, d.DefiningCode.Validate(attrPath)...)

	// Validate formatting
	if d.Formatting.IsSet() {
		attrPath = path + ".formatting"
		validFormats := []string{"plain", "plain_no_newlines", "markdown"}
		if !slices.Contains(validFormats, d.Formatting.Unwrap()) {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid formatting field: %s", d.Formatting.Unwrap()),
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

	if d.Formatting.IsSet() && d.Formatting.Unwrap() == "plain_no_newlines" {
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
	if d.Mappings.IsSet() {
		attrPath = path + ".mappings"
		for i, v := range d.Mappings.Unwrap() {
			errors = append(errors, v.Validate(fmt.Sprintf("%s[%d]", attrPath, i))...)
		}
	}

	// Validate language
	if d.Language.IsSet() {
		attrPath = path + ".language"
		v := d.Language.Unwrap()
		if !terminology.IsValidLanguageTerminologyID(v.TerminologyId.Value) {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language field: %s", d.Language.Unwrap()),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}

		if !terminology.IsValidLanguageCode(v.CodeString) {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language field: %s", d.Language.Unwrap()),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}
		errors = append(errors, d.Language.Unwrap().Validate(attrPath)...)
	}

	// Validate encoding
	if d.Encoding.IsSet() {
		attrPath = path + ".encoding"
		v := d.Encoding.Unwrap()
		if !terminology.IsValidCharsetTerminologyID(v.TerminologyId.Value) {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding field: %s", d.Encoding.Unwrap()),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}

		if !terminology.IsValidCharset(v.CodeString) {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding field: %s", d.Encoding.Unwrap()),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}
		errors = append(errors, d.Encoding.Unwrap().Validate(attrPath)...)
	}

	return errors
}
