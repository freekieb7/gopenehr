package util

import "reflect"

type AbstractType interface {
	GetAbstractType() reflect.Type
}
