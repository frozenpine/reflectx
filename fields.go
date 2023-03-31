package reflectx

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type fieldOffset struct {
	offset uintptr
	typ    reflect.Type
}

func FieldsByName[T any](fields ...string) (func(*T) []interface{}, error) {
	count := len(fields)
	if count == 0 {
		return nil, errors.New("no field name specified")
	}

	obj := new(T)
	objType := reflect.TypeOf(obj).Elem()

	fieldOffsets := make([]fieldOffset, count)

	for idx, name := range fields {
		if field, ok := objType.FieldByName(name); !ok {
			return nil, fmt.Errorf("no field[%s] for %s", name, objType.Name())
		} else {
			fieldOffsets[idx] = fieldOffset{
				offset: field.Offset,
				typ:    field.Type,
			}
		}
	}

	return func(data *T) []interface{} {
		basePtr := reflect.Indirect(
			reflect.ValueOf(data),
		).Addr().Pointer()

		results := make([]interface{}, count)

		for idx, define := range fieldOffsets {
			results[idx] = reflect.Indirect(reflect.NewAt(
				define.typ, unsafe.Pointer(basePtr+define.offset),
			)).Interface()
		}

		return results
	}, nil
}

func FieldsByTag[T any](tag string) (func(*T) []interface{}, error) {
	if tag == "" {
		return nil, errors.New("no tag specifed")
	}

	obj := new(T)
	objType := reflect.TypeOf(obj).Elem()

	fieldOffsets := []fieldOffset{}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)

		tagValue := field.Tag.Get(tag)

		if tagValue == "" {
			continue
		}

		fieldOffsets = append(fieldOffsets, fieldOffset{
			offset: field.Offset,
			typ:    field.Type,
		})
	}

	if len(fieldOffsets) == 0 {
		return nil, fmt.Errorf("no field found with tag[%s]", tag)
	}

	return func(data *T) []interface{} {
		basePtr := reflect.Indirect(
			reflect.ValueOf(data),
		).Addr().Pointer()

		result := make([]interface{}, len(fieldOffsets))

		for idx, define := range fieldOffsets {
			result[idx] = reflect.Indirect(reflect.NewAt(
				define.typ, unsafe.Pointer(basePtr+define.offset),
			)).Interface()
		}

		return result
	}, nil
}
