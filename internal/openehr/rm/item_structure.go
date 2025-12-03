package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ITEM_STRUCTURE_TYPE string = "ITEM_STRUCTURE"

type ItemStructureKind int

const (
	ITEM_STRUCTURE_kind_unknown ItemStructureKind = iota
	ITEM_STRUCTURE_kind_ITEM_SINGLE
	ITEM_STRUCTURE_kind_ITEM_LIST
	ITEM_STRUCTURE_kind_ITEM_TABLE
	ITEM_STRUCTURE_kind_ITEM_TREE
)

type ItemStructureUnion struct {
	Kind  ItemStructureKind
	Value any
}

func (i *ItemStructureUnion) SetModelName() {
	switch i.Kind {
	case ITEM_STRUCTURE_kind_ITEM_SINGLE:
		i.Value.(*ITEM_SINGLE).SetModelName()
	case ITEM_STRUCTURE_kind_ITEM_LIST:
		i.Value.(*ITEM_LIST).SetModelName()
	case ITEM_STRUCTURE_kind_ITEM_TABLE:
		i.Value.(*ITEM_TABLE).SetModelName()
	case ITEM_STRUCTURE_kind_ITEM_TREE:
		i.Value.(*ITEM_TREE).SetModelName()
	}
}

func (x *ItemStructureUnion) Validate(path string) util.ValidateError {
	switch x.Kind {
	case ITEM_STRUCTURE_kind_ITEM_SINGLE:
		return x.Value.(*ITEM_SINGLE).Validate(path)
	case ITEM_STRUCTURE_kind_ITEM_LIST:
		return x.Value.(*ITEM_LIST).Validate(path)
	case ITEM_STRUCTURE_kind_ITEM_TABLE:
		return x.Value.(*ITEM_TABLE).Validate(path)
	case ITEM_STRUCTURE_kind_ITEM_TREE:
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
		i.Kind = ITEM_STRUCTURE_kind_ITEM_SINGLE
		i.Value = &ITEM_SINGLE{}
	case ITEM_LIST_TYPE:
		i.Kind = ITEM_STRUCTURE_kind_ITEM_LIST
		i.Value = &ITEM_LIST{}
	case ITEM_TABLE_TYPE:
		i.Kind = ITEM_STRUCTURE_kind_ITEM_TABLE
		i.Value = &ITEM_TABLE{}
	case ITEM_TREE_TYPE:
		i.Kind = ITEM_STRUCTURE_kind_ITEM_TREE
		i.Value = &ITEM_TREE{}
	default:
		i.Kind = ITEM_STRUCTURE_kind_unknown
		return nil
	}

	return json.Unmarshal(data, i.Value)
}

func (i *ItemStructureUnion) ITEM_SINGLE() *ITEM_SINGLE {
	if i.Kind != ITEM_STRUCTURE_kind_ITEM_SINGLE {
		return nil
	}
	return i.Value.(*ITEM_SINGLE)
}

func (i *ItemStructureUnion) ITEM_LIST() *ITEM_LIST {
	if i.Kind != ITEM_STRUCTURE_kind_ITEM_LIST {
		return nil
	}
	return i.Value.(*ITEM_LIST)
}

func (i *ItemStructureUnion) ITEM_TABLE() *ITEM_TABLE {
	if i.Kind != ITEM_STRUCTURE_kind_ITEM_TABLE {
		return nil
	}
	return i.Value.(*ITEM_TABLE)
}

func (i *ItemStructureUnion) ITEM_TREE() *ITEM_TREE {
	if i.Kind != ITEM_STRUCTURE_kind_ITEM_TREE {
		return nil
	}
	return i.Value.(*ITEM_TREE)
}

func ITEM_STRUCTURE_from_ITEM_SINGLE(itemSingle ITEM_SINGLE) ItemStructureUnion {
	itemSingle.Type_ = utils.Some(ITEM_SINGLE_TYPE)
	return ItemStructureUnion{
		Kind:  ITEM_STRUCTURE_kind_ITEM_SINGLE,
		Value: &itemSingle,
	}
}

func ITEM_STRUCTURE_from_ITEM_LIST(itemList ITEM_LIST) ItemStructureUnion {
	itemList.Type_ = utils.Some(ITEM_LIST_TYPE)
	return ItemStructureUnion{
		Kind:  ITEM_STRUCTURE_kind_ITEM_LIST,
		Value: &itemList,
	}
}

func ITEM_STRUCTURE_from_ITEM_TABLE(itemTable ITEM_TABLE) ItemStructureUnion {
	itemTable.Type_ = utils.Some(ITEM_TABLE_TYPE)
	return ItemStructureUnion{
		Kind:  ITEM_STRUCTURE_kind_ITEM_TABLE,
		Value: &itemTable,
	}
}

func ITEM_STRUCTURE_from_ITEM_TREE(itemTree ITEM_TREE) ItemStructureUnion {
	itemTree.Type_ = utils.Some(ITEM_TREE_TYPE)
	return ItemStructureUnion{
		Kind:  ITEM_STRUCTURE_kind_ITEM_TREE,
		Value: &itemTree,
	}
}
