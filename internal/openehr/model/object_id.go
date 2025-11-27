package model

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const OBJECT_ID_MODEL_NAME string = "OBJECT_ID"

// Abstract
type OBJECT_ID struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Value string                 `json:"value"`
}

func (o OBJECT_ID) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal abstract OBJECT_ID type")
}

func (o *OBJECT_ID) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("cannot unmarshal abstract OBJECT_ID type")
}

// ========== Union of OBJECT_ID ==========

type ObjectIDModel interface {
	isObjectIDModel()
	HasModelName() bool
	GetModelName() string
	SetModelName()
	Validate(path string) util.ValidateError
}

type X_OBJECT_ID struct {
	Value ObjectIDModel
}

func (x *X_OBJECT_ID) SetModelName() {
	x.Value.SetModelName()
}

func (x *X_OBJECT_ID) Validate(path string) util.ValidateError {
	if x.Value == nil {
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          OBJECT_ID_MODEL_NAME,
					Path:           path,
					Message:        "value is not known OBJECT_ID subtype",
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
			Model:          OBJECT_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	validateErr.Errs = append(validateErr.Errs, x.Value.Validate(path).Errs...)
	return validateErr
}

func (o X_OBJECT_ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

func (o *X_OBJECT_ID) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.Type_
	switch t {
	case HIER_OBJECT_ID_MODEL_NAME:
		o.Value = new(HIER_OBJECT_ID)
	case OBJECT_VERSION_ID_MODEL_NAME:
		o.Value = new(OBJECT_VERSION_ID)
	case ARCHETYPE_ID_MODEL_NAME:
		o.Value = new(ARCHETYPE_ID)
	case TEMPLATE_ID_MODEL_NAME:
		o.Value = new(TEMPLATE_ID)
	case GENERIC_ID_MODEL_NAME:
		o.Value = new(GENERIC_ID)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          OBJECT_ID_MODEL_NAME,
					Path:           "$.**._type",
					Message:        fmt.Sprintf("unexpected OBJECT_ID _type '%s'", t),
					Recommendation: "Ensure _type field is one of the known OBJECT_ID subtypes",
				},
			},
		}
	}

	return json.Unmarshal(data, o.Value)
}
