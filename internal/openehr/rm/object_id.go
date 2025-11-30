package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const OBJECT_ID_TYPE string = "OBJECT_ID"

type ObjectIDKind int

const (
	ObjectIDKind_Unknown ObjectIDKind = iota
	ObjectIDKind_HIER_OBJECT_ID
	ObjectIDKind_OBJECT_VERSION_ID
	ObjectIDKind_ARCHETYPE_ID
	ObjectIDKind_TEMPLATE_ID
	ObjectIDKind_GENERIC_ID
)

type ObjectIDUnion struct {
	Kind  ObjectIDKind
	Value any
}

func (o *ObjectIDUnion) SetModelName() {
	switch o.Kind {
	case ObjectIDKind_HIER_OBJECT_ID:
		o.Value.(*HIER_OBJECT_ID).SetModelName()
	case ObjectIDKind_OBJECT_VERSION_ID:
		o.Value.(*OBJECT_VERSION_ID).SetModelName()
	case ObjectIDKind_ARCHETYPE_ID:
		o.Value.(*ARCHETYPE_ID).SetModelName()
	case ObjectIDKind_TEMPLATE_ID:
		o.Value.(*TEMPLATE_ID).SetModelName()
	case ObjectIDKind_GENERIC_ID:
		o.Value.(*GENERIC_ID).SetModelName()
	}
}

func (o *ObjectIDUnion) Validate(path string) util.ValidateError {
	switch o.Kind {
	case ObjectIDKind_HIER_OBJECT_ID:
		return o.Value.(*HIER_OBJECT_ID).Validate(path)
	case ObjectIDKind_OBJECT_VERSION_ID:
		return o.Value.(*OBJECT_VERSION_ID).Validate(path)
	case ObjectIDKind_ARCHETYPE_ID:
		return o.Value.(*ARCHETYPE_ID).Validate(path)
	case ObjectIDKind_TEMPLATE_ID:
		return o.Value.(*TEMPLATE_ID).Validate(path)
	case ObjectIDKind_GENERIC_ID:
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
		o.Kind = ObjectIDKind_HIER_OBJECT_ID
		o.Value = &HIER_OBJECT_ID{}
	case OBJECT_VERSION_ID_TYPE:
		o.Kind = ObjectIDKind_OBJECT_VERSION_ID
		o.Value = &OBJECT_VERSION_ID{}
	case ARCHETYPE_ID_TYPE:
		o.Kind = ObjectIDKind_ARCHETYPE_ID
		o.Value = &ARCHETYPE_ID{}
	case TEMPLATE_ID_TYPE:
		o.Kind = ObjectIDKind_TEMPLATE_ID
		o.Value = &TEMPLATE_ID{}
	case GENERIC_ID_TYPE:
		o.Kind = ObjectIDKind_GENERIC_ID
		o.Value = &GENERIC_ID{}
	default:
		o.Kind = ObjectIDKind_Unknown
		return nil
	}

	return json.Unmarshal(data, o.Value)
}

func (o *ObjectIDUnion) HierObjectID() *HIER_OBJECT_ID {
	if o.Kind == ObjectIDKind_HIER_OBJECT_ID {
		return o.Value.(*HIER_OBJECT_ID)
	}
	return nil
}

func (o *ObjectIDUnion) ObjectVersionID() *OBJECT_VERSION_ID {
	if o.Kind == ObjectIDKind_OBJECT_VERSION_ID {
		return o.Value.(*OBJECT_VERSION_ID)
	}
	return nil
}

func (o *ObjectIDUnion) ArchetypeID() *ARCHETYPE_ID {
	if o.Kind == ObjectIDKind_ARCHETYPE_ID {
		return o.Value.(*ARCHETYPE_ID)
	}
	return nil
}

func (o *ObjectIDUnion) TemplateID() *TEMPLATE_ID {
	if o.Kind == ObjectIDKind_TEMPLATE_ID {
		return o.Value.(*TEMPLATE_ID)
	}
	return nil
}

func (o *ObjectIDUnion) GenericID() *GENERIC_ID {
	if o.Kind == ObjectIDKind_GENERIC_ID {
		return o.Value.(*GENERIC_ID)
	}
	return nil
}

func ObjectIDFromHierObjectID(hierObjectID HIER_OBJECT_ID) ObjectIDUnion {
	hierObjectID.Type_ = utils.Some(HIER_OBJECT_ID_TYPE)

	return ObjectIDUnion{
		Kind:  ObjectIDKind_HIER_OBJECT_ID,
		Value: &hierObjectID,
	}
}

func ObjectIDFromObjectVersionID(objectVersionID OBJECT_VERSION_ID) ObjectIDUnion {
	objectVersionID.Type_ = utils.Some(OBJECT_VERSION_ID_TYPE)

	return ObjectIDUnion{
		Kind:  ObjectIDKind_OBJECT_VERSION_ID,
		Value: &objectVersionID,
	}
}

func ObjectIDFromArchetypeID(archetypeID ARCHETYPE_ID) ObjectIDUnion {
	archetypeID.Type_ = utils.Some(ARCHETYPE_ID_TYPE)

	return ObjectIDUnion{
		Kind:  ObjectIDKind_ARCHETYPE_ID,
		Value: &archetypeID,
	}
}

func ObjectIDFromTemplateID(templateID TEMPLATE_ID) ObjectIDUnion {
	templateID.Type_ = utils.Some(TEMPLATE_ID_TYPE)

	return ObjectIDUnion{
		Kind:  ObjectIDKind_TEMPLATE_ID,
		Value: &templateID,
	}
}

func ObjectIDFromGenericID(genericID GENERIC_ID) ObjectIDUnion {
	genericID.Type_ = utils.Some(GENERIC_ID_TYPE)

	return ObjectIDUnion{
		Kind:  ObjectIDKind_GENERIC_ID,
		Value: &genericID,
	}
}
