package rm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
	"github.com/google/uuid"
)

const UID_BASED_ID_MODEL_NAME string = "UID_BASED_ID"

// Abstract type
type UID_BASED_ID struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Value string                 `json:"value"`
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
	Validate(path string) util.ValidateError
}

type X_UID_BASED_ID struct {
	Value UIDBasedIDModel
}

func (x *X_UID_BASED_ID) SetModelName() {
	x.Value.SetModelName()
}

func (x *X_UID_BASED_ID) Validate(path string) util.ValidateError {
	if x.Value == nil {
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          UID_BASED_ID_MODEL_NAME,
					Path:           path,
					Message:        "value is not known UID_BASED_ID subtype",
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
			Model:          UID_BASED_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	validateErr.Errs = append(validateErr.Errs, x.Value.Validate(path).Errs...)
	return validateErr
}

func (a X_UID_BASED_ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Value)
}

func (a *X_UID_BASED_ID) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.Type_
	switch t {
	case HIER_OBJECT_ID_MODEL_NAME:
		a.Value = new(HIER_OBJECT_ID)
	case OBJECT_VERSION_ID_MODEL_NAME:
		a.Value = new(OBJECT_VERSION_ID)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          UID_BASED_ID_MODEL_NAME,
					Path:           "$.**._type",
					Message:        fmt.Sprintf("unexpected UID_BASED_ID _type %s", t),
					Recommendation: "Ensure _type field is one of the known UID_BASED_ID subtypes",
				},
			},
		}
	}

	return json.Unmarshal(data, a.Value)
}

func (a *X_UID_BASED_ID) ValueAsString() string {
	switch v := a.Value.(type) {
	case *HIER_OBJECT_ID:
		return v.Value
	case *OBJECT_VERSION_ID:
		return v.Value
	default:
		return ""
	}
}

func (a *X_UID_BASED_ID) UUID() uuid.UUID {
	switch v := a.Value.(type) {
	case *HIER_OBJECT_ID:
		return uuid.MustParse(v.Value)
	case *OBJECT_VERSION_ID:
		return uuid.MustParse(strings.Split(v.Value, "::")[0])
	default:
		return uuid.Nil
	}
}
