package di

import (
	"context"
	"reflect"
)

type structStrategy struct{}

func (structStrategy) furnishValue(ctx context.Context, v reflect.Value) (bool, error) {
	if v.Kind() != reflect.Struct {
		return false, nil
	}

	t := v.Type()
	for i := t.NumField() - 1; i >= 0; i-- {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		if field.Tag != "" {
			continue
		}
		err := FurnishValue(ctx, v.FieldByIndex(field.Index))
		if err != nil {
			return true, err
		}
	}

	return true, nil
}
