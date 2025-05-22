package json

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	err := marshal(reflect.ValueOf(v), &buf)
	return buf.Bytes(), err
}

func marshal(rv reflect.Value, buf *bytes.Buffer) error {
	switch rv.Kind() {
	case reflect.Bool:
		{
			var err error
			if rv.Bool() {
				_, err = buf.Write([]byte("true"))
			} else {
				_, err = buf.Write([]byte("false"))
			}
			return err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			_, err := buf.WriteString(strconv.Itoa(int(rv.Int())))
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		{
			_, err := buf.WriteString(strconv.FormatUint(uint64(rv.Float()), 10))
			return err
		}
	case reflect.Float32, reflect.Float64:
		{
			_, err := buf.WriteString(strconv.FormatFloat(rv.Float(), 'f', -1, 64))
			return err
		}
	case reflect.String:
		{
			buf.WriteByte('"')
			_, err := buf.WriteString(rv.String())
			buf.WriteByte('"')
			return err
		}
	case reflect.Pointer:
		{
			if rv.IsNil() {
				return errors.New("cannot handle nil pointers")
			}

			if rv.Type().NumMethod() > 0 && rv.CanInterface() {
				if u, ok := rv.Interface().(Marshaler); ok {
					bytes, err := u.Marshal()
					if err != nil {
						return err
					}

					if len(bytes) == 0 {
						return nil
					}

					buf.Write(bytes)
					return nil
				}
			}

			return marshal(rv.Elem(), buf)
		}
	case reflect.Array:
		{
			if rv.Type().NumMethod() > 0 && rv.CanInterface() {
				if u, ok := rv.Interface().(Marshaler); ok {
					bytes, err := u.Marshal()
					if err != nil {
						return err
					}

					if len(bytes) == 0 {
						return nil
					}

					buf.Write(bytes)
					return nil
				}
			}

			buf.WriteByte('[')
			n := rv.Len()
			for i := 0; i < n; i++ {
				if i > 0 {
					buf.WriteByte(',')
				}
				if err := marshal(rv.Index(i), buf); err != nil {
					return err
				}
			}
			buf.WriteByte(']')
			return nil
		}
	case reflect.Struct:
		{
			if rv.Type().NumMethod() > 0 && rv.CanInterface() {
				if u, ok := rv.Interface().(Marshaler); ok {
					bytes, err := u.Marshal()
					if err != nil {
						return err
					}

					if len(bytes) == 0 {
						return nil
					}

					buf.Write(bytes)
					return nil
				}
			}

			buf.WriteByte('{')
			rt := rv.Type()

			addComma := false
			for i := range rt.NumField() {
				field := rt.Field(i)
				tag := field.Tag.Get("json")
				if tag == "" {
					return fmt.Errorf("openehr: Marshal missing 'json' tag for field %s in struct %s", field.Type, rv.Type().Name())
				}

				fv := rv.Field(i)
				var fieldBuf bytes.Buffer
				if err := marshal(fv, &fieldBuf); err != nil {
					return err
				}

				if fieldBuf.Len() == 0 {
					continue
				}

				if addComma {
					buf.WriteByte(',')
				}
				addComma = true

				buf.WriteByte('"')
				buf.WriteString(tag)
				buf.WriteByte('"')
				buf.WriteByte(':')

				buf.Write(fieldBuf.Bytes())
			}
			buf.WriteByte('}')
			return nil
		}
	default:
		{
			return fmt.Errorf("openehr: Marshal with unsupported type %s", rv.Type().String())
		}
	}
}
