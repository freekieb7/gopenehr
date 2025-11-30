package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const UID_BASED_ID_TYPE string = "UID_BASED_ID"

type UIDBasedIDKind int

const (
	UIDBasedIDKind_Unknown UIDBasedIDKind = iota
	UIDBasedIDKind_HIER_OBJECT_ID
	UIDBasedIDKind_OBJECT_VERSION_ID
)

type UIDBasedIDUnion struct {
	Kind  UIDBasedIDKind
	Value any
}

func (u *UIDBasedIDUnion) SetModelName() {
	switch u.Kind {
	case UIDBasedIDKind_HIER_OBJECT_ID:
		u.Value.(*HIER_OBJECT_ID).SetModelName()
	case UIDBasedIDKind_OBJECT_VERSION_ID:
		u.Value.(*OBJECT_VERSION_ID).SetModelName()
	}
}

func (u *UIDBasedIDUnion) Validate(path string) util.ValidateError {
	switch u.Kind {
	case UIDBasedIDKind_HIER_OBJECT_ID:
		return u.Value.(*HIER_OBJECT_ID).Validate(path)
	case UIDBasedIDKind_OBJECT_VERSION_ID:
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
		u.Kind = UIDBasedIDKind_HIER_OBJECT_ID
		u.Value = new(HIER_OBJECT_ID)
	case OBJECT_VERSION_ID_TYPE:
		u.Kind = UIDBasedIDKind_OBJECT_VERSION_ID
		u.Value = new(OBJECT_VERSION_ID)
	default:
		u.Kind = UIDBasedIDKind_Unknown
		return nil
	}

	return json.Unmarshal(data, u.Value)
}

func (u *UIDBasedIDUnion) HierObjectID() *HIER_OBJECT_ID {
	if u.Kind == UIDBasedIDKind_HIER_OBJECT_ID {
		return u.Value.(*HIER_OBJECT_ID)
	}
	return nil
}

func (u *UIDBasedIDUnion) ObjectVersionID() *OBJECT_VERSION_ID {
	if u.Kind == UIDBasedIDKind_OBJECT_VERSION_ID {
		return u.Value.(*OBJECT_VERSION_ID)
	}
	return nil
}

func UIDBasedIDFromHierObjectID(id *HIER_OBJECT_ID) UIDBasedIDUnion {
	id.Type_ = utils.Some(HIER_OBJECT_ID_TYPE)

	return UIDBasedIDUnion{
		Kind:  UIDBasedIDKind_HIER_OBJECT_ID,
		Value: id,
	}
}

func UIDBasedIDFromObjectVersionID(id *OBJECT_VERSION_ID) UIDBasedIDUnion {
	id.Type_ = utils.Some(OBJECT_VERSION_ID_TYPE)

	return UIDBasedIDUnion{
		Kind:  UIDBasedIDKind_OBJECT_VERSION_ID,
		Value: id,
	}
}
