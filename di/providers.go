package di

import (
	"reflect"
)

func Instance(x interface{}) Provider {
	return &instanceProvider{reflect.ValueOf(x)}
}

func Factory(f interface{}) (Provider, error) {
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
	return &factoryProvider{fv}, nil
}

var errorType = reflect.TypeOf(new(error)).Elem()

type instanceProvider struct {
	v reflect.Value
}

func (p *instanceProvider) Apply(t reflect.Type, found []Provider) []Provider {
	if p.v.Type().AssignableTo(t) {
		return append(found, p)
	}
	return found
}

func (p *instanceProvider) Provide(s *Scope, v reflect.Value) error {
	v.Set(p.v)
	return nil
}

type factoryProvider struct {
	f reflect.Value
}

func (p *factoryProvider) Apply(t reflect.Type, found []Provider) []Provider {
	if p.f.Type().Out(0).AssignableTo(t) {
		return append(found, p)
	}
	return found
}

func (p *factoryProvider) Provide(s *Scope, v reflect.Value) error {
	t := p.f.Type()
	in := make([]reflect.Value, t.NumIn())
	for i := range in {
		v := reflect.New(t.In(i))
		if err := s.requireValue(v.Elem()); err != nil {
			return err
		}
		in[i] = v.Elem()
	}
	out := p.f.Call(in)
	if len(out) == 2 && !out[1].IsNil() {
		return out[1].Interface().(error)
	}
	v.Set(out[0])
	return nil
}

func genericProvider(k reflect.Kind, p Provider, t reflect.Type, found []Provider) []Provider {
	if t.Kind() != k {
		return found
	}
	if len(found) != 0 {
		return found
	}
	return append(found, p)
}

type pointerProvider struct {
}

func (p *pointerProvider) Apply(t reflect.Type, found []Provider) []Provider {
	return genericProvider(reflect.Ptr, p, t, found)
}

func (p *pointerProvider) Provide(s *Scope, v reflect.Value) error {
	if v.CanAddr() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return s.requireValue(v.Elem())
}

type sliceProvider struct {
}

func (p *sliceProvider) Apply(t reflect.Type, found []Provider) []Provider {
	return genericProvider(reflect.Slice, p, t, found)
}

func (p *sliceProvider) Provide(s *Scope, v reflect.Value) error {
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

type structProvider struct {
}

func (p *structProvider) Apply(t reflect.Type, found []Provider) []Provider {
	return genericProvider(reflect.Struct, p, t, found)
}

func (p *structProvider) Provide(s *Scope, v reflect.Value) error {
	t := v.Type()
	value := reflect.New(t).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		if err := s.requireValue(value.FieldByIndex(f.Index)); err != nil {
			return err
		}
	}

	v.Set(value)
	return nil
}
