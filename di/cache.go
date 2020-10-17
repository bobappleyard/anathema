package di

import (
	"reflect"
	"sync"
)

type cache struct {
	lock   sync.RWMutex
	values map[reflect.Type][]Furnisher
}

func (c *cache) get(t reflect.Type) ([]Furnisher, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	fs, ok := c.values[t]
	return fs, ok
}

func (c *cache) put(t reflect.Type, fs []Furnisher) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.values == nil {
		c.values = make(map[reflect.Type][]Furnisher)
	}
	c.values[t] = fs
}
