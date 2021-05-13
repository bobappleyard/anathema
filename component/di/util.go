package di

import (
	"reflect"
)

func injectedCall(s Process, f reflect.Value) ([]reflect.Value, error) {
	t := f.Type()
	in := make([]reflect.Value, t.NumIn())
	for i := range in {
		v, err := s.RequireValue(t.In(i))
		if err != nil {
			return nil, err
		}
		in[i] = v
	}
	return f.Call(in), nil
}
