package di

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type testInterface interface {
	method()
}

type testImplementationA struct{}
type testImplementationB struct{}
type testImplementationC struct{}

func (testImplementationA) method() {}
func (testImplementationB) method() {}

func assertEqual(t *testing.T, got, expecting interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expecting) {
		t.Fatalf("got %v, expecting %v", got, expecting)
	}
}

func TestInstanceProvider(t *testing.T) {
	s := &Scope{
		providers: []Provider{
			&instanceProvider{reflect.ValueOf(10)},
		},
	}
	var x int
	s.requireValue(reflect.ValueOf(&x).Elem())

	assertEqual(t, x, 10)
}

func TestFactoryProvider(t *testing.T) {
	s := &Scope{
		providers: []Provider{
			&instanceProvider{reflect.ValueOf(10)},
			&factoryProvider{reflect.ValueOf(func(x int) string {
				return fmt.Sprint(x)
			})},
		},
	}

	var x string
	s.requireValue(reflect.ValueOf(&x).Elem())

	assertEqual(t, x, "10")
}

func TestFactoryConstructor(t *testing.T) {
	f, err := Factory(func(x int) (string, error) {
		return fmt.Sprint(x), nil
	})
	assertEqual(t, err, nil)

	s := &Scope{
		providers: []Provider{
			&instanceProvider{reflect.ValueOf(10)},
			f,
		},
	}

	var x string
	s.requireValue(reflect.ValueOf(&x).Elem())

	assertEqual(t, x, "10")
}

func TestFactoryError(t *testing.T) {
	e := errors.New("test error")
	f, err := Factory(func(x int) (string, error) {
		return fmt.Sprint(x), e
	})
	assertEqual(t, err, nil)

	s := &Scope{
		providers: []Provider{
			&instanceProvider{reflect.ValueOf(10)},
			f,
		},
	}

	var x string
	err = s.requireValue(reflect.ValueOf(&x).Elem())
	assertEqual(t, err, e)
	assertEqual(t, x, "")
}

func TestPointerProvider(t *testing.T) {
	s := &Scope{
		providers: []Provider{
			&pointerProvider{},
			&instanceProvider{reflect.ValueOf(10)},
		},
	}
	var x int
	s.Require(&x)

	assertEqual(t, x, 10)
}

func TestSliceProvider(t *testing.T) {
	s := &Scope{
		providers: []Provider{
			&instanceProvider{reflect.ValueOf(testImplementationA{})},
			&instanceProvider{reflect.ValueOf(testImplementationB{})},
			&pointerProvider{},
			&sliceProvider{},
		},
	}
	var x []testInterface
	s.Require(&x)

	assertEqual(t, x, []testInterface{testImplementationA{}, testImplementationB{}})
}

func TestStructProvider(t *testing.T) {
	s := &Scope{
		providers: []Provider{
			&instanceProvider{reflect.ValueOf(10)},
			&pointerProvider{},
			&sliceProvider{},
			&structProvider{},
		},
	}

	var x struct {
		value  int
		Value  int
		PValue *int
	}
	s.Require(&x)

	assertEqual(t, x.value, 0)
	assertEqual(t, x.Value, 10)
	assertEqual(t, *x.PValue, 10)
}
