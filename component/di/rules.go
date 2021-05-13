package di

import (
	"reflect"
)

// Instance introduces a dependency that has already been created outside the DI system.
func Instance(x interface{}) Rule {
	return &instanceRule{reflect.ValueOf(x)}
}

type instanceRule struct {
	v reflect.Value
}

func (r *instanceRule) Apply(b Builder, t reflect.Type) {
	if !r.v.Type().AssignableTo(t) {
		return
	}
	b.Constructor(r)
	b.Complete()
}

func (r *instanceRule) WillCreate() reflect.Type {
	return r.v.Type()
}

func (r *instanceRule) Create(p Process) (reflect.Value, error) {
	return r.v, nil
}

// Factory introduces a dependency that is furnished by calling a factory function. This function can take any number of
// inputs and one or two outputs. The first output is the dependency itself, the second output, if present, should be an
// error. The inputs will be injected as part of the factory being called.
func Factory(f interface{}) (Rule, error) {
	fv := reflect.ValueOf(f)
	ft := fv.Type()
	if fv.Kind() != reflect.Func {
		return nil, ErrInvalidProvider
	}
	switch ft.NumOut() {
	case 1:
	case 2:
		if !ft.Out(1).AssignableTo(errorType) {
			return nil, ErrInvalidProvider
		}
	default:
		return nil, ErrInvalidProvider
	}
	return &factoryRule{fv}, nil
}

var errorType = reflect.TypeOf(new(error)).Elem()

type factoryRule struct {
	f reflect.Value
}

func (r *factoryRule) Apply(b Builder, t reflect.Type) {
	if !r.f.Type().Out(0).AssignableTo(t) {
		return
	}
	b.Constructor(r)
	b.Complete()
}

func (r *factoryRule) WillCreate() reflect.Type {
	return r.f.Type().Out(0)
}

func (r *factoryRule) Create(p Process) (reflect.Value, error) {
	out, err := injectedCall(p, r.f)
	if err != nil {
		return reflect.Value{}, err
	}
	if len(out) == 2 && !out[1].IsNil() {
		return reflect.Value{}, out[1].Interface().(error)
	}
	return out[0], nil
}

// baseScope installs the built-in rules
var baseScope = func() *Scope {
	res := new(Scope)
	res.AddRule(&structRule{})
	res.AddRule(&sliceRule{})
	res.AddRule(&ptrRule{})
	res.AddRule(&injectedRule{})
	return res
}()

// injectedRule applies to any dependency that provides an Inject method. This method takes any number of inputs and
// zero or one output. The inputs to the method are injected into a call following the creation of the dependency. The
// output, if present, should be an error.
type injectedRule struct{}

func (r *injectedRule) Apply(b Builder, t reflect.Type) {
	_, ok := t.MethodByName("Inject")
	if !ok {
		return
	}
	b.Mutator(r)
}

func (m *injectedRule) Update(p Process, v reflect.Value) error {
	method := v.MethodByName("Inject")
	out, err := injectedCall(p, method)
	if err != nil {
		return err
	}
	if len(out) == 1 && !out[0].IsNil() {
		return out[0].Interface().(error)
	}
	return nil
}

// ptrRule applies to dependencies that are pointer types.
type ptrRule struct{}

type ptrInjector struct {
	t reflect.Type
}

func (r *ptrRule) Apply(b Builder, t reflect.Type) {
	if t.Kind() != reflect.Ptr {
		return
	}
	b.FallbackConstructor(&ptrInjector{t})
}

func (r *ptrInjector) WillCreate() reflect.Type {
	return r.t
}

func (r *ptrInjector) Create(p Process) (reflect.Value, error) {
	pv := reflect.New(r.t.Elem())
	v, err := p.RequireValue(r.t.Elem())
	if err != nil {
		return reflect.Value{}, err
	}
	pv.Elem().Set(v)
	return pv, nil
}

// Slice rules applies to dependencies that are slices. The injection process will find all dependencies that match the
// slice's element type, the resulting slice contains those dependencies. This is useful for creating pluggable systems.
type sliceRule struct{}

type sliceInjector struct {
	t reflect.Type
}

func (r *sliceRule) Apply(b Builder, t reflect.Type) {
	if t.Kind() != reflect.Slice {
		return
	}
	injector := &sliceInjector{t}
	b.FallbackConstructor(injector)
}

func (r *sliceInjector) WillCreate() reflect.Type {
	return r.t
}

func (r *sliceInjector) Create(p Process) (reflect.Value, error) {
	pv := reflect.New(r.t).Elem()
	pv.Set(reflect.MakeSlice(r.t, 0, 0))
	var builder sliceBuilder
	p.Apply(&builder, r.t.Elem())
	for _, t := range builder.types {
		ev, err := p.RequireValue(t)
		if err != nil {
			return reflect.Value{}, err
		}
		pv.Set(reflect.Append(pv, ev))
	}
	return pv, nil
}

type sliceBuilder struct {
	types []reflect.Type
}

func (s *sliceBuilder) Constructor(c Constructor) {
	s.types = append(s.types, c.WillCreate())
}

// These are intentionally left blank to ignore these types

func (s *sliceBuilder) Cache(scope string) {}

func (s *sliceBuilder) FallbackConstructor(c Constructor) {}

func (s *sliceBuilder) Mutator(m Mutator) {}

func (s *sliceBuilder) Complete() {}

// structRule applies to dependencies that are structs. This will inject all exported fields on the struct that lack
// struct tags. You can therefore add your own rules for tagged fields.
type structRule struct{}

type structInjector struct {
	t reflect.Type
}

func (r *structRule) Apply(b Builder, t reflect.Type) {
	if t.Kind() != reflect.Struct {
		return
	}
	injector := &structInjector{t}
	b.FallbackConstructor(injector)
	b.Mutator(injector)
}

func (r *structInjector) WillCreate() reflect.Type {
	return r.t
}

func (r *structInjector) Create(p Process) (reflect.Value, error) {
	return reflect.New(r.t).Elem(), nil
}

func (r *structInjector) Update(p Process, v reflect.Value) error {
	t := v.Type()
	n := t.NumField()

	for i := 0; i < n; i++ {
		f := t.Field(i)

		if f.PkgPath != "" || f.Tag != "" {
			continue
		}

		if err := r.injectField(p, v.FieldByIndex(f.Index), f.Type); err != nil {
			return err
		}
	}

	return nil
}

func (r *structInjector) injectField(p Process, v reflect.Value, t reflect.Type) error {
	if t.Kind() == reflect.Struct {
		return r.Update(p, v)
	}
	fv, err := p.RequireValue(t)
	if err != nil {
		return err
	}
	v.Set(fv)
	return nil
}
