package di

import (
	"context"
	"reflect"
	"testing"
)

func TestContextStrategy(t *testing.T) {
	var x context.Context
	s := contextStrategy{}
	ctx := context.Background()
	done, err := s.furnishValue(ctx, reflect.ValueOf(&x).Elem())
	if !done {
		t.Fail()
	}
	if err != nil {
		t.Error(err)
	}
	if x != ctx {
		t.Fail()
	}
}
