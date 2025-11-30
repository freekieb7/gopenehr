package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ITEM_STRUCTURE_TYPE string = "ITEM_STRUCTURE"

type ItemStructureKind int

const (
	ItemStructureKind_Unknown ItemStructureKind = iota
	ItemStructureKind_ITEM_SINGLE
	ItemStructureKind_ITEM_LIST
	ItemStructureKind_ITEM_TABLE
	ItemStructureKind_ITEM_TREE
)

type ItemStructureUnion struct {
	Kind  ItemStructureKind
	Value any
}

func (i *ItemStructureUnion) SetModelName() {
	switch i.Kind {
	case ItemStructureKind_ITEM_SINGLE:
		i.Value.(*ITEM_SINGLE).SetModelName()
	case ItemStructureKind_ITEM_LIST:
		i.Value.(*ITEM_LIST).SetModelName()
	case ItemStructureKind_ITEM_TABLE:
		i.Value.(*ITEM_TABLE).SetModelName()
	case ItemStructureKind_ITEM_TREE:
		i.Value.(*ITEM_TREE).SetModelName()
	}
}

func (x *ItemStructureUnion) Validate(path string) util.ValidateError {
	switch x.Kind {
	case ItemStructureKind_ITEM_SINGLE:
		return x.Value.(*ITEM_SINGLE).Validate(path)
	case ItemStructureKind_ITEM_LIST:
		return x.Value.(*ITEM_LIST).Validate(path)
	case ItemStructureKind_ITEM_TABLE:
		return x.Value.(*ITEM_TABLE).Validate(path)
	case ItemStructureKind_ITEM_TREE:
		return x.Value.(*ITEM_TREE).Validate(path)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          ITEM_STRUCTURE_TYPE,
					Path:           path,
					Message:        "value is not known ITEM_STRUCTURE subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}
}

func (i ItemStructureUnion) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Value)
}

func (i *ItemStructureUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case ITEM_SINGLE_TYPE:
		i.Kind = ItemStructureKind_ITEM_SINGLE
		i.Value = &ITEM_SINGLE{}
	case ITEM_LIST_TYPE:
		i.Kind = ItemStructureKind_ITEM_LIST
		i.Value = &ITEM_LIST{}
	case ITEM_TABLE_TYPE:
		i.Kind = ItemStructureKind_ITEM_TABLE
		i.Value = &ITEM_TABLE{}
	case ITEM_TREE_TYPE:
		i.Kind = ItemStructureKind_ITEM_TREE
		i.Value = &ITEM_TREE{}
	default:
		i.Kind = ItemStructureKind_Unknown
		return nil
	}

	return json.Unmarshal(data, i.Value)
}
