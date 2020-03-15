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
)

// Require calls the function f, using extract to furnish its inputs. If any of
// the factories that are invoked in the course of retrieving these values
// fails, the error is returned by Require. Additionally, f may be declared to
// return an error, which is duly returned by Require.
func Require(ctx context.Context, f interface{}) error {
	return GetRegistry(ctx).Require(ctx, f)
}

// ProvideValue registers a value which, when requested using Require, will be
// passed into the requiring function.
func ProvideValue(ctx context.Context, x interface{}) context.Context {
	r, ctx := New(ctx)
	r.ProvideValue(x)
	return ctx
}

// Provide registers a function with the purpose of constructing values of a
// particular type. This is given by the return type of said function (which may
// optionally also return an error type). Any declared inputs to this function
// are dependencies that are resolved using Require.
func Provide(ctx context.Context, f interface{}) context.Context {
	r, ctx := New(ctx)
	r.Provide(f)
	return ctx
}

// New creates a new Registry and binds it to the provided context.
func New(ctx context.Context) (*Registry, context.Context) {
	r := GetRegistry(ctx).Extend()
	return r, r.Bind(ctx)
}
