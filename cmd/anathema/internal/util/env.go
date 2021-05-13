package util

import (
	"github.com/bobappleyard/anathema/component/di"
	"github.com/bobappleyard/anathema/server/a"
	"github.com/bobappleyard/anathema/server/resource"
	"os"
	"reflect"
)

type envField struct {
	a.Service
	encodings []resource.FieldEncoding
}

func (e *envField) Inject(encodings []resource.FieldEncoding) {
	e.encodings = encodings
}

func (e *envField) Apply(b di.Builder, t reflect.Type) {
	if t.Kind() != reflect.Struct {
		return
	}
	b.Mutator(e)
}

func (e *envField) Update(p di.Process, v reflect.Value) error {
	t := v.Type()
	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)

		if f.PkgPath != "" {
			continue
		}

		name := f.Tag.Get("env")
		if name == "" {
			continue
		}
		value := os.Getenv(name)
		if value == "" {
			continue
		}

		if err := e.updateField(v.FieldByIndex(f.Index), f.Type, value); err != nil {
			return err
		}
	}

	return nil
}

func (e *envField) updateField(f reflect.Value, t reflect.Type, value string) error {
	for _, encoding := range e.encodings {
		if encoding.Accept(t) {
			return encoding.Decode(value, f)
		}
	}
	return resource.ErrNoFieldEncodingFound
}
