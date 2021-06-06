package di

import (
	"context"
	"errors"
	"reflect"
)

var (
	ErrNoFurnisher      = errors.New("no furnisher has been installed")
	ErrTooManyFactories = errors.New("too many factories")
	ErrUnsupported      = errors.New("unsupported type")
)

func Furnish(ctx context.Context, ref interface{}) error {
	return FurnishValue(ctx, reflect.ValueOf(ref).Elem())
}

func FurnishValue(ctx context.Context, v reflect.Value) error {
	f := fromContext(ctx)
	if f == nil {
		return ErrNoFurnisher
	}
	return f.furnishValue(ctx, v)
}

func FurnishArgs(ctx context.Context, m reflect.Value) ([]reflect.Value, error) {
	t := m.Type()
	args := make([]reflect.Value, t.NumIn())
	for i := range args {
		a := reflect.New(t.In(i)).Elem()
		err := FurnishValue(ctx, a)
		if err != nil {
			return nil, err
		}
		args[i] = a
	}
	return args, nil
}

var furnisherKey = reflect.TypeOf(new(furnisher))

type furnisher struct {
	strategies []strategy
}

type strategy interface {
	furnishValue(ctx context.Context, v reflect.Value) (bool, error)
}

func newFurnisher(factories []factory) *furnisher {
	fs := &factoryRepository{factories: factories}

	return &furnisher{
		strategies: []strategy{
			contextStrategy{},
			&sliceStrategy{fs},
			&factoryStrategy{fs},
			ptrStrategy{},
			structStrategy{},
		},
	}
}

func fromContext(ctx context.Context) *furnisher {
	f := ctx.Value(furnisherKey)
	if f == nil {
		return nil
	}
	return f.(*furnisher)
}

func toContext(ctx context.Context, f *furnisher) context.Context {
	return context.WithValue(ctx, furnisherKey, f)
}

func (f *furnisher) furnishValue(ctx context.Context, v reflect.Value) error {
	for _, s := range f.strategies {
		done, err := s.furnishValue(ctx, v)
		if done {
			return err
		}
	}
	return ErrUnsupported
}
