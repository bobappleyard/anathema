package di

import (
	"fmt"
	"reflect"
)

type process struct {
	scope   *Scope
	created map[reflect.Type]reflect.Value
	pending []pendingWork
}

type pendingWork struct {
	t reflect.Type
	b builder
}

type valueProvider interface {
	RequireValue(t reflect.Type) (reflect.Value, error)
}

func furnish(p valueProvider, ptr interface{}) error {
	target := reflect.ValueOf(ptr)
	v, err := p.RequireValue(target.Type())
	if err != nil {
		return err
	}
	target.Elem().Set(v.Elem())
	return nil
}

func (p *process) run(b builder, t reflect.Type) (reflect.Value, error) {
	v, err := p.createTarget(b, t)
	if err != nil {
		return reflect.Value{}, err
	}

	for len(p.pending) != 0 {
		next := p.next()
		nextTarget := p.created[next.t]
		if err := p.updateTarget(next.b, nextTarget); err != nil {
			return reflect.Value{}, err
		}
	}

	return v, nil
}

func (p *process) Apply(b Builder, t reflect.Type) {
	p.scope.Apply(b, t)
}

func (p *process) Furnish(ptr interface{}) error {
	return furnish(p, ptr)
}

func (p *process) RequireValue(t reflect.Type) (reflect.Value, error) {
	if v, ok := p.created[t]; ok {
		return v, nil
	}

	var b builder
	p.Apply(&b, t)

	return p.createTarget(b, t)
}

func (p *process) createTarget(b builder, t reflect.Type) (reflect.Value, error) {
	constructors := b.getConstructors()
	if len(constructors) != 1 {
		err := fmt.Errorf("injecting type %v: %w", t, ErrInjectionFailed)
		return reflect.Value{}, err
	}

	v, err := constructors[0].Create(p)
	if err != nil {
		return reflect.Value{}, err
	}

	p.created[t] = v
	p.pending = append(p.pending, pendingWork{t, b})
	return v, nil
}

func (p *process) updateTarget(builder builder, v reflect.Value) error {
	if builder.complete {
		return nil
	}
	for _, m := range builder.mutators {
		if err := m.Update(p, v); err != nil {
			return err
		}
	}
	return nil
}

func (p *process) next() pendingWork {
	res := p.pending[len(p.pending)-1]
	p.pending = p.pending[:len(p.pending)-1]
	return res
}
