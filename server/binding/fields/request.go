package fields

import (
	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/server/binding"
	"net/http"
	"reflect"
)

type RequestInjectors struct {
	a.Service

	encodings binding.EncodingSource
}

type requestSource func(r *http.Request, name string) string

type requestInjector struct {
	name     string
	field    reflect.StructField
	source   requestSource
	encoding binding.Encoding
}

var sources = []struct {
	tag string
	fn  requestSource
}{
	{"get", func(r *http.Request, name string) string { return r.URL.Query().Get(name) }},
	{"head", func(r *http.Request, name string) string { return r.Header.Get(name) }},
}

func (p *RequestInjectors) Inject(encodings binding.EncodingSource) {
	p.encodings = encodings
}

func (p *RequestInjectors) GetInjector(t reflect.Type, field reflect.StructField) (binding.FieldInjector, error) {
	for _, source := range sources {
		name := field.Tag.Get(source.tag)
		if name == "" {
			continue
		}
		encoding, err := p.encodings.GetEncoding(field.Type)
		if err != nil {
			return nil, err
		}
		return &requestInjector{name, field, source.fn, encoding}, nil
	}
	return nil, nil
}

func (p *requestInjector) InjectField(s *di.Scope, v reflect.Value) error {
	var r *http.Request
	if err := s.Require(&r); err != nil {
		return err
	}
	w, err := p.encoding.Decode(p.source(r, p.name))
	if err != nil {
		return err
	}

	v.FieldByIndex(p.field.Index).Set(w)
	return nil
}
