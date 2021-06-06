package di

import (
	"context"
	"reflect"
)

func Install(ctx context.Context, module interface{}) context.Context {
	v := reflect.ValueOf(module)
	var mi moduleInstaller
	for i := v.NumMethod() - 1; i >= 0; i-- {
		mi.analyzeFactoryMethod(v.Method(i))
	}
	return mi.install(ctx)
}

type moduleInstaller struct {
	factories []factory
}

var errType = reflect.TypeOf(new(error)).Elem()

func (mi *moduleInstaller) analyzeFactoryMethod(m reflect.Value) {
	switch m.Type().NumOut() {
	case 1:
		mi.guaranteedResultFactory(m)

	case 2:
		mi.errorResultFactory(m)

	default:
		panic("malformed factory method")
	}
}

func (mi *moduleInstaller) install(ctx context.Context) context.Context {
	return toContext(ctx, newFurnisher(mi.factories))
}

func (mi *moduleInstaller) guaranteedResultFactory(m reflect.Value) {
	mi.factories = append(mi.factories, factory{
		forType: m.Type().Out(0),
		impl: func(ctx context.Context) (reflect.Value, error) {
			args, err := FurnishArgs(ctx, m)
			if err != nil {
				return reflect.Value{}, err
			}
			return m.Call(args)[0], nil
		},
	})
}

func (mi *moduleInstaller) errorResultFactory(m reflect.Value) {
	if !m.Type().Out(1).AssignableTo(errType) {
		panic("malformed factory method")
	}
	mi.factories = append(mi.factories, factory{
		forType: m.Type().Out(0),
		impl: func(ctx context.Context) (reflect.Value, error) {
			args, err := FurnishArgs(ctx, m)
			if err != nil {
				return reflect.Value{}, err
			}
			res := m.Call(args)
			if !res[1].IsNil() {
				return reflect.Value{}, res[1].Interface().(error)
			}
			return res[0], nil
		},
	})
}
