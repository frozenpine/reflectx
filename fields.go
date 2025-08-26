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
	isPtr  bool
	tag    string
}

func getNamedFields[T any](fields ...string) ([]fieldOffset, error) {
	count := len(fields)
	if count == 0 {
		return nil, errors.New("no field name specified")
	}

	obj := new(T)
	objType := reflect.TypeOf(obj).Elem()

	switch objType.Kind() {
	case reflect.Struct:
	default:
		return nil, errors.New("object must be a struct")
	}

	fieldOffsets := make([]fieldOffset, count)

	for idx, name := range fields {
		if field, ok := objType.FieldByName(name); !ok {
			return nil, fmt.Errorf("no field[%s] for %s", name, objType.Name())
		} else {
			fieldOffsets[idx] = fieldOffset{
				offset: field.Offset,
				typ:    field.Type,
				isPtr:  field.Type.Kind() == reflect.Pointer,
			}

			if fieldOffsets[idx].isPtr {
				typeName := fieldOffsets[idx].typ.
					Elem().Name()

				if _, exist := typedPtrPools[typeName]; !exist {
					return nil, errors.New(
						"struct ptr field has no value pool registered",
					)
				}
			}
		}
	}

	return fieldOffsets, nil
}

func FieldsPtrByName[T any](fields ...string) (func(*T) []any, error) {
	fieldOffsets, err := getNamedFields[T](fields...)
	if err != nil {
		return nil, err
	}

	return func(data *T) []any {
		basePtr := uintptr(unsafe.Pointer(data))

		results := make([]interface{}, len(fieldOffsets))

		for idx, define := range fieldOffsets {
			v := reflect.NewAt(
				define.typ, unsafe.Pointer(basePtr+define.offset),
			)
			if define.isPtr {
				v = v.Elem()
			}

			r := v.Interface()
			if v.IsNil() || v.IsZero() {
				ptrBaseType := define.typ.Elem().Kind().String()
				r = typedPtrPools[ptrBaseType].Get()
			}
			results[idx] = r
		}

		return results
	}, nil
}

func FieldsByName[T any](fields ...string) (func(*T) []interface{}, error) {
	fieldOffsets, err := getNamedFields[T](fields...)
	if err != nil {
		return nil, err
	}

	return func(data *T) []interface{} {
		basePtr := reflect.Indirect(
			reflect.ValueOf(data),
		).Addr().Pointer()

		results := make([]interface{}, len(fieldOffsets))

		for idx, define := range fieldOffsets {
			results[idx] = reflect.Indirect(reflect.NewAt(
				define.typ, unsafe.Pointer(basePtr+define.offset),
			)).Interface()
		}

		return results
	}, nil
}

type TagField struct {
	Tag   string
	Value interface{}
}

func FieldsByTag[T any](tag string) (func(*T) []TagField, error) {
	if tag == "" {
		return nil, errors.New("no tag specifed")
	}

	obj := new(T)
	objType := reflect.TypeOf(obj).Elem()

	switch objType.Kind() {
	case reflect.Struct:
	default:
		return nil, errors.New("object must be a struct")
	}

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
			tag:    tagValue,
		})
	}

	if len(fieldOffsets) == 0 {
		return nil, fmt.Errorf("no field found with tag[%s]", tag)
	}

	return func(data *T) []TagField {
		basePtr := reflect.Indirect(
			reflect.ValueOf(data),
		).Addr().Pointer()

		result := make([]TagField, len(fieldOffsets))

		for idx, define := range fieldOffsets {
			result[idx] = TagField{
				Tag: define.tag,
				Value: reflect.Indirect(reflect.NewAt(
					define.typ, unsafe.Pointer(basePtr+define.offset),
				)).Interface(),
			}
		}

		return result
	}, nil
}
