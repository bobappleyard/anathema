package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/typereg"
)

var ErrInjectionFailed = errors.New("injection failed")

func Furnish(ctx context.Context, ptr interface{}) error {
	pv := reflect.ValueOf(ptr)
	if pv.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: injection target is not a pointer", ErrInjectionFailed)
	}
	pv = pv.Elem()
	return getScope(ctx).furnish(pv)
}

type Option func(*scope) error

func Enter(ctx context.Context, options ...Option) (context.Context, error) {
	s := &scope{
		next: getScope(ctx),
	}
	for _, o := range options {
		if err := o(s); err != nil {
			return nil, err
		}
	}
	return context.WithValue(ctx, scopeKey, s), nil
}

func Scan(scan, name string) Option {
	return func(s *scope) error {
		for _, p := range scanProviders(scan, name) {
			rule, err := newProviderRule(p)
			if err != nil {
				return err
			}
			s.rules = append(s.rules, rule)
		}
		return nil
	}
}

var providerType = reflect.TypeOf(new(a.Provider)).Elem()
var rulesType = reflect.TypeOf(new(a.Rules)).Elem()
var scopeKey = new(struct{})

func scanProviders(scan, name string) []reflect.Type {
	return typereg.ListTypes(
		typereg.InPackage(scan),
		typereg.AssignableTo(providerType),
		hasScope(providerType, name),
	)
}

func scanRules(scan, name string) []reflect.Type {
	return typereg.ListTypes(
		typereg.InPackage(scan),
		typereg.AssignableTo(rulesType),
		hasScope(rulesType, name),
	)
}

func hasScope(marker reflect.Type, name string) typereg.Option {
	return func(t reflect.Type) bool {
		if name == "" {
			return true
		}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Type == marker {
				return f.Tag.Get("scope") == name
			}
		}
		return false
	}
}

type scope struct {
	next  *scope
	rules []rule
	cache map[reflect.Type][]furnisher
}

type rule interface {
	apply(start *scope, t reflect.Type, results []furnisher) []furnisher
}

type furnisher interface {
	furnish(start *scope, p unsafe.Pointer) error
}

func getScope(ctx context.Context) *scope {
	s := ctx.Value(scopeKey)
	if s == nil {
		return nil
	}
	return s.(*scope)
}

func (s *scope) furnish(pv reflect.Value) error {
	ps := s.getFurnishers(pv.Type())
	if len(ps) != 1 {
		return fmt.Errorf("%w: unable to furnish value of type %v, found %v", ErrInjectionFailed, pv.Type(), ps)
	}
	return ps[0].furnish(s, unsafe.Pointer(pv.UnsafeAddr()))
}

func (s *scope) getFurnishers(t reflect.Type) []furnisher {
	if furnishers, ok := s.cache[t]; ok {
		return furnishers
	}

	var furnishers []furnisher
	for cur := s; cur != nil; cur = cur.next {
		for _, r := range cur.rules {
			furnishers = r.apply(s, t, furnishers)
		}
	}

	s.cache[t] = furnishers
	return furnishers
}
