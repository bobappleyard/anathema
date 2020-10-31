package di

import (
	"errors"
	"reflect"
)

var ErrInjectionFailed = errors.New("injection failed")
var ErrInvalidProvider = errors.New("invalid provider")

type Scope struct {
	next  *Scope
	rules []Rule
}

type Rule interface {
	Apply(t reflect.Type, found []Provider) []Provider
}

type Provider interface {
	Provide(s *Scope, v reflect.Value) error
}

func NewScope(next *Scope) *Scope {
	return &Scope{next: next}
}

func (s *Scope) AddRule(p Rule) {
	s.rules = append(s.rules, p)
}

func (s *Scope) Require(ptr interface{}) error {
	return s.RequireValue(reflect.ValueOf(ptr))
}

func (s *Scope) RequireValue(v reflect.Value) error {
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
		// Iterate backwards so that the property of later additions overriding earlier ones is consistent across the
		// piece.
		for i := len(cur.rules) - 1; i >= 0; i-- {
			rule := cur.rules[i]
			providers = rule.Apply(t, providers)
		}
	}
	return providers
}
