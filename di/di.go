// Package di provides a simple dependency injection container
//
// This is a simple dependency injection container implemented in terms of
// context.Context.
//
// There are two flavours of function. Most uses should be in term of Require
// and Provide. These afford a simple, declarative API in terms of funcs. The
// low-level functions Insert and Extract operate in terms of the reflect
// package. This gives you more control, but with more verbosity and less type
// safety.
package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// Errors
var (
	ErrNoFactory = errors.New("no registered factory")
)

type containerKey struct {
	typ reflect.Type
}

type factoryFunc func(context.Context) (interface{}, error)

// Insert extends a context with a value.
//
// This can be considered the low-level version of Provide.
func Insert(ctx context.Context, typ reflect.Type, value reflect.Value) context.Context {
	vi := value.Interface()
	f := func(context.Context) (interface{}, error) {
		return vi, nil
	}
	return context.WithValue(ctx, containerKey{typ}, factoryFunc(f))
}

// Extract queries the context for a value of the required type within a DI
// container. This may cause the execution of an arbitrarily long chain of
// factories. If this is not possible for whatever reason, said reason is
// returned as an error.
//
// This can be considered the low-level version of Require.
func Extract(ctx context.Context, typ reflect.Type) (interface{}, error) {
	f, ok := ctx.Value(containerKey{typ}).(factoryFunc)
	if !ok {
		return nil, fmt.Errorf("%s: %w", typ, ErrNoFactory)
	}
	return f(ctx)
}

// Require calls the function f, using extract to furnish its inputs. If any of
// the factories that are invoked in the course of retrieving these values
// fails, the error is returned by Require. Additionally, f may be declared to
// return an error, which is duly returned by Require.
func Require(ctx context.Context, f interface{}) error {
	fv := reflect.ValueOf(f)
	a := parseArgs(fv.Type(), false)
	in, err := a.assembleInputs(ctx)
	if err != nil {
		return err
	}
	out := fv.Call(in)
	return a.extractError(out, 0)
}

// Provide registers a function with the purpose of constructing values of a
// particular type. This is given by the return type of said function (which may
// optionally also return an error type). Any declared inputs to this function
// are dependencies that are resolved using Require.
func Provide(ctx context.Context, f interface{}) context.Context {
	fv := reflect.ValueOf(f)
	args := parseArgs(fv.Type(), true)
	factory := buildFactory(fv, args)
	return context.WithValue(ctx, containerKey{args.out}, factory)
}

func buildFactory(fv reflect.Value, a args) factoryFunc {
	return func(ctx context.Context) (interface{}, error) {
		in, err := a.assembleInputs(ctx)
		if err != nil {
			return nil, err
		}
		out := fv.Call(in)
		ret := out[0].Interface()
		return ret, a.extractError(out, 1)
	}
}

type args struct {
	in     []reflect.Type
	out    reflect.Type
	hasErr bool
}

var errorType = reflect.TypeOf(new(error)).Elem()

func parseArgs(ft reflect.Type, hasRet bool) args {
	if ft.Kind() != reflect.Func {
		panic("expecting a function")
	}
	var res args
	res.in = make([]reflect.Type, ft.NumIn())
	for i := range res.in {
		res.in[i] = ft.In(i)
	}
	n := ft.NumOut()
	if hasRet {
		assertOutput(n, 1, 2)
		res.out = ft.Out(0)
		res.hasErr = n == 2
		return res
	}
	assertOutput(n, 0, 1)
	res.hasErr = n == 1
	return res
}

func (a args) extractError(out []reflect.Value, at int) error {
	if !a.hasErr {
		return nil
	}
	e := out[at].Interface()
	if e == nil {
		return nil
	}
	return e.(error)
}

func (a args) assembleInputs(ctx context.Context) ([]reflect.Value, error) {
	in := make([]reflect.Value, len(a.in))
	for i, t := range a.in {
		x, err := Extract(ctx, t)
		if err != nil {
			return nil, err
		}
		in[i] = reflect.ValueOf(x)
	}
	return in, nil
}

func assertOutput(n, min, max int) {
	if n < min || n > max {
		panic("wrong number of outputs")
	}
}
