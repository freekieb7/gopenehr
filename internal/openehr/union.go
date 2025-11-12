package openehr

import "reflect"

type UnionValue interface {
	GetBaseType() reflect.Type
}
