package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ITEM_TYPE string = "ITEM"

type ItemKind int

const (
	ItemKind_Unknown ItemKind = iota
	ItemKind_CLUSTER
	ItemKind_ELEMENT
)

type ItemUnion struct {
	Kind  ItemKind
	Value any
}

func (i *ItemUnion) SetModelName() {
	switch i.Kind {
	case ItemKind_CLUSTER:
		i.Value.(*CLUSTER).SetModelName()
	case ItemKind_ELEMENT:
		i.Value.(*ELEMENT).SetModelName()
	}
}

func (i *ItemUnion) Validate(path string) util.ValidateError {
	switch i.Kind {
	case ItemKind_CLUSTER:
		return i.Value.(*CLUSTER).Validate(path)
	case ItemKind_ELEMENT:
		return i.Value.(*ELEMENT).Validate(path)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          ITEM_TYPE,
					Path:           path,
					Message:        "value is not known ITEM subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}
}

func (i ItemUnion) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Value)
}

func (i *ItemUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case CLUSTER_TYPE:
		i.Kind = ItemKind_CLUSTER
		i.Value = new(CLUSTER)
	case ELEMENT_TYPE:
		i.Kind = ItemKind_ELEMENT
		i.Value = new(ELEMENT)
	default:
		i.Kind = ItemKind_Unknown
		return nil
	}

	return json.Unmarshal(data, i.Value)
}
