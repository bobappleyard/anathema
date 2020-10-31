package binding

import (
	"errors"
	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/router"
	"reflect"
)

var ErrMissingField = errors.New("missing field")
var ErrUnknownEncoding = errors.New("unknown encoding")

type FieldInjector interface {
	InjectField(s *di.Scope, v reflect.Value) error
}

type Injectors struct {
	Route  *router.Route
	Fields []FieldInjector
}

type InjectionSource interface {
	GetInjectors(t reflect.Type) (Injectors, error)
}

type EncodingSource interface {
	GetEncoding(t reflect.Type) (Encoding, error)
}

type Encoding interface {
	Accept(t reflect.Type) bool
	Decode(s string) (reflect.Value, error)
}
