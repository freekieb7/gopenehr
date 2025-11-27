package model

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ITEM_MODEL_NAME string = "ITEM"

// Abstract
type ITEM struct {
	Type_            utils.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                      `json:"name"`
	ArchetypeNodeID  string                         `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
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
	Validate(path string) util.ValidateError
}

type X_ITEM struct {
	Value ItemModel
}

func (i *X_ITEM) SetModelName() {
	i.Value.SetModelName()
}

func (i *X_ITEM) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Abstract model requires _type to be defined
	if !i.Value.HasModelName() {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:   "ITEM",
			Path:    attrPath,
			Message: "missing _type field",
		})
	}

	return validateErr
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
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          ITEM_MODEL_NAME,
					Path:           "$.**._type",
					Message:        fmt.Sprintf("unexpected ITEM _type %s", t),
					Recommendation: "Ensure _type field is one of the known ITEM subtypes",
				},
			},
		}
	}

	return json.Unmarshal(data, i.Value)
}
