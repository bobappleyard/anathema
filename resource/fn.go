package resource

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/hterror"
)

var errType = reflect.TypeOf(new(error)).Elem()

var (
	errNotFound   = hterror.WithStatusCode(http.StatusNotFound, errors.New("not found"))
	errBadRequest = hterror.WithStatusCode(http.StatusBadRequest, errors.New("bad request"))
)

func Func(f interface{}, requestBody bool, bind func(context.Context) error) http.Handler {
	ft := reflect.TypeOf(f)
	res, err := parseFunctionOutputs(ft)
	return &funcHandler{
		invoke: f,
		res:    res,
		err:    err,
		body:   requestBody,
		bind:   bind,
	}
}

type funcHandler struct {
	invoke         interface{}
	bind           func(context.Context) error
	res, err, body bool
}

func parseFunctionOutputs(ft reflect.Type) (res, err bool) {
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
	return res, err
}

func (h *funcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.interpretRequest(r)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	bs, contentType, err := h.marshalResponse(r, out)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	if len(bs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(bs)))
	w.Write(bs)
}

func (h *funcHandler) interpretRequest(r *http.Request) error {
	ctx := r.Context()
	reg := di.GetScope(ctx)
	err := h.bind(ctx)
	if err != nil {
		return errNotFound
	}
	ft := reflect.TypeOf(h.invoke)
	if h.body {
		rt := ft.In(1)
		req := reflect.New(rt)
		var e Encoding
		err = reg.Require(&e)
		if err != nil {
			return errBadRequest
		}
		err = e.Decode(r, req.Interface())
		if err != nil {
			return errBadRequest
		}
		reg.AddProvider(di.Instance(req.Elem()))
	}
	return reg.Require(h.invoke)
}

func (h *funcHandler) marshalResponse(r *http.Request, out []reflect.Value) ([]byte, string, error) {
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
		return nil, "", err
	}
	if !h.res {
		return nil, "", nil
	}
	var bs []byte
	var contentType string
	err = di.Require(r.Context(), func(e Encoding) {
		bs, contentType, err = e.Encode(r, res)
	})
	if err != nil {
		return nil, "", err
	}
	return bs, contentType, nil
}

func (h *funcHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if di.Require(r.Context(), func(h hterror.Handler) {
		h.HandleError(w, r, err)
	}) != nil {
		w.WriteHeader(500)
	}
}
