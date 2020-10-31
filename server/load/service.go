package load

import (
	"github.com/bobappleyard/anathema/di"
	"reflect"
	"sync"
)

var serviceLock sync.RWMutex

type serviceRule struct {
	service reflect.Type
	cache   reflect.Value
}

func (r *serviceRule) Apply(t reflect.Type, found []di.Provider) []di.Provider {
	if r.service != t && r.service.AssignableTo(t) {
		return append(found, r)
	}
	return found
}

func (r *serviceRule) Provide(s *di.Scope, v reflect.Value) error {
	service := r.getCachedService()
	if !service.IsValid() {
		service = reflect.New(r.service).Elem()
		if err := s.RequireValue(service); err != nil {
			return err
		}
		service = r.setCachedService(service)
	}
	v.Set(service)
	return nil
}

func (r *serviceRule) getCachedService() reflect.Value {
	serviceLock.RLock()
	defer serviceLock.RUnlock()

	return r.cache
}

func (r *serviceRule) setCachedService(service reflect.Value) reflect.Value {
	serviceLock.Lock()
	defer serviceLock.Unlock()

	if r.cache.IsValid() {
		return r.cache
	}
	r.cache = service
	return service
}
