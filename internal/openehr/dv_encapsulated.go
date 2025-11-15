package openehr

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_ENCAPSULATED_MODEL_NAME string = "DV_ENCAPSULATED"

// Abstact
type DV_ENCAPSULATED struct {
	Type_    util.Optional[string]      `json:"_type,omitzero"`
	Charset  util.Optional[CODE_PHRASE] `json:"charset,omitzero"`
	Language util.Optional[CODE_PHRASE] `json:"language,omitzero"`
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
	Validate(path string) []util.ValidationError
}

// Abstract
type X_DV_ENCAPSULATED struct {
	Value DvEncapsulatedModel
}

func (x *X_DV_ENCAPSULATED) SetModelName() {
	x.Value.SetModelName()
}

func (x X_DV_ENCAPSULATED) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Abstract model requires _type to be defined
	if !x.Value.HasModelName() {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          DV_ENCAPSULATED_MODEL_NAME,
			Path:           attrPath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	return append(errs, x.Value.Validate(path)...)
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
	case "":
		return fmt.Errorf("missing DV_ENCAPSULATED _type field")
	default:
		return fmt.Errorf("DV_ENCAPSULATED unexpected _type %s", t)
	}

	return json.Unmarshal(data, a.Value)
}
