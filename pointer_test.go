package reflectx_test

import (
	"reflect"
	"testing"
)

func TestPointer(t *testing.T) {
	var a *int

	obj := reflect.TypeOf(a)

	t.Log(obj, obj.Name(), obj.Elem().Name())
}
