package di

import (
	"context"
	"reflect"
	"testing"
)

func TestPtrStrategy(t *testing.T) {
	var x *int
	s := ptrStrategy{}
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
	if *x != 10 {
		t.Fail()
	}
}
