package di

import (
	"reflect"
)

func Instance(x interface{}) Rule {
	return &instanceRule{reflect.ValueOf(x)}
}

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

type instanceRule struct {
	v reflect.Value
}

func (r *instanceRule) Apply(t reflect.Type, found []Provider) []Provider {
	if r.v.Type().AssignableTo(t) {
		return append(found, r)
	}
	return found
}

func (r *instanceRule) Provide(s *Scope, v reflect.Value) error {
	v.Set(r.v)
	return nil
}

type factoryRule struct {
	f reflect.Value
}

func (r *factoryRule) Apply(t reflect.Type, found []Provider) []Provider {
	if r.f.Type().Out(0).AssignableTo(t) {
		return append(found, r)
	}
	return found
}

func (r *factoryRule) Provide(s *Scope, v reflect.Value) error {
	out, err := injectedCall(s, r.f)
	if err != nil {
		return err
	}
	if len(out) == 2 && !out[1].IsNil() {
		return out[1].Interface().(error)
	}
	v.Set(out[0])
	return nil
}

func injectedCall(s *Scope, f reflect.Value) ([]reflect.Value, error) {
	t := f.Type()
	in := make([]reflect.Value, t.NumIn())
	for i := range in {
		v := reflect.New(t.In(i))
		if err := s.RequireValue(v.Elem()); err != nil {
			return nil, err
		}
		in[i] = v.Elem()
	}
	return f.Call(in), nil
}

func genericRule(k reflect.Kind, p Provider, t reflect.Type, found []Provider) []Provider {
	if t.Kind() != k {
		return found
	}
	if len(found) != 0 {
		return found
	}
	return append(found, p)
}

type pointerRule struct {
}

func (r *pointerRule) Apply(t reflect.Type, found []Provider) []Provider {
	return genericRule(reflect.Ptr, r, t, found)
}

func (r *pointerRule) Provide(s *Scope, v reflect.Value) error {
	if v.CanAddr() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return s.RequireValue(v.Elem())
}

type sliceRule struct {
}

func (r *sliceRule) Apply(t reflect.Type, found []Provider) []Provider {
	return genericRule(reflect.Slice, r, t, found)
}

func (r *sliceRule) Provide(s *Scope, v reflect.Value) error {
	t := v.Type()
	providers := s.findProviders(t.Elem())
	slice := reflect.MakeSlice(t, len(providers), len(providers))
	for i, p := range providers {
		if err := p.Provide(s, slice.Index(i)); err != nil {
			return err
		}
	}
	v.Set(slice)
	return nil
}

type structRule struct {
}

func (r *structRule) Apply(t reflect.Type, found []Provider) []Provider {
	return genericRule(reflect.Struct, r, t, found)
}

func (r *structRule) Provide(s *Scope, v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		if err := s.RequireValue(v.FieldByIndex(f.Index)); err != nil {
			return err
		}
	}

	return nil
}

type injectRule struct {
}

type injectProvider struct {
	fallback []Provider
}

func (r *injectRule) Apply(t reflect.Type, found []Provider) []Provider {
	if _, ok := t.MethodByName("Inject"); ok {
		return []Provider{&injectProvider{found}}
	}
	return found
}

func (r *injectProvider) Provide(s *Scope, v reflect.Value) error {
	if len(r.fallback) != 1 {
		return ErrInjectionFailed
	}
	if err := r.fallback[0].Provide(s, v); err != nil {
		return err
	}
	out, err := injectedCall(s, v.MethodByName("Inject"))
	if err != nil {
		return err
	}
	if len(out) == 1 && !out[0].IsNil() {
		return out[0].Interface().(error)
	}
	return nil
}
