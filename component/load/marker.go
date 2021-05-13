package load

import (
	"github.com/bobappleyard/anathema/component"
	"github.com/bobappleyard/anathema/component/di"
	"reflect"
)

type markerRule struct {
}

type markerInjector struct {
	v reflect.Value
}

func (m markerInjector) WillCreate() reflect.Type {
	return m.v.Type()
}

func (m markerInjector) Create(p di.Process) (reflect.Value, error) {
	return m.v, nil
}

var markerType = reflect.TypeOf(new(component.Marker)).Elem()

func (r *markerRule) Apply(b di.Builder, t reflect.Type) {
	if t.Kind() != reflect.Interface {
		return
	}
	if !t.AssignableTo(markerType) {
		return
	}
	n := t.NumMethod()
	for i := 0; i < n; i++ {
		m := t.Method(i)
		if m.PkgPath == "" {
			return
		}
	}
	b.Constructor(&markerInjector{reflect.New(t).Elem()})
}

