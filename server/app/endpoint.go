package app

import (
	"github.com/bobappleyard/anathema/component/di"
	"net/http"
	"reflect"
)

type endpoint struct {
	app                 *Config
	impl                reflect.Value
	recv, input, output reflect.Type
	hasErr              bool
}

var errType = reflect.TypeOf(new(error)).Elem()

func (g *serverGroup) addEndpoint(method string, f reflect.Value) {
	t := f.Type()
	e := &endpoint{app: g.app, impl: f}
	e.parseInputs(t)
	e.parseOutputs(t)

	r := g.route.WithHandler(e)
	if err := g.app.router.AddRoute(method, r); err != nil {
		panic(err)
	}
}

func (e *endpoint) parseInputs(t reflect.Type) {
	e.recv = t.In(0)
	switch t.NumIn() {
	case 1:
		// do nothing
	case 2:
		e.input = t.In(1)
	default:
		panic("wrong number of inputs to endpoint implementation")
	}
}

func (e *endpoint) parseOutputs(t reflect.Type) {
	switch t.NumOut() {
	case 0:
		// do nothing
	case 1:
		if t.Out(0) == errType {
			e.hasErr = true
			break
		}
		e.output = t.Out(0)
	case 2:
		e.output = t.Out(0)
		if t.Out(1) != errType {
			panic("second output of two-output form should be error")
		}
		e.hasErr = true
	default:
		panic("wrong number of inputs to endpoint implementation")
	}
}

func (e *endpoint) handleError(w http.ResponseWriter, r *http.Request, err error) {
	e.app.handler.HandleError(w, r, err)
}

func (e *endpoint) encode(r *http.Request, entity interface{}) ([]byte, string, error) {
	return e.app.encoding.Encode(r, entity)
}

func (e *endpoint) decode(r *http.Request, entity interface{}) error {
	return e.app.encoding.Decode(r, entity)
}

func (e *endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	in, err := e.parseRequest(r)
	if err != nil {
		e.handleError(w, r, err)
		return
	}
	out := e.impl.Call(in)
	e.marshalResponse(w, r, out)
}

func (e *endpoint) parseRequest(r *http.Request) ([]reflect.Value, error) {
	scope := di.GetScope(r.Context())

	recv, err := scope.RequireValue(e.recv)
	if err != nil {
		return nil, err
	}
	if e.input == nil {
		return []reflect.Value{recv}, nil
	}

	ent := reflect.New(e.input)
	if err := e.decode(r, ent.Interface()); err != nil {
		return nil, err
	}

	return []reflect.Value{recv, ent}, nil
}

func (e *endpoint) marshalResponse(w http.ResponseWriter, r *http.Request, out []reflect.Value) {
	if e.hasErr {
		errv := out[len(out)-1]
		if !errv.IsNil() {
			e.handleError(w, r, errv.Interface().(error))
			return
		}
	}
	if e.output == nil {
		return
	}

	buf, contentType, err := e.encode(r, out[0].Interface())
	if err != nil {
		e.handleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(buf)
}
