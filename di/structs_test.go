package di

import (
	"context"
	"reflect"
	"testing"
)

func TestStructStrategy(t *testing.T) {
	var x struct {
		X int
		x int
		Y int `tag:"value"`
	}
	s := structStrategy{}
	ctx := toContext(context.Background(), &furnisher{
		[]strategy{&testStrategy{reflect.ValueOf(10)}},
	})
	done, err := s.furnishValue(ctx, reflect.ValueOf(&x).Elem())
	if !done {
		t.Fail()
	}
	if err != nil {
		t.Error(err)
	}
	if x.X != 10 {
		t.Fail()
	}
}
