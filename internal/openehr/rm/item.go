package rm

import (
	"github.com/bytedance/sonic"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ITEM_TYPE string = "ITEM"

type ItemKind int

const (
	ITEM_kind_unknown ItemKind = iota
	ITEM_kind_CLUSTER
	ITEM_kind_ELEMENT
)

type ItemUnion struct {
	Kind  ItemKind
	Value any
}

func (i *ItemUnion) SetModelName() {
	switch i.Kind {
	case ITEM_kind_CLUSTER:
		i.Value.(*CLUSTER).SetModelName()
	case ITEM_kind_ELEMENT:
		i.Value.(*ELEMENT).SetModelName()
	}
}

func (i *ItemUnion) Validate(path string) util.ValidateError {
	switch i.Kind {
	case ITEM_kind_CLUSTER:
		return i.Value.(*CLUSTER).Validate(path)
	case ITEM_kind_ELEMENT:
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
	return sonic.Marshal(i.Value)
}

func (i *ItemUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case CLUSTER_TYPE:
		i.Kind = ITEM_kind_CLUSTER
		i.Value = new(CLUSTER)
	case ELEMENT_TYPE:
		i.Kind = ITEM_kind_ELEMENT
		i.Value = new(ELEMENT)
	default:
		i.Kind = ITEM_kind_unknown
		return nil
	}

	return sonic.Unmarshal(data, i.Value)
}

func (i *ItemUnion) CLUSTER() *CLUSTER {
	if i.Kind != ITEM_kind_CLUSTER {
		return nil
	}
	return i.Value.(*CLUSTER)
}

func (i *ItemUnion) ELEMENT() *ELEMENT {
	if i.Kind != ITEM_kind_ELEMENT {
		return nil
	}
	return i.Value.(*ELEMENT)
}

func ITEM_from_CLUSTER(cluster CLUSTER) ItemUnion {
	cluster.Type_ = utils.Some(CLUSTER_TYPE)
	return ItemUnion{
		Kind:  ITEM_kind_CLUSTER,
		Value: &cluster,
	}
}

func ITEM_from_ELEMENT(element ELEMENT) ItemUnion {
	element.Type_ = utils.Some(ELEMENT_TYPE)
	return ItemUnion{
		Kind:  ITEM_kind_ELEMENT,
		Value: &element,
	}
}
