package di

import (
	"context"
	"reflect"
	"testing"
)

func TestSliceStrategyElemFactory(t *testing.T) {
	var xs []int
	s := sliceStrategy{&factoryRepository{factories: []factory{
		{
			forType: reflect.TypeOf(0),
			impl: func(ctx context.Context) (reflect.Value, error) {
				return reflect.ValueOf(10), nil
			},
		},
	}}}
	done, err := s.furnishValue(context.Background(), reflect.ValueOf(&xs).Elem())
	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Fail()
	}
	if !reflect.DeepEqual(xs, []int{10}) {
		t.Fail()
	}
}

func TestSliceStrategySliceFactory(t *testing.T) {
	var xs []int
	ss := sliceStrategy{&factoryRepository{factories: []factory{
		{
			forType: reflect.SliceOf(reflect.TypeOf(0)),
			impl: func(ctx context.Context) (reflect.Value, error) {
				return reflect.ValueOf([]int{10}), nil
			},
		},
	}}}
	done, err := ss.furnishValue(context.Background(), reflect.ValueOf(&xs).Elem())
	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Fail()
	}
	if !reflect.DeepEqual(xs, []int{10}) {
		t.Fail()
	}
}

func TestSliceFactoryMixedFactory(t *testing.T) {
	var xs []int
	ss := sliceStrategy{&factoryRepository{factories: []factory{
		{
			forType: reflect.SliceOf(reflect.TypeOf(0)),
			impl: func(ctx context.Context) (reflect.Value, error) {
				return reflect.ValueOf([]int{10}), nil
			},
		},
		{
			forType: reflect.TypeOf(0),
			impl: func(ctx context.Context) (reflect.Value, error) {
				return reflect.ValueOf(11), nil
			},
		},
	}}}
	done, err := ss.furnishValue(context.Background(), reflect.ValueOf(&xs).Elem())
	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Fail()
	}
	if !reflect.DeepEqual(xs, []int{10, 11}) {
		t.Fail()
	}
}
