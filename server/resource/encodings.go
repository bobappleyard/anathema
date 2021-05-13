package resource

import (
	"encoding"
	"github.com/bobappleyard/anathema/server/a"
	"reflect"
	"strconv"
)

type intEncoding struct {
	a.Service
}

func (e *intEncoding) Accept(t reflect.Type) bool {
	return t.Kind() == reflect.Int
}

func (e *intEncoding) Decode(s string, v reflect.Value) error {
	x, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	v.SetInt(int64(x))
	return nil
}

type stringEncoding struct {
	a.Service
}

func (e *stringEncoding) Accept(t reflect.Type) bool {
	return t.Kind() == reflect.String
}

func (e *stringEncoding) Decode(s string, v reflect.Value) error {
	v.SetString(s)
	return nil
}

type methodEncoding struct {
	a.Service
}

var unmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

func (e *methodEncoding) Accept(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
	}
	return t.AssignableTo(unmarshalerType)
}

func (e *methodEncoding) Decode(s string, v reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		v = v.Addr()
	}
	return v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(s))
}
