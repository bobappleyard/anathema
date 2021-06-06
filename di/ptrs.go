package di

import (
	"context"
	"reflect"
)

type ptrStrategy struct{}

func (ptrStrategy) furnishValue(ctx context.Context, v reflect.Value) (bool, error) {
	if v.Kind() != reflect.Ptr {
		return false, nil
	}
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return true, FurnishValue(ctx, v.Elem())
}
