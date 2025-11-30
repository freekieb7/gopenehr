package rm

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const CONTENT_ITEM_MODEL_NAME string = "CONTENT_ITEM"

// Abstract
type CONTENT_ITEM struct {
	Type_            utils.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                      `json:"name"`
	ArchetypeNodeID  string                         `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
}

func (c CONTENT_ITEM) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal abstract CONTENT_ITEM type")
}

func (c *CONTENT_ITEM) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("cannot unmarshal abstract CONTENT_ITEM type")
}

// ======== Union of CONTENT_ITEM subtypes ========

type ContentItemModel interface {
	isContentItemModel()
	HasModelName() bool
	SetModelName()
	Validate(path string) util.ValidateError
}

type X_CONTENT_ITEM struct {
	Value ContentItemModel
}

func (x *X_CONTENT_ITEM) SetModelName() {
	x.Value.SetModelName()
}

func (x *X_CONTENT_ITEM) Validate(path string) util.ValidateError {
	if x.Value == nil {
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          CONTENT_ITEM_MODEL_NAME,
					Path:           path,
					Message:        "value is not known CONTENT_ITEM subtype",
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
			Model:          CONTENT_ITEM_MODEL_NAME,
			Path:           attrPath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}
	validateErr.Errs = append(validateErr.Errs, x.Value.Validate(path).Errs...)

	return validateErr
}

func (x X_CONTENT_ITEM) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Value)
}

func (x *X_CONTENT_ITEM) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.Type_
	switch t {
	case SECTION_MODEL_NAME:
		x.Value = new(SECTION)
	case ADMIN_ENTRY_MODEL_NAME:
		x.Value = new(ADMIN_ENTRY)
	case OBSERVATION_MODEL_NAME:
		x.Value = new(OBSERVATION)
	case EVALUATION_MODEL_NAME:
		x.Value = new(EVALUATION)
	case INSTRUCTION_MODEL_NAME:
		x.Value = new(INSTRUCTION)
	case ACTIVITY_MODEL_NAME:
		x.Value = new(ACTIVITY)
	case ACTION_MODEL_NAME:
		x.Value = new(ACTION)
	case GENERIC_ENTRY_MODEL_NAME:
		x.Value = new(GENERIC_ENTRY)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          CONTENT_ITEM_MODEL_NAME,
					Path:           "$.**._type",
					Message:        fmt.Sprintf("unexpected CONTENT_ITEM _type %s", t),
					Recommendation: "Ensure _type field is one of the known CONTENT_ITEM subtypes",
				},
			},
		}
	}

	return json.Unmarshal(data, x.Value)
}
