package di

import (
	"context"
	"reflect"
)

// A Factory creates values to inject into functions that require them.
type Factory func(context.Context) (interface{}, error)

// A Registry maps types to factories
type Registry struct {
	next    *Registry
	entries map[reflect.Type]Factory
}

var (
	contextType = reflect.TypeOf(new(context.Context)).Elem()
	marker      = new(struct{})
)

// GetRegistry searches the provided context for a Registry and returns it, or
// nil otherwise.
func GetRegistry(ctx context.Context) *Registry {
	if r := ctx.Value(marker); r != nil {
		return r.(*Registry)
	}
	return nil
}

// Extend creates a new Registry that will fall back on the current reciever if
// it cannot locate a required type.
func (r *Registry) Extend() *Registry {
	return &Registry{r, map[reflect.Type]Factory{}}
}

// Bind attaches the current receiver to the provided context, returning the new
// context formed from that operation.
func (r *Registry) Bind(ctx context.Context) context.Context {
	return context.WithValue(ctx, marker, r)
}

// Apply uses the inputs to a function in order to determine what values to
// extract. The outputs of the function are returned.
func (r *Registry) Apply(ctx context.Context, f interface{}) ([]reflect.Value, error) {
	fv := reflect.ValueOf(f)
	ft := fv.Type()

	in := make([]reflect.Value, ft.NumIn())

	for i := 0; i < ft.NumIn(); i++ {
		x, err := Extract(ctx, ft.In(i))
		if err != nil {
			return nil, err
		}
		in[i] = reflect.ValueOf(x)
	}

	out := fv.Call(in)
	return out, nil
}

// Require applies a function before making some assertions about what that
// function returns.
func (r *Registry) Require(ctx context.Context, f interface{}) error {
	vs, err := r.Apply(ctx, f)
	if err != nil {
		return err
	}
	if len(vs) == 0 {
		return nil
	}
	if vs[0].IsNil() {
		return nil
	}
	return vs[0].Interface().(error)
}

// AddFactory registers a function as a factory for the given type.
func (r *Registry) AddFactory(t reflect.Type, f Factory) {
	r.entries[t] = f
}

// Insert registers a value for the given type.
func (r *Registry) Insert(t reflect.Type, v reflect.Value) {
	vi := v.Interface()
	r.AddFactory(t, func(context.Context) (interface{}, error) {
		return vi, nil
	})
}

// ProvideValue registers a value, inferring its type.
func (r *Registry) ProvideValue(x interface{}) {
	v := reflect.ValueOf(x)
	r.Insert(v.Type(), v)
}

// Provide registers a factory function that may itself require values of
// particular types.
func (r *Registry) Provide(f interface{}) {
	ft := reflect.TypeOf(f)
	r.AddFactory(ft.Out(0), func(ctx context.Context) (interface{}, error) {
		vs, err := r.Apply(ctx, f)
		if err != nil {
			return nil, err
		}
		if len(vs) == 2 && !vs[1].IsNil() {
			return nil, vs[1].Interface().(error)
		}
		return vs[0].Interface(), nil
	})
}
