package openehr

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_TEXT_MODEL_NAME string = "DV_TEXT"

type DV_TEXT struct {
	Type_      util.Optional[string]         `json:"_type,omitzero"`
	Value      string                        `json:"value"`
	Formatting util.Optional[string]         `json:"formatting,omitzero"`
	Mappings   util.Optional[[]TERM_MAPPING] `json:"mappings,omitzero"`
	Language   util.Optional[CODE_PHRASE]    `json:"language,omitzero"`
	Encoding   util.Optional[CODE_PHRASE]    `json:"encoding,omitzero"`
}

func (d DV_TEXT) isDataValueModel() {}

func (d DV_TEXT) isDvTextModel() {}

func (d DV_TEXT) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_TEXT) SetModelName() {
	d.Type_ = util.Some(DV_TEXT_MODEL_NAME)
	if d.Mappings.E {
		for i := range d.Mappings.V {
			d.Mappings.V[i].SetModelName()
		}
	}
	if d.Language.E {
		d.Language.V.SetModelName()
	}
	if d.Encoding.E {
		d.Encoding.V.SetModelName()
	}
}

func (d DV_TEXT) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_TEXT_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_TEXT_MODEL_NAME, d.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_TEXT_MODEL_NAME),
		})
	}

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
		for i := range d.Mappings.V {
			attrPath := fmt.Sprintf("%s[%d]", attrPath, i)
			errors = append(errors, d.Mappings.V[i].Validate(attrPath)...)
		}
	}

	// Validate language
	if d.Language.E {
		attrPath = path + ".language"
		if !terminology.IsValidLanguageTerminologyID(d.Language.V.TerminologyID.Value) {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language field: %s", d.Language.V.TerminologyID.Value),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}

		if !terminology.IsValidLanguageCode(d.Language.V.CodeString) {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language field: %s", d.Language.V.CodeString),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}
		errors = append(errors, d.Language.V.Validate(attrPath)...)
	}

	// Validate encoding
	if d.Encoding.E {
		attrPath = path + ".encoding"
		if !terminology.IsValidCharsetTerminologyID(d.Encoding.V.TerminologyID.Value) {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding field: %s", d.Encoding.V.TerminologyID.Value),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}

		if !terminology.IsValidCharset(d.Encoding.V.CodeString) {
			errors = append(errors, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding field: %s", d.Encoding.V.CodeString),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}
		errors = append(errors, d.Encoding.V.Validate(attrPath)...)
	}

	return errors
}

// ========== Union of DV_TEXT ==========

type DvTextModel interface {
	isDvTextModel()
	HasModelName() bool
	SetModelName()
	Validate(path string) []util.ValidationError
}

type X_DV_TEXT struct {
	Value DvTextModel
}

func (x X_DV_TEXT) SetModelName() {
	x.Value.SetModelName()
}

func (x X_DV_TEXT) Validate(path string) []util.ValidationError {
	return x.Value.Validate(path)
}

func (a X_DV_TEXT) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Value)
}

func (a *X_DV_TEXT) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.Type_
	switch t {
	case DV_TEXT_MODEL_NAME, "":
		a.Value = new(DV_TEXT)
	case DV_CODED_TEXT_MODEL_NAME:
		a.Value = new(DV_CODED_TEXT)
	default:
		return fmt.Errorf("DV_TEXT unexpected _type %s", t)
	}

	return json.Unmarshal(data, a.Value)
}
