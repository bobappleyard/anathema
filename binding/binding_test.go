package binding

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Str string
	Int int
}

func TestDecode(t *testing.T) {
	b := ForStruct(reflect.TypeOf(testStruct{})).Slice([]string{"Str", "Int"})
	s, _ := b.FromStrings([]string{"1", "2"})
	v := s.Interface().(testStruct)
	if v.Int != 2 {
		t.Fail()
	}
	if v.Str != "1" {
		t.Fail()
	}
}

func TestEncode(t *testing.T) {
	b := ForStruct(reflect.TypeOf(testStruct{})).Slice([]string{"Str", "Int"})
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
