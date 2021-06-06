package di

import (
	"context"
	"reflect"
)

type sliceStrategy struct {
	factories *factoryRepository
}

func (s *sliceStrategy) furnishValue(ctx context.Context, v reflect.Value) (bool, error) {
	if v.Kind() != reflect.Slice {
		return false, nil
	}

	items := reflect.New(v.Type()).Elem()

	err := s.furnishSlices(ctx, v.Type(), items)
	if err != nil {
		return true, err
	}

	err = s.furnishIndividuals(ctx, v.Type().Elem(), items)
	if err != nil {
		return true, err
	}

	v.Set(items)
	return true, nil
}

func (s *sliceStrategy) furnishSlices(ctx context.Context, t reflect.Type, items reflect.Value) error {
	for _, f := range s.factories.findFactories(ctx, t) {
		s, err := s.factories.triggerFactory(ctx, f)
		if err != nil {
			return err
		}

		items.Set(reflect.AppendSlice(items, s))
	}

	return nil
}

func (s *sliceStrategy) furnishIndividuals(ctx context.Context, t reflect.Type, items reflect.Value) error {
	for _, f := range s.factories.findFactories(ctx, t) {
		s, err := s.factories.triggerFactory(ctx, f)
		if err != nil {
			return err
		}

		items.Set(reflect.Append(items, s))
	}

	return nil
}
