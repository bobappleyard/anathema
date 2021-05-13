package di

import (
	"reflect"
	"sync"
)

type cache struct {
	lock  sync.RWMutex
	items map[reflect.Type]reflect.Value
}

func (c *cache) retrieve(key reflect.Type) (reflect.Value, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	v, ok := c.items[key]
	return v, ok
}

func (c *cache) store(key reflect.Type, b reflect.Value) reflect.Value {
	c.lock.Lock()
	defer c.lock.Unlock()

	if v, ok := c.items[key]; ok {
		return v
	}

	c.items[key] = b
	return b
}

func (c *cache) values() []reflect.Value {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var res []reflect.Value
	for _, v := range c.items {
		res = append(res, v)
	}

	return res
}