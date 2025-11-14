package openehr

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

// ========== Abstract of OBJECT_ID ==========

var _ util.ReferenceModel = (*X_OBJECT_ID)(nil)

// Abstract
type X_OBJECT_ID struct {
	Value util.ReferenceModel
}

func (x X_OBJECT_ID) HasModelName() bool {
	return x.Value.HasModelName()
}

func (x X_OBJECT_ID) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var valuePath string

	// Abstract model requires _type to be defined
	if !x.HasModelName() {
		valuePath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          OBJECT_REF_MODEL_NAME,
			Path:           valuePath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	return append(errs, x.Value.Validate(path)...)
}

func (o X_OBJECT_ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

func (o *X_OBJECT_ID) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.MetaType
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
	case "":
		return fmt.Errorf("missing OBJECT_ID _type field")
	default:
		return fmt.Errorf("OBJECT_ID unexpected _type %s", t)
	}

	return json.Unmarshal(data, o.Value)
}
