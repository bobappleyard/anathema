package binding

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Str string
	Int int
}

func testBinding() Binding {
	b := Fields().ForStruct(reflect.TypeOf(testStruct{}))
	return b.Slice([]string{"Str", "Int"})
}

func TestDecode(t *testing.T) {
	b := testBinding()
	var v testStruct
	b.FromStrings([]string{"1", "2"}, reflect.ValueOf(&v))
	if v.Int != 2 {
		t.Fail()
	}
	if v.Str != "1" {
		t.Fail()
	}
}

func TestEncode(t *testing.T) {
	b := testBinding()
	ss, _ := b.ToStrings(reflect.ValueOf(&testStruct{
		Str: "hello",
		Int: 10,
	}).Elem())
	if len(ss) != 2 {
		t.FailNow()
	}
	if ss[0] != "hello" {
		t.Fail()
	}
	if ss[1] != "10" {
		t.Fail()
	}
}
