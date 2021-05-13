package resource

import (
	"errors"
	"fmt"
	"github.com/bobappleyard/anathema/component/di"
	"github.com/bobappleyard/anathema/server/a"
	"github.com/bobappleyard/anathema/server/router"
	"reflect"
)

var resourceType = reflect.TypeOf(new(a.Resource)).Elem()

var ErrNoFieldEncodingFound = errors.New("no encoding found for field")

type FieldEncoding interface {
	Accept(t reflect.Type) bool
	Decode(s string, v reflect.Value) error
}

type TagSource interface {
	NameFromTag(tag reflect.StructTag) string
	ValueFromName(m router.Match, name string) string
}

type tagRule struct {
	a.Service

	encodings []FieldEncoding
	sources   []TagSource
}

func (r *tagRule) Inject(encodings []FieldEncoding, sources []TagSource) {
	r.encodings = encodings
	r.sources = sources
}

func (r *tagRule) Apply(b di.Builder, t reflect.Type) {
	if t.Kind() != reflect.Struct || !t.AssignableTo(resourceType) {
		return
	}
	b.Mutator(r)
}

func (r *tagRule) Update(p di.Process, v reflect.Value) error {
	var match router.Match
	if err := p.Furnish(&match); err != nil {
		return err
	}

	t := v.Type()
	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if f.PkgPath != "" || f.Tag == "" {
			continue
		}
		if f.Type.AssignableTo(resourceType) {
			continue
		}

		if err := r.injectField(match, v.FieldByIndex(f.Index), f); err != nil {
			return err
		}
	}

	return nil
}

func (r *tagRule) injectField(match router.Match, v reflect.Value, f reflect.StructField) error {
	value := r.valueFromTag(match, f.Tag)
	if value == "" {
		return nil
	}

	encoding := r.encodingForType(f.Type)
	if encoding == nil {
		return fmt.Errorf("%v: %w", f.Type, ErrNoFieldEncodingFound)
	}

	if err := encoding.Decode(value, v); err != nil {
		return fmt.Errorf("decoding %q: %w", value, err)
	}

	return nil
}

func (r *tagRule) valueFromTag(match router.Match, tag reflect.StructTag) string {
	for _, source := range r.sources {
		name := source.NameFromTag(tag)
		if name == "" {
			continue
		}
		return source.ValueFromName(match, name)
	}
	return ""
}

func (r *tagRule) encodingForType(t reflect.Type) FieldEncoding {
	for _, encoding := range r.encodings {
		if encoding.Accept(t) {
			return encoding
		}
	}
	return nil
}
