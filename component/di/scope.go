package di

import (
	"io"
	"reflect"
	"strings"
)

type Scope struct {
	next  *Scope
	rules []Rule
	name  string
	cache cache
}

func (s *Scope) AddRule(p Rule) {
	s.rules = append(s.rules, p)
}

func (s *Scope) Furnish(ptr interface{}) error {
	return furnish(s, ptr)
}

func (s *Scope) RequireValue(t reflect.Type) (reflect.Value, error) {
	var b builder
	s.Apply(&b, t)

	cs := s.cacheScope(b.cache)

	if cs != nil {
		v, ok := cs.cache.retrieve(t)
		if ok {
			return v, nil
		}
	}

	process := &process{
		scope:   s,
		created: map[reflect.Type]reflect.Value{},
	}
	v, err := process.run(b, t)
	if err != nil {
		return reflect.Value{}, err
	}

	if cs != nil {
		v = cs.cache.store(t, v)
	}

	return v, nil
}

func (s *Scope) Apply(b Builder, t reflect.Type) {
	if s == nil {
		return
	}
	s.next.Apply(b, t)
	for _, rule := range s.rules {
		rule.Apply(b, t)
	}
}

type multiError struct {
	errs []error
}

func (m *multiError) Error() string {
	var sb strings.Builder
	sb.WriteString(m.errs[0].Error())
	for _, e := range m.errs[1:] {
		sb.WriteString("; ")
		sb.WriteString(e.Error())
	}
	return sb.String()
}

func (s *Scope) Close() error {
	var errs []error
	for _, v := range s.cache.values() {
		if v, ok := v.Interface().(io.Closer); ok {
			if err := v.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return &multiError{errs}
}

func (s *Scope) cacheScope(names []string) *Scope {
	if s == nil {
		return nil
	}
	cs := s.next.cacheScope(names)
	if cs != nil {
		return cs
	}
	for _, n := range names {
		if s.name == n {
			return s
		}
	}
	return nil
}
