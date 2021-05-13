package load

import (
	"github.com/bobappleyard/anathema/component/di"
	"reflect"
	"sync"
)

var serviceLock sync.RWMutex

type serviceRule struct {
	service reflect.Type
}

func (r *serviceRule) Apply(b di.Builder, t reflect.Type) {
	if t == serviceType {
		return
	}
	if !r.service.AssignableTo(t) {
		return
	}
	b.Constructor(r)
}

func (r *serviceRule) WillCreate() reflect.Type {
	return r.service
}

func (r *serviceRule) Create(p di.Process) (reflect.Value, error) {
	v, err := p.RequireValue(r.service.Elem())
	if err != nil {
		return reflect.Value{}, err
	}
	return v.Addr(), nil
}
