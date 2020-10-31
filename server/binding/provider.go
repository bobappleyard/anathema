package binding

import (
	"errors"
	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/router"
	"reflect"
)

var ErrBadRoute = errors.New("route mismatch")

type Rule struct {
	a.Service

	source InjectionSource
}

var resourceType = reflect.TypeOf(new(a.Resource)).Elem()

func (p *Rule) Inject(source InjectionSource) {
	p.source = source
}

func (p *Rule) Apply(t reflect.Type, providers []di.Provider) []di.Provider {
	if t.Kind() == reflect.Struct && t.AssignableTo(resourceType) {
		return append(providers, p)
	}
	return providers
}

func (p *Rule) Provide(s *di.Scope, v reflect.Value) error {
	injectors, err := p.source.GetInjectors(v.Type())
	if err != nil {
		return err
	}
	var m *router.Match
	if err := s.Require(&m); err != nil {
		return err
	}
	if !m.Route.EqualPath(injectors.Route) {
		return ErrBadRoute
	}
	for _, injector := range injectors.Fields {
		if err := injector.InjectField(s, v); err != nil {
			return err
		}
	}
	return nil
}
