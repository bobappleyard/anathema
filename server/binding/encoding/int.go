package encoding

import (
	"github.com/bobappleyard/anathema/a"
	"reflect"
	"strconv"
)

type IntEncoding struct {
	a.Service
}

func (e *IntEncoding) Accept(t reflect.Type) bool {
	return t.Kind() == reflect.Int
}

func (e *IntEncoding) Decode(s string) (reflect.Value, error) {
	x, err := strconv.Atoi(s)
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(x), nil
}
