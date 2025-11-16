package openehr

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

// Abstract
type ITEM struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
}

func (i ITEM) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal abstract ITEM type")
}

func (i *ITEM) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("cannot unmarshal abstract ITEM type")
}

// ======== Union of ITEM subtypes ========

type ItemModel interface {
	isItemModel()
	HasModelName() bool
	SetModelName()
	Validate(path string) []util.ValidationError
}

type X_ITEM struct {
	Value ItemModel
}

func (i *X_ITEM) SetModelName() {
	i.Value.SetModelName()
}

func (i *X_ITEM) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Abstract model requires _type to be defined
	if !i.Value.HasModelName() {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:   "ITEM",
			Path:    attrPath,
			Message: "missing _type field",
		})
	}

	return errs
}

func (i X_ITEM) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Value)
}

func (i *X_ITEM) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.Type_
	switch t {
	case CLUSTER_MODEL_NAME:
		i.Value = new(CLUSTER)
	case ELEMENT_MODEL_NAME:
		i.Value = new(ELEMENT)
	case "":
		return fmt.Errorf("missing ITEM _type field")
	default:
		return fmt.Errorf("ITEM unexpected _type %s", t)
	}

	return json.Unmarshal(data, i.Value)
}
