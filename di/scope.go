package di

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrInjectionFailed = errors.New("injection failed")

type Scope struct {
	next  *Scope
	rules []Rule
	cache cache
}

type Rule interface {
	Apply(s *Scope, t reflect.Type, results []Furnisher) []Furnisher
}

type Furnisher interface {
	Furnish(s *Scope, p reflect.Value) error
}

func (s *Scope) Furnish(ptr reflect.Value) (err error) {
	t := ptr.Type().Elem()
	defer addErrorScope(&err, "injecting %v", t)

	furnishers := s.findFurnishers(t)
	if len(furnishers) != 1 {
		return ErrInjectionFailed
	}

	if err := furnishers[0].Furnish(s, ptr); err != nil {
		return err
	}
	s.cacheInjection(t, ptr.Elem())

	return nil
}

func (s *Scope) findFurnishers(t reflect.Type) []Furnisher {
	if fs, ok := s.cache.get(t); ok {
		return fs
	}
	var res []Furnisher
	for cur := s; cur != nil; cur = cur.next {
		for _, r := range cur.rules {
			res = r.Apply(s, t, res)
		}
	}
	return res
}

func (s *Scope) cacheInjection(t reflect.Type, value reflect.Value) {
	s.cache.put(t, []Furnisher{&instanceFurnisher{value}})
}

func addErrorScope(errp *error, fmtStr string, args ...interface{}) {
	err := *errp
	if err == nil {
		return
	}
	args = append(args, err)
	*errp = fmt.Errorf(fmtStr+": %w", args...)
}
