package fields

import (
	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/server/binding"
	"reflect"
	"sync"
)

type CacheService struct {
	a.Service

	lock  sync.RWMutex
	items map[reflect.Type][]binding.FieldInjector
}

func (c *CacheService) Get(t reflect.Type) []binding.FieldInjector {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.items[t]
}

func (c *CacheService) Set(t reflect.Type, injectors []binding.FieldInjector) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.items == nil {
		c.items = map[reflect.Type][]binding.FieldInjector{}
	}
	c.items[t] = injectors
}
