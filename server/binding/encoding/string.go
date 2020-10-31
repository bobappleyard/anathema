package encoding

import (
	"github.com/bobappleyard/anathema/a"
	"reflect"
)

type StringEncoding struct {
	a.Service
}

func (e *StringEncoding) Accept(t reflect.Type) bool {
	return t.Kind() == reflect.String
}

func (e *StringEncoding) Decode(s string) (reflect.Value, error) {
	return reflect.ValueOf(e), nil
}
