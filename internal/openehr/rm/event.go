package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const EVENT_TYPE string = "EVENT"

type EventKind int

const (
	EventKind_Unknown EventKind = iota
	EventKind_POINT_EVENT
	EventKind_INTERVAL_EVENT
)

type EventUnion struct {
	Kind  EventKind
	Value any
}

func (e *EventUnion) SetModelName() {
	switch e.Kind {
	case EventKind_POINT_EVENT:
		e.Value.(*POINT_EVENT).SetModelName()
	case EventKind_INTERVAL_EVENT:
		e.Value.(*INTERVAL_EVENT).SetModelName()
	}
}

func (e *EventUnion) Validate(path string) util.ValidateError {
	switch e.Kind {
	case EventKind_POINT_EVENT:
		return e.Value.(*POINT_EVENT).Validate(path)
	case EventKind_INTERVAL_EVENT:
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
		e.Kind = EventKind_POINT_EVENT
		e.Value = &POINT_EVENT{}
	case INTERVAL_EVENT_TYPE:
		e.Kind = EventKind_INTERVAL_EVENT
		e.Value = &INTERVAL_EVENT{}
	default:
		e.Kind = EventKind_Unknown
		return nil
	}

	return json.Unmarshal(data, e.Value)
}
