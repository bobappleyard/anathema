package di

import (
	"context"
	"reflect"
	"testing"
)

func TestFactoryStrategy(t *testing.T) {
	var x int
	s := &factoryStrategy{&factoryRepository{factories: []factory{
		{
			forType: reflect.TypeOf(0),
			impl: func(ctx context.Context) (reflect.Value, error) {
				return reflect.ValueOf(10), nil
			},
		},
	}}}
	done, err := s.furnishValue(context.Background(), reflect.ValueOf(&x).Elem())
	if !done {
		t.Fail()
	}
	if err != nil {
		t.Error(err)
	}
	if x != 10 {
		t.Fail()
	}
}

func TestMultiFactoryStrategy(t *testing.T) {
	var x int
	s := &factoryStrategy{&factoryRepository{factories: []factory{
		{
			forType: reflect.TypeOf(0),
			impl: func(ctx context.Context) (reflect.Value, error) {
				return reflect.ValueOf(10), nil
			},
		},
		{
			forType: reflect.TypeOf(0),
			impl: func(ctx context.Context) (reflect.Value, error) {
				return reflect.ValueOf(10), nil
			},
		},
	}}}
	_, err := s.furnishValue(context.Background(), reflect.ValueOf(&x).Elem())
	if err != ErrTooManyFactories {
		t.Error(err)
	}
}
