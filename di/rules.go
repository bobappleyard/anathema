package di

import (
	"fmt"
	"reflect"
	"unsafe"
)

type providerRule struct {
	provider reflect.Type
	provided []reflect.Type
}

type directTypeFurnisher struct {
	rule     *providerRule
	provided reflect.Type
	method   int
}

func newProviderRule(provider reflect.Type) (*providerRule, error) {
	var provided []reflect.Type
	for i := 0; i < provider.NumMethod(); i++ {
		m := provider.Method(i)
		if m.Type.NumIn() != 1 {
			return nil, fmt.Errorf("%w: provider method %v has wrong number of inputs", ErrInjectionFailed, m)
		}
		switch m.Type.NumOut() {
		case 1, 2:
			provided = append(provided, m.Type.Out(0))
		default:
			return nil, fmt.Errorf("%w: provider method %v has wrong number of outputs", ErrInjectionFailed, m)
		}
	}
	return &providerRule{provider, provided}, nil
}

func (r *providerRule) apply(start *scope, t reflect.Type, results []furnisher) []furnisher {
	for i, p := range r.provided {
		if p == t {
			return append(results, &directTypeFurnisher{r, p, i})
		}
	}
	return results
}

func (r *providerRule) instantiateProvider(start *scope) (reflect.Value, error) {
	pv := reflect.New(r.provider)
	err := start.furnish(pv.Elem())
	return pv, err
}

func (f *directTypeFurnisher) furnish(start *scope, p unsafe.Pointer) error {
	provider, err := f.rule.instantiateProvider(start)
	if err != nil {
		return err
	}
	res := provider.Method(f.method).Call(nil)
	if len(res) == 2 && !res[1].IsNil() {
		return res[1].Interface().(error)
	}
	pv := reflect.NewAt(f.provided, p)
	pv.Elem().Set(res[0])
	return nil
}
