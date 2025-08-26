package reflectx

import (
	"errors"
	"reflect"
	"sync"
)

var (
	typedPtrPools = map[string]*sync.Pool{}

	baseTypeCreator = map[reflect.Kind]func() any{
		reflect.Bool:    func() any { return new(bool) },
		reflect.Uint8:   func() any { return new(uint8) },
		reflect.Int8:    func() any { return new(int8) },
		reflect.Uint:    func() any { return new(uint) },
		reflect.Int:     func() any { return new(int) },
		reflect.Uint16:  func() any { return new(uint16) },
		reflect.Int16:   func() any { return new(int16) },
		reflect.Uint32:  func() any { return new(uint32) },
		reflect.Int32:   func() any { return new(int32) },
		reflect.Uint64:  func() any { return new(uint64) },
		reflect.Int64:   func() any { return new(int64) },
		reflect.Float32: func() any { return new(float32) },
		reflect.Float64: func() any { return new(float64) },
		reflect.String:  func() any { return new(string) },
	}
)

func RegisterTypedPool[T any](pool *sync.Pool) error {
	if pool == nil {
		return errors.New("invalid pool for type")
	}

	obj := new(T)

	objType := reflect.TypeOf(obj).Elem()
	typeName := objType.Name()

	typedPtrPools[typeName] = pool

	return nil
}

func init() {
	for _, baseType := range []reflect.Kind{
		reflect.Bool, reflect.Uint8, reflect.Int8, reflect.Uint, reflect.Int,
		reflect.Uint16, reflect.Int16, reflect.Uint32, reflect.Int32,
		reflect.Uint64, reflect.Int64, reflect.Float32, reflect.Float64,
		reflect.String,
	} {
		creatorFn := baseTypeCreator[baseType]
		if creatorFn == nil {
			panic("no base type create found")
		}

		typedPtrPools[baseType.String()] = &sync.Pool{
			New: creatorFn,
		}
	}
}
