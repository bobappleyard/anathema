package di

import (
	"context"
	"errors"
	"reflect"
)

var ErrInjectionFailed = errors.New("injection failed")
var ErrInvalidProvider = errors.New("invalid provider")

type Builder interface {
	Cache(scope string)
	Constructor(c Constructor)
	FallbackConstructor(c Constructor)
	Mutator(m Mutator)
	Complete()
}

type Process interface {
	Furnish(p interface{}) error
	RequireValue(t reflect.Type) (reflect.Value, error)
	Apply(b Builder, t reflect.Type)
}

type Rule interface {
	Apply(b Builder, t reflect.Type)
}

type Constructor interface {
	WillCreate() reflect.Type
	Create(p Process) (reflect.Value, error)
}

type Mutator interface {
	Update(p Process, v reflect.Value) error
}

var scopeKey = &struct{}{}

func Furnish(ctx context.Context, ref interface{}) error {
	s := GetScope(ctx)
	return s.Furnish(ref)
}

func GetScope(ctx context.Context) *Scope {
	scope := ctx.Value(scopeKey)
	if scope == nil {
		return nil
	}
	return scope.(*Scope)
}

func EnterScope(ctx context.Context, name string) context.Context {
	next := GetScope(ctx)
	if next == nil {
		next = baseScope
	}
	return context.WithValue(ctx, scopeKey, &Scope{next: next, name: name})
}
