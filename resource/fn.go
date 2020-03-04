package resource

import (
	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/errors"
	"net/http"
	"reflect"
)

var errType = reflect.TypeOf(new(error)).Elem()

func Func(f interface{}) http.Handler {
	ft := reflect.TypeOf(f)
	var res, err bool
	switch ft.NumOut() {
	case 0:
		err = false
		res = false
	case 1:
		err = ft.Out(0) == errType
		res = !err
	case 2:
		err = true
		res = true
	default:
		panic("wrong nmber of outputs")
	}
	return &funcHandler{
		invoke: reflect.ValueOf(f),
		res:    res,
		err:    err,
	}
}

type funcHandler struct {
	invoke   reflect.Value
	res, err bool
}

func (h *funcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	in, err := h.interpretRequest(r)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	out := h.invoke.Call(in)
	bs, err := h.marshalResponse(r, out)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	if len(bs) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}
	w.Write(bs)
}

func (h *funcHandler) interpretRequest(r *http.Request) ([]reflect.Value, error) {
	ctx := r.Context()
	ft := h.invoke.Type()
	in := make([]reflect.Value, ft.NumIn())
	for i := ft.NumIn() - 1; i >= 0; i-- {
		arg, err := di.Extract(ctx, ft.In(i))
		if err != nil {
			return nil, err
		}
		in[i] = reflect.ValueOf(arg)
	}
	return in, nil
}

func (h *funcHandler) marshalResponse(r *http.Request, out []reflect.Value) ([]byte, error) {
	var res interface{}
	if h.res {
		res = out[0].Interface()
		out = out[1:]
	}
	var err error
	if h.err {
		e, ok := out[0].Interface().(error)
		if ok {
			err = e
		}
	}
	if err != nil {
		return nil, err
	}
	if !h.res {
		return nil, nil
	}
	var bs []byte
	err = di.Require(r.Context(), func(e Encoding) {
		bs, err = e.Encode(r, res)
	})
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (h *funcHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if di.Require(r.Context(), func(h errors.Handler) {
		h.HandleError(w, r, err)
	}) != nil {
		w.WriteHeader(500)
	}
}
