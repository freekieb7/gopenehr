package openehr

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ITEM_STRUCTURE_MODEL_NAME string = "ITEM_STRUCTURE"

// Abstract
type ITEM_STRUCTURE struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[*FEEDER_AUDIT]  `json:"feeder_audit,omitzero"`
}

func (i ITEM_STRUCTURE) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal abstract ITEM_STRUCTURE type")
}

func (i *ITEM_STRUCTURE) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("cannot unmarshal abstract ITEM_STRUCTURE type")
}

// ======== Union of ITEM_STRUCTURE subtypes ========

type ItemStructureModel interface {
	isItemStructureModel()
	HasModelName() bool
	SetModelName()
	Validate(path string) []util.ValidationError
}

type X_ITEM_STRUCTURE struct {
	Value ItemStructureModel
}

func (x *X_ITEM_STRUCTURE) SetModelName() {
	x.Value.SetModelName()
}

func (x X_ITEM_STRUCTURE) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Abstract model requires _type to be defined
	if !x.Value.HasModelName() {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          ITEM_STRUCTURE_MODEL_NAME,
			Path:           attrPath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	return append(errs, x.Value.Validate(path)...)
}

func (x X_ITEM_STRUCTURE) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Value)
}

func (x *X_ITEM_STRUCTURE) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.MetaType
	switch t {
	case ITEM_SINGLE_MODEL_NAME:
		x.Value = new(ITEM_SINGLE)
	case ITEM_LIST_MODEL_NAME:
		x.Value = new(ITEM_LIST)
	case ITEM_TABLE_MODEL_NAME:
		x.Value = new(ITEM_TABLE)
	case ITEM_TREE_MODEL_NAME:
		x.Value = new(ITEM_TREE)
	case "":
		return fmt.Errorf("missing ITEM_STRUCTURE _type field")
	default:
		return fmt.Errorf("ITEM_STRUCTURE unexpected _type %s", t)
	}

	return json.Unmarshal(data, x.Value)
}
