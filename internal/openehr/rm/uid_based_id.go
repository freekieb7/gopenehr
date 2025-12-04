package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const UID_BASED_ID_TYPE string = "UID_BASED_ID"

type UIDBasedIDKind int

const (
	UID_BASED_ID_kind_unknown UIDBasedIDKind = iota
	UID_BASED_ID_kind_HIER_OBJECT_ID
	UID_BASED_ID_kind_OBJECT_VERSION_ID
)

func (k UIDBasedIDKind) String() string {
	switch k {
	case UID_BASED_ID_kind_HIER_OBJECT_ID:
		return HIER_OBJECT_ID_TYPE
	case UID_BASED_ID_kind_OBJECT_VERSION_ID:
		return OBJECT_VERSION_ID_TYPE
	default:
		return "unknown"
	}
}

type UIDBasedIDUnion struct {
	Kind  UIDBasedIDKind
	Value any
}

func (u *UIDBasedIDUnion) SetModelName() {
	switch u.Kind {
	case UID_BASED_ID_kind_HIER_OBJECT_ID:
		u.Value.(*HIER_OBJECT_ID).SetModelName()
	case UID_BASED_ID_kind_OBJECT_VERSION_ID:
		u.Value.(*OBJECT_VERSION_ID).SetModelName()
	}
}

func (u *UIDBasedIDUnion) Validate(path string) util.ValidateError {
	switch u.Kind {
	case UID_BASED_ID_kind_HIER_OBJECT_ID:
		return u.Value.(*HIER_OBJECT_ID).Validate(path)
	case UID_BASED_ID_kind_OBJECT_VERSION_ID:
		return u.Value.(*OBJECT_VERSION_ID).Validate(path)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          UID_BASED_ID_TYPE,
					Path:           path,
					Message:        "value is not known UID_BASED_ID subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}
}

func (u UIDBasedIDUnion) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Value)
}

func (u *UIDBasedIDUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case HIER_OBJECT_ID_TYPE:
		u.Kind = UID_BASED_ID_kind_HIER_OBJECT_ID
		u.Value = &HIER_OBJECT_ID{}
	case OBJECT_VERSION_ID_TYPE:
		u.Kind = UID_BASED_ID_kind_OBJECT_VERSION_ID
		u.Value = &OBJECT_VERSION_ID{}
	default:
		u.Kind = UID_BASED_ID_kind_unknown
		return nil
	}

	return json.Unmarshal(data, u.Value)
}

func (u *UIDBasedIDUnion) HIER_OBJECT_ID() HIER_OBJECT_ID {
	if u.Kind == UID_BASED_ID_kind_HIER_OBJECT_ID {
		return *u.Value.(*HIER_OBJECT_ID)
	}
	return HIER_OBJECT_ID{}
}

func (u *UIDBasedIDUnion) OBJECT_VERSION_ID() OBJECT_VERSION_ID {
	if u.Kind == UID_BASED_ID_kind_OBJECT_VERSION_ID {
		return *u.Value.(*OBJECT_VERSION_ID)
	}
	return OBJECT_VERSION_ID{}
}

func UID_BASED_ID_from_HIER_OBJECT_ID(id *HIER_OBJECT_ID) UIDBasedIDUnion {
	id.Type_ = utils.Some(HIER_OBJECT_ID_TYPE)

	return UIDBasedIDUnion{
		Kind:  UID_BASED_ID_kind_HIER_OBJECT_ID,
		Value: id,
	}
}

func UID_BASED_ID_from_OBJECT_VERSION_ID(id *OBJECT_VERSION_ID) UIDBasedIDUnion {
	id.Type_ = utils.Some(OBJECT_VERSION_ID_TYPE)

	return UIDBasedIDUnion{
		Kind:  UID_BASED_ID_kind_OBJECT_VERSION_ID,
		Value: id,
	}
}
