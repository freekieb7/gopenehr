package utils

import (
	"database/sql/driver"
	"encoding/json"
)

type OptionalValue interface {
	IsSet() bool
}

type Optional[T any] struct {
	V T
	E bool
}

func Some[T any](v T) Optional[T] {
	return Optional[T]{V: v, E: true}
}

func None[T any]() Optional[T] {
	return Optional[T]{}
}

// func (o Optional[T]) IsZero() bool {
// 	return !o.E
// }

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.E {
		return []byte("null"), nil
	}
	return json.Marshal(o.V)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.E = false
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.E = true
	o.V = v
	return nil
}

// Scan implements the SQL driver.Scanner interface.
func (o *Optional[T]) Scan(value any) error {
	if value == nil {
		o.E = false
		return nil
	}

	var v T
	switch t := any(&v).(type) {
	case interface{ Scan(any) error }:
		if err := t.Scan(value); err != nil {
			return err
		}
	default:
		v = value.(T)
	}

	o.V = v
	o.E = true

	return nil
}

// Value implements the driver Valuer interface.
func (o Optional[T]) Value() (driver.Value, error) {
	if !o.E {
		return nil, nil
	}
	switch t := any(o.V).(type) {
	case interface{ Value() (any, error) }:
		return t.Value()
	default:
		return o.V, nil
	}
}
