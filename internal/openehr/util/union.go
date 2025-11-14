package util

import "reflect"

type UnionValue interface {
	GetBaseType() reflect.Type
}
