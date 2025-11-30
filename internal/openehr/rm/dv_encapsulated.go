package rm

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_ENCAPSULATED_MODEL_NAME string = "DV_ENCAPSULATED"

// Abstact
type DV_ENCAPSULATED struct {
	Type_    utils.Optional[string]      `json:"_type,omitzero"`
	Charset  utils.Optional[CODE_PHRASE] `json:"charset,omitzero"`
	Language utils.Optional[CODE_PHRASE] `json:"language,omitzero"`
}

func (d DV_ENCAPSULATED) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal abstract DV_ENCAPSULATED type")
}

func (d *DV_ENCAPSULATED) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("cannot unmarshal abstract DV_ENCAPSULATED type")
}

// ========== Union of DV_ENCAPSULATED ==========

type DvEncapsulatedModel interface {
	HasModelName() bool
	SetModelName()
	Validate(path string) util.ValidateError
}

// Abstract
type X_DV_ENCAPSULATED struct {
	Value DvEncapsulatedModel
}

func (x *X_DV_ENCAPSULATED) SetModelName() {
	x.Value.SetModelName()
}

func (x *X_DV_ENCAPSULATED) Validate(path string) util.ValidateError {
	if x.Value == nil {
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          DV_ENCAPSULATED_MODEL_NAME,
					Path:           path,
					Message:        "value is not known DV_ENCAPSULATED subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}

	var validateErr util.ValidateError
	var attrPath string

	// Abstract model requires _type to be defined
	if !x.Value.HasModelName() {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_ENCAPSULATED_MODEL_NAME,
			Path:           attrPath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	validateErr.Errs = append(validateErr.Errs, x.Value.Validate(path).Errs...)
	return validateErr
}

func (a X_DV_ENCAPSULATED) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Value)
}

func (a *X_DV_ENCAPSULATED) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.Type_
	switch t {
	case DV_MULTIMEDIA_MODEL_NAME:
		a.Value = new(DV_MULTIMEDIA)
	case DV_EHR_URI_MODEL_NAME:
		a.Value = new(DV_EHR_URI)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          DV_ENCAPSULATED_MODEL_NAME,
					Path:           "$.**._type",
					Message:        fmt.Sprintf("unexpected DV_ENCAPSULATED _type %s", t),
					Recommendation: "Ensure _type field is one of the known DV_ENCAPSULATED subtypes",
				},
			},
		}
	}

	return json.Unmarshal(data, a.Value)
}
