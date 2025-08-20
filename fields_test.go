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
		d *float64
	}

	v := test{
		a: "A",
		b: "B",
		c: 123,
	}

	nameGetter, err := reflectx.FieldsByName[test]("a", "b", "c", "d")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nameGetter(&v)...)

	ptrGetter, err := reflectx.FieldsPtrByName[test]("a", "b", "c", "d")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ptrGetter(&v)...)

	_, err = reflectx.FieldsByTag[test]("sql")
	if err == nil {
		t.Fatal("tag reflecter fail")
	}

	tagGetter, err := reflectx.FieldsByTag[test]("get")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range tagGetter(&v) {
		t.Logf("%+v", f)
	}
}
