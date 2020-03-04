package binding

import (
	"errors"
	"reflect"
	"strconv"
	"unsafe"
)

type mechanism interface {
	read(unsafe.Pointer) (string, error)
	write(unsafe.Pointer, string) error
	defined() bool
}

var (
	intType    = reflect.TypeOf(0)
	stringType = reflect.TypeOf("")
)

func mechanismFor(t reflect.Type) mechanism {
	switch t {
	case intType:
		return intValue{}
	case stringType:
		return stringValue{}
	}
	return unsupported{}
}

type unsupported struct{}

var ErrUnsupportedType = errors.New("unsupported type")

func (unsupported) read(unsafe.Pointer) (string, error) {
	return "", ErrUnsupportedType
}

func (unsupported) write(unsafe.Pointer, string) error {
	return ErrUnsupportedType
}

func (unsupported) defined() bool {
	return true
}

type undefined struct{}

func (undefined) read(unsafe.Pointer) (string, error) {
	return "", nil
}

func (undefined) write(unsafe.Pointer, string) error {
	return nil
}

func (undefined) defined() bool {
	return false
}

type stringValue struct{}

func (stringValue) read(p unsafe.Pointer) (string, error) {
	return *(*string)(p), nil
}

func (stringValue) write(p unsafe.Pointer, s string) error {
	*(*string)(p) = s
	return nil
}

func (stringValue) defined() bool {
	return true
}

type intValue struct{}

func (intValue) read(p unsafe.Pointer) (string, error) {
	return strconv.Itoa(*(*int)(p)), nil
}

func (intValue) write(p unsafe.Pointer, s string) error {
	x, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*(*int)(p) = x
	return nil
}

func (intValue) defined() bool {
	return true
}
