package load

import (
	"context"
	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/component"
	"github.com/bobappleyard/anathema/di"
	"reflect"
)

var (
	serviceType  = reflect.TypeOf(new(a.Service)).Elem()
	providerType = reflect.TypeOf(new(a.Provider)).Elem()
)

func ServerScope(ctx context.Context) (*di.Scope, error) {
	s := di.EnterScope(ctx)
	loadServices(s)
	err := loadProviders(s)
	if err != nil {
		return nil, err
	}
	err = loadRules(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func loadServices(s *di.Scope) {
	for _, service := range component.ListTypes(component.AssignableTo(serviceType)) {
		if service.Kind() != reflect.Struct {
			continue
		}
		s.AddRule(&serviceRule{service: service})
	}
}

func loadProviders(s *di.Scope) error {
	for _, provider := range component.ListTypes(component.AssignableTo(providerType)) {
		if provider.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < provider.NumMethod(); i++ {
			m := provider.Method(i)
			if m.PkgPath != "" {
				continue
			}
			if m.Name == "Inject" {
				continue
			}
			factory, err := di.Factory(m.Func.Interface())
			if err != nil {
				return err
			}
			s.AddRule(factory)
		}
	}
	return nil
}

func loadRules(s *di.Scope) error {
	var ruleServices []di.Rule
	if err := s.Require(&ruleServices); err != nil {
		return err
	}
	for _, rule := range ruleServices {
		s.AddRule(rule)
	}
	return nil
}
