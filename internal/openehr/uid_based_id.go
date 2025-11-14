package openehr

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const UID_BASED_ID_MODEL_NAME string = "UID_BASED_ID"

// Abstract type
type UID_BASED_ID struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
	Value string                `json:"value"`
}

func (u UID_BASED_ID) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal abstract UID_BASED_ID type")
}

func (u *UID_BASED_ID) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("cannot unmarshal abstract UID_BASED_ID type")
}

// ========== Union of UID_BASED_ID ==========

type UIDBasedIDModel interface {
	isUidBasedIDModel()
	HasModelName() bool
	SetModelName()
	Validate(path string) []util.ValidationError
}

type X_UID_BASED_ID struct {
	Value UIDBasedIDModel
}

func (x *X_UID_BASED_ID) SetModelName() {
	x.Value.SetModelName()
}

func (x X_UID_BASED_ID) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError

	// Abstract model requires _type to be defined
	if !x.Value.HasModelName() {
		errs = append(errs, util.ValidationError{
			Model:          OBJECT_REF_MODEL_NAME,
			Path:           "._type",
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	return append(errs, x.Value.Validate(path)...)
}

func (a X_UID_BASED_ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Value)
}

func (a *X_UID_BASED_ID) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.MetaType
	switch t {
	case HIER_OBJECT_ID_MODEL_NAME:
		a.Value = new(HIER_OBJECT_ID)
	case OBJECT_VERSION_ID_MODEL_NAME:
		a.Value = new(OBJECT_VERSION_ID)
	case "":
		return fmt.Errorf("missing DV_TEXT _type field")
	default:
		return fmt.Errorf("DV_TEXT unexpected _type %s", t)
	}

	return json.Unmarshal(data, a.Value)
}
