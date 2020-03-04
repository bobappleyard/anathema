package binding

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestMechanism(t *testing.T) {
	for _, test := range []struct {
		name string
		str  string
		m    mechanism
		x    interface{}
	}{
		{
			name: "string",
			str:  "hello",
			m:    stringValue{},
			x:    "hello",
		},
		{
			name: "int",
			str:  "1234",
			m:    intValue{},
			x:    1234,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			targ := reflect.New(reflect.TypeOf(test.x)).Elem()
			err := test.m.write(unsafe.Pointer(targ.UnsafeAddr()), test.str)
			if err != nil {
				t.Fatalf("unexpected error %q", err)
			}
			x := targ.Interface()
			if x != test.x {
				t.Errorf("got %v, expecting %v", x, test.x)
			}
		})
	}
}
