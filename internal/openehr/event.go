package openehr

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const EVENT_MODEL_NAME string = "EVENT"

type EventDataModel interface {
	isEventModel()
	HasModelName() bool
	SetModelName()
	Validate(path string) []util.ValidationError
}

// Abstract
type EVENT struct {
	Type_            util.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                       `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Time             DV_DATE_TIME                    `json:"time"`
	State            util.Optional[X_ITEM_STRUCTURE] `json:"state,omitzero"`
	Data             X_ITEM_STRUCTURE                `json:"data"`
}

func (e EVENT) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal abstract EVENT type")
}

func (e *EVENT) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("cannot unmarshal abstract EVENT type")
}

// ======== Union of EVENT subtypes ========

type X_EVENT struct {
	Value EventDataModel
}

func (x *X_EVENT) SetModelName() {
	x.Value.SetModelName()
}

func (x X_EVENT) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Abstract model requires _type to be defined
	if !x.Value.HasModelName() {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          EVENT_MODEL_NAME,
			Path:           attrPath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	errs = append(errs, x.Value.Validate(path)...)

	return errs
}

func (x X_EVENT) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Value)
}

func (x *X_EVENT) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.Type_
	switch t {
	case POINT_EVENT_MODEL_NAME:
		x.Value = new(POINT_EVENT)
	case INTERVAL_EVENT_MODEL_NAME:
		x.Value = new(INTERVAL_EVENT)
	case "":
		return fmt.Errorf("missing EVENT _type field")
	default:
		return fmt.Errorf("EVENT unexpected _type %s", t)
	}

	return json.Unmarshal(data, x.Value)
}
