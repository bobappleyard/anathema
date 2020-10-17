package di

import "reflect"

type instanceFurnisher struct {
	value reflect.Value
}

func (f *instanceFurnisher) Furnish(s *Scope, p reflect.Value) error {
	p.Elem().Set(f.value)
	return nil
}

type methodFurnisher struct {
	method reflect.Method
}

func (f *methodFurnisher) Furnish(s *Scope, p reflect.Value) error {
	t := f.method.Type
	args := make([]reflect.Value, t.NumIn())

	for i := range args {
		aptr := reflect.New(t.In(i))
		if err := s.Furnish(aptr); err != nil {
			return err
		}
		args[i] = aptr.Elem()
	}

	out := f.method.Func.Call(args)
	if len(out) == 2 && !out[1].IsNil() {
		return out[1].Interface().(error)
	}

	p.Elem().Set(out[0])
	return nil
}
