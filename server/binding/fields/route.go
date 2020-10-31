package fields

import (
	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/router"
	"github.com/bobappleyard/anathema/server/binding"
	"reflect"
)

type RoutePropertyInjectors struct {
	a.Service

	encodings binding.EncodingSource
}

type matchInjector struct {
	nameIndex int
	field     reflect.StructField
	encoding  binding.Encoding
}

func (p *RoutePropertyInjectors) Inject(encodings binding.EncodingSource) {
	p.encodings = encodings
}

func (p *RoutePropertyInjectors) GetInjector(t reflect.Type, field reflect.StructField) (binding.FieldInjector, error) {
	f, _ := t.FieldByName("Resource")
	route, err := router.ParseRoute(f.Tag.Get("path"))
	if err != nil {
		return nil, err
	}
	for i, name := range route.Names() {
		if field.Name != name {
			continue
		}
		encoding, err := p.encodings.GetEncoding(field.Type)
		if err != nil {
			return nil, err
		}
		return &matchInjector{i, field, encoding}, nil
	}
	return nil, nil
}

func (p *matchInjector) InjectField(s *di.Scope, v reflect.Value) error {
	var m router.Match
	if err := s.Require(&m); err != nil {
		return err
	}
	v, err := p.encoding.Decode(m.Values[p.nameIndex])
	if err != nil {
		return err
	}
	v.FieldByIndex(p.field.Index).Set(v)
	return nil
}
