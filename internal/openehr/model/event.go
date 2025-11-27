package model

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const EVENT_MODEL_NAME string = "EVENT"

type EventDataModel interface {
	isEventModel()
	HasModelName() bool
	SetModelName()
	Validate(path string) util.ValidateError
}

// Abstract
type EVENT struct {
	Type_            utils.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                        `json:"name"`
	ArchetypeNodeID  string                           `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Time             DV_DATE_TIME                     `json:"time"`
	State            utils.Optional[X_ITEM_STRUCTURE] `json:"state,omitzero"`
	Data             X_ITEM_STRUCTURE                 `json:"data"`
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

func (x *X_EVENT) Validate(path string) util.ValidateError {
	if x.Value == nil {
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          EVENT_MODEL_NAME,
					Path:           path,
					Message:        "value is not known EVENT subtype",
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
			Model:          EVENT_MODEL_NAME,
			Path:           attrPath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	validateErr.Errs = append(validateErr.Errs, x.Value.Validate(path).Errs...)

	return validateErr
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
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          EVENT_MODEL_NAME,
					Path:           "$.**._type",
					Message:        fmt.Sprintf("unexpected EVENT _type %s", t),
					Recommendation: "Ensure _type field is one of the known EVENT subtypes",
				},
			},
		}
	}

	return json.Unmarshal(data, x.Value)
}
