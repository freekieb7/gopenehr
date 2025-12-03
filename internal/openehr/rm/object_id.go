package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const OBJECT_ID_TYPE string = "OBJECT_ID"

type ObjectIDKind int

const (
	OBJECT_ID_kind_unknown ObjectIDKind = iota
	OBJECT_ID_kind_HIER_OBJECT_ID
	OBJECT_ID_kind_OBJECT_VERSION_ID
	OBJECT_ID_kind_ARCHETYPE_ID
	OBJECT_ID_kind_TEMPLATE_ID
	OBJECT_ID_kind_GENERIC_ID
)

type ObjectIDUnion struct {
	Kind  ObjectIDKind
	Value any
}

func (o *ObjectIDUnion) SetModelName() {
	switch o.Kind {
	case OBJECT_ID_kind_HIER_OBJECT_ID:
		o.Value.(*HIER_OBJECT_ID).SetModelName()
	case OBJECT_ID_kind_OBJECT_VERSION_ID:
		o.Value.(*OBJECT_VERSION_ID).SetModelName()
	case OBJECT_ID_kind_ARCHETYPE_ID:
		o.Value.(*ARCHETYPE_ID).SetModelName()
	case OBJECT_ID_kind_TEMPLATE_ID:
		o.Value.(*TEMPLATE_ID).SetModelName()
	case OBJECT_ID_kind_GENERIC_ID:
		o.Value.(*GENERIC_ID).SetModelName()
	}
}

func (o *ObjectIDUnion) Validate(path string) util.ValidateError {
	switch o.Kind {
	case OBJECT_ID_kind_HIER_OBJECT_ID:
		return o.Value.(*HIER_OBJECT_ID).Validate(path)
	case OBJECT_ID_kind_OBJECT_VERSION_ID:
		return o.Value.(*OBJECT_VERSION_ID).Validate(path)
	case OBJECT_ID_kind_ARCHETYPE_ID:
		return o.Value.(*ARCHETYPE_ID).Validate(path)
	case OBJECT_ID_kind_TEMPLATE_ID:
		return o.Value.(*TEMPLATE_ID).Validate(path)
	case OBJECT_ID_kind_GENERIC_ID:
		return o.Value.(*GENERIC_ID).Validate(path)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          OBJECT_ID_TYPE,
					Path:           path,
					Message:        "value is not known OBJECT_ID subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}
}

func (o ObjectIDUnion) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

func (o *ObjectIDUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case HIER_OBJECT_ID_TYPE:
		o.Kind = OBJECT_ID_kind_HIER_OBJECT_ID
		o.Value = &HIER_OBJECT_ID{}
	case OBJECT_VERSION_ID_TYPE:
		o.Kind = OBJECT_ID_kind_OBJECT_VERSION_ID
		o.Value = &OBJECT_VERSION_ID{}
	case ARCHETYPE_ID_TYPE:
		o.Kind = OBJECT_ID_kind_ARCHETYPE_ID
		o.Value = &ARCHETYPE_ID{}
	case TEMPLATE_ID_TYPE:
		o.Kind = OBJECT_ID_kind_TEMPLATE_ID
		o.Value = &TEMPLATE_ID{}
	case GENERIC_ID_TYPE:
		o.Kind = OBJECT_ID_kind_GENERIC_ID
		o.Value = &GENERIC_ID{}
	default:
		o.Kind = OBJECT_ID_kind_unknown
		return nil
	}

	return json.Unmarshal(data, o.Value)
}

func (o *ObjectIDUnion) HIER_OBJECT_ID() *HIER_OBJECT_ID {
	if o.Kind == OBJECT_ID_kind_HIER_OBJECT_ID {
		return o.Value.(*HIER_OBJECT_ID)
	}
	return nil
}

func (o *ObjectIDUnion) OBJECT_VERSION_ID() *OBJECT_VERSION_ID {
	if o.Kind == OBJECT_ID_kind_OBJECT_VERSION_ID {
		return o.Value.(*OBJECT_VERSION_ID)
	}
	return nil
}

func (o *ObjectIDUnion) ARCHETYPE_ID() *ARCHETYPE_ID {
	if o.Kind == OBJECT_ID_kind_ARCHETYPE_ID {
		return o.Value.(*ARCHETYPE_ID)
	}
	return nil
}

func (o *ObjectIDUnion) TEMPLATE_ID() *TEMPLATE_ID {
	if o.Kind == OBJECT_ID_kind_TEMPLATE_ID {
		return o.Value.(*TEMPLATE_ID)
	}
	return nil
}

func (o *ObjectIDUnion) GENERIC_ID() *GENERIC_ID {
	if o.Kind == OBJECT_ID_kind_GENERIC_ID {
		return o.Value.(*GENERIC_ID)
	}
	return nil
}

func OBJECT_ID_from_HIER_OBJECT_ID(hierObjectID HIER_OBJECT_ID) ObjectIDUnion {
	hierObjectID.Type_ = utils.Some(HIER_OBJECT_ID_TYPE)

	return ObjectIDUnion{
		Kind:  OBJECT_ID_kind_HIER_OBJECT_ID,
		Value: &hierObjectID,
	}
}

func OBJECT_ID_from_OBJECT_VERSION_ID(objectVersionID OBJECT_VERSION_ID) ObjectIDUnion {
	objectVersionID.Type_ = utils.Some(OBJECT_VERSION_ID_TYPE)

	return ObjectIDUnion{
		Kind:  OBJECT_ID_kind_OBJECT_VERSION_ID,
		Value: &objectVersionID,
	}
}

func OBJECT_ID_from_ARCHETYPE_ID(archetypeID ARCHETYPE_ID) ObjectIDUnion {
	archetypeID.Type_ = utils.Some(ARCHETYPE_ID_TYPE)

	return ObjectIDUnion{
		Kind:  OBJECT_ID_kind_ARCHETYPE_ID,
		Value: &archetypeID,
	}
}

func OBJECT_ID_from_TEMPLATE_ID(templateID TEMPLATE_ID) ObjectIDUnion {
	templateID.Type_ = utils.Some(TEMPLATE_ID_TYPE)

	return ObjectIDUnion{
		Kind:  OBJECT_ID_kind_TEMPLATE_ID,
		Value: &templateID,
	}
}

func OBJECT_ID_from_GENERIC_ID(genericID GENERIC_ID) ObjectIDUnion {
	genericID.Type_ = utils.Some(GENERIC_ID_TYPE)

	return ObjectIDUnion{
		Kind:  OBJECT_ID_kind_GENERIC_ID,
		Value: &genericID,
	}
}
