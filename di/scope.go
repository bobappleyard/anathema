package di

import (
	"errors"
	"reflect"
)

var ErrInjectionFailed = errors.New("injection failed")
var ErrInvalidProvider = errors.New("invalid provider")

type Scope struct {
	next      *Scope
	providers []Provider
}

type Provider interface {
	Apply(t reflect.Type, found []Provider) []Provider
	Provide(s *Scope, v reflect.Value) error
}

func (s *Scope) AddProvider(p Provider) {
	s.providers = append(s.providers, p)
}

func (s *Scope) Require(ptr interface{}) error {
	return s.requireValue(reflect.ValueOf(ptr))
}

func (s *Scope) requireValue(v reflect.Value) error {
	t := v.Type()
	providers := s.findProviders(t)

	if len(providers) != 1 {
		return ErrInjectionFailed
	}

	return providers[0].Provide(s, v)
}

func (s *Scope) findProviders(t reflect.Type) []Provider {
	var providers []Provider
	for cur := s; cur != nil; cur = cur.next {
		for _, provider := range cur.providers {
			providers = provider.Apply(t, providers)
		}
	}
	return providers
}
