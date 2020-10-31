package binding

import (
	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/di"
	"reflect"
)

type InjectionCache interface {
	Get(t reflect.Type) []FieldInjector
	Set(t reflect.Type, injectors []FieldInjector)
}

type FieldInjectionProvider interface {
	GetInjector(t reflect.Type, f reflect.StructField) (FieldInjector, error)
}

type FieldInjectionService struct {
	a.Service

	cache   InjectionCache
	sources []FieldInjectionProvider
}

type scopeProvidedField struct {
	field reflect.StructField
}

func (p *scopeProvidedField) InjectField(s *di.Scope, v reflect.Value) error {
	return s.RequireValue(v)
}

func (p *FieldInjectionService) Inject(cache InjectionCache, sources []FieldInjectionProvider) {
	p.cache = cache
	p.sources = sources
}

func (p *FieldInjectionService) GetInjectors(t reflect.Type) (Injectors, error) {
	if injectors := p.cache.Get(t); injectors != nil {
		return injectors, nil
	}
	injectors, err := p.buildInjectors(t)
	if err != nil {
		return nil, err
	}
	p.cache.Set(t, injectors)
	return injectors, nil
}

func (p *FieldInjectionService) buildInjectors(t reflect.Type) ([]FieldInjector, error) {
	var injectors []FieldInjector
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		injector, err := p.findInjector(t, field)
		if err != nil {
			return nil, err
		}
		injectors = append(injectors, injector)
	}
	return injectors, nil
}

func (p *FieldInjectionService) findInjector(t reflect.Type, field reflect.StructField) (FieldInjector, error) {
	for _, source := range p.sources {
		injector, err := source.GetInjector(t, field)
		if err != nil {
			return nil, err
		}
		if injector != nil {
			return injector, nil
		}
	}
	return &scopeProvidedField{field}, nil
}
