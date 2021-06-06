package di

import (
	"context"
	"reflect"
)

type contextStrategy struct{}

func (contextStrategy) furnishValue(ctx context.Context, v reflect.Value) (bool, error) {
	if v.Type() != contextType {
		return false, nil
	}
	v.Set(reflect.ValueOf(ctx))
	return true, nil
}
