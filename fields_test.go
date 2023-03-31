package reflectx_test

import (
	"testing"

	"github.com/frozenpine/reflectx"
)

func TestGetFields(t *testing.T) {
	type test struct {
		a string `get:"a"`
		b string
		c int `get:"c"`
	}

	v := test{
		a: "A",
		b: "B",
		c: 123,
	}

	nameGetter, err := reflectx.FieldsByName[test]("a", "b", "c")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nameGetter(&v)...)

	_, err = reflectx.FieldsByTag[test]("sql")
	if err == nil {
		t.Fatal("tag reflecter fail")
	}

	tagGetter, err := reflectx.FieldsByTag[test]("get")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tagGetter(&v)...)
}
