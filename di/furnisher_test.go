package di

import (
	"context"
	"reflect"
	"testing"
)

func TestToContext(t *testing.T) {
	f := &furnisher{}
	ctx := toContext(context.Background(), f)
	if fromContext(ctx) != f {
		t.Fail()
	}
}

type testStrategy struct {
	v reflect.Value
}

func (s *testStrategy) furnishValue(ctx context.Context, v reflect.Value) (bool, error) {
	v.Set(s.v)
	return true, nil
}

func TestFurnisher_FurnishValue(t *testing.T) {
	f := &furnisher{
		strategies: []strategy{&testStrategy{reflect.ValueOf(10)}},
	}
	var x int
	err := f.furnishValue(context.Background(), reflect.ValueOf(&x).Elem())
	if err != nil {
		t.Error(err)
	}
	if x != 10 {
		t.Fail()
	}
}

func TestFurnisher_FurnishValueUnsupported(t *testing.T) {
	f := &furnisher{}
	var x int
	err := f.furnishValue(context.Background(), reflect.ValueOf(&x).Elem())
	if err != ErrUnsupported {
		t.Fail()
	}
}

func TestFurnish(t *testing.T) {
	f := &furnisher{
		strategies: []strategy{&testStrategy{reflect.ValueOf(10)}},
	}
	var x int
	ctx := toContext(context.Background(), f)
	err := Furnish(ctx, &x)
	if err != nil {
		t.Error(err)
	}
	if x != 10 {
		t.Fail()
	}
}

func TestFurnishNoFurnisher(t *testing.T) {
	var x int
	err := Furnish(context.Background(), &x)
	if err != ErrNoFurnisher {
		t.Fail()
	}
}

func TestFurnishIntegration(t *testing.T) {
	var x struct {
		Xs []int
	}
	ctx := toContext(context.Background(), newFurnisher([]factory{
		{
			forType: reflect.TypeOf(0),
			impl: func(ctx context.Context) (reflect.Value, error) {
				return reflect.ValueOf(10), nil
			},
		},
	}))
	err := Furnish(ctx, &x)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x.Xs, []int{10}) {
		t.Fail()
	}
}
