package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const EVENT_TYPE string = "EVENT"

type EventKind int

const (
	EVENT_kind_unknown EventKind = iota
	EVENT_kind_POINT_EVENT
	EVENT_kind_INTERVAL_EVENT
)

type EventUnion struct {
	Kind  EventKind
	Value any
}

func (e *EventUnion) SetModelName() {
	switch e.Kind {
	case EVENT_kind_POINT_EVENT:
		e.Value.(*POINT_EVENT).SetModelName()
	case EVENT_kind_INTERVAL_EVENT:
		e.Value.(*INTERVAL_EVENT).SetModelName()
	}
}

func (e *EventUnion) Validate(path string) util.ValidateError {
	switch e.Kind {
	case EVENT_kind_POINT_EVENT:
		return e.Value.(*POINT_EVENT).Validate(path)
	case EVENT_kind_INTERVAL_EVENT:
		return e.Value.(*INTERVAL_EVENT).Validate(path)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          EVENT_TYPE,
					Path:           path,
					Message:        "value is not known EVENT subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}
}

func (e EventUnion) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Value)
}

func (e *EventUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case POINT_EVENT_TYPE:
		e.Kind = EVENT_kind_POINT_EVENT
		e.Value = &POINT_EVENT{}
	case INTERVAL_EVENT_TYPE:
		e.Kind = EVENT_kind_INTERVAL_EVENT
		e.Value = &INTERVAL_EVENT{}
	default:
		e.Kind = EVENT_kind_unknown
		return nil
	}

	return json.Unmarshal(data, e.Value)
}

func (o *EventUnion) POINT_EVENT() *POINT_EVENT {
	if o.Kind == EVENT_kind_POINT_EVENT {
		return o.Value.(*POINT_EVENT)
	}
	return nil
}

func (o *EventUnion) INTERVAL_EVENT() *INTERVAL_EVENT {
	if o.Kind == EVENT_kind_INTERVAL_EVENT {
		return o.Value.(*INTERVAL_EVENT)
	}
	return nil
}

func EVENT_from_POINT_EVENT(pointEvent POINT_EVENT) EventUnion {
	pointEvent.Type_ = utils.Some(POINT_EVENT_TYPE)
	return EventUnion{
		Kind:  EVENT_kind_POINT_EVENT,
		Value: &pointEvent,
	}
}

func EVENT_from_INTERVAL_EVENT(intervalEvent INTERVAL_EVENT) EventUnion {
	intervalEvent.Type_ = utils.Some(INTERVAL_EVENT_TYPE)
	return EventUnion{
		Kind:  EVENT_kind_INTERVAL_EVENT,
		Value: &intervalEvent,
	}
}
