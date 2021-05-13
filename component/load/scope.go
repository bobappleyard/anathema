package load

import (
	"github.com/bobappleyard/anathema/component/di"
	"github.com/bobappleyard/anathema/component/registry"
	"github.com/bobappleyard/anathema/server/a"
	"reflect"
)

func Services(s *di.Scope) error {
	load := scopeLoader{s, nil}
	load.markers()
	load.services()
	load.providers()
	load.rules()
	return load.err
}

var (
	serviceType  = reflect.TypeOf(new(a.Service)).Elem()
	providerType = reflect.TypeOf(new(a.Provider)).Elem()
)

type scopeLoader struct {
	s   *di.Scope
	err error
}

func (l *scopeLoader) markers() {
	l.addRule(&markerRule{})
}

func (l *scopeLoader) services() {
	if l.err != nil {
		return
	}

	for _, service := range registry.ListTypes(registry.AssignableTo(serviceType)) {
		if service.Kind() != reflect.Struct {
			continue
		}

		l.addRule(&serviceRule{service: reflect.PtrTo(service)})
	}
}

func (l *scopeLoader) providers() {
	if l.err != nil {
		return
	}

	for _, provider := range registry.ListTypes(registry.AssignableTo(providerType)) {
		if provider.Kind() != reflect.Struct {
			continue
		}

		provider = reflect.PtrTo(provider)
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
				l.err = err
				return
			}
			l.addRule(factory)
		}
	}
}

func (l *scopeLoader) rules() {
	if l.err != nil {
		return
	}

	var ruleServices []di.Rule
	if err := l.s.Furnish(&ruleServices); err != nil {
		l.err = err
		return
	}
	for _, rule := range ruleServices {
		l.addRule(rule)
	}
}

func (l *scopeLoader) addRule(r di.Rule) {
	l.s.AddRule(r)
}
