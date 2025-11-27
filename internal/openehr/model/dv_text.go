package model

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_TEXT_MODEL_NAME string = "DV_TEXT"

type DV_TEXT struct {
	Type_      utils.Optional[string]         `json:"_type,omitzero"`
	Value      string                         `json:"value"`
	Formatting utils.Optional[string]         `json:"formatting,omitzero"`
	Mappings   utils.Optional[[]TERM_MAPPING] `json:"mappings,omitzero"`
	Language   utils.Optional[CODE_PHRASE]    `json:"language,omitzero"`
	Encoding   utils.Optional[CODE_PHRASE]    `json:"encoding,omitzero"`
}

func (d *DV_TEXT) isDataValueModel() {}

func (d *DV_TEXT) isDvTextModel() {}

func (d *DV_TEXT) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_TEXT) SetModelName() {
	d.Type_ = utils.Some(DV_TEXT_MODEL_NAME)
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

func (d *DV_TEXT) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_TEXT_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
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
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
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
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field is required",
			Recommendation: "Ensure value field is not empty",
		})
	} else if len(d.Value) > 10000 {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field exceeds maximum length of 10000 characters",
			Recommendation: "Ensure value field does not exceed 10000 characters",
		})
	}

	if d.Formatting.E && d.Formatting.V == "plain_no_newlines" {
		if strings.ContainsAny(d.Value, "\n\r") {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
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
			validateErr.Errs = append(validateErr.Errs, d.Mappings.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate language
	if d.Language.E {
		attrPath = path + ".language"
		if !terminology.IsValidLanguageTerminologyID(d.Language.V.TerminologyID.Value) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language field: %s", d.Language.V.TerminologyID.Value),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}

		if !terminology.IsValidLanguageCode(d.Language.V.CodeString) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language field: %s", d.Language.V.CodeString),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}
		validateErr.Errs = append(validateErr.Errs, d.Language.V.Validate(attrPath).Errs...)
	}

	// Validate encoding
	if d.Encoding.E {
		attrPath = path + ".encoding"
		if !terminology.IsValidCharsetTerminologyID(d.Encoding.V.TerminologyID.Value) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding field: %s", d.Encoding.V.TerminologyID.Value),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}

		if !terminology.IsValidCharset(d.Encoding.V.CodeString) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding field: %s", d.Encoding.V.CodeString),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}
		validateErr.Errs = append(validateErr.Errs, d.Encoding.V.Validate(attrPath).Errs...)
	}

	return validateErr
}

// ========== Union of DV_TEXT ==========

type DvTextModel interface {
	isDvTextModel()
	HasModelName() bool
	SetModelName()
	Validate(path string) util.ValidateError
}

type X_DV_TEXT struct {
	Value DvTextModel
}

func (x X_DV_TEXT) SetModelName() {
	x.Value.SetModelName()
}

func (x X_DV_TEXT) Validate(path string) util.ValidateError {
	if x.Value == nil {
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          DV_TEXT_MODEL_NAME,
					Path:           path,
					Message:        "value is not known DV_TEXT subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}

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
