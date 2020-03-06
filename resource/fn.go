package resource

import (
	"context"
	"errors"
	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/hterror"
	"net/http"
	"reflect"
	"strconv"
)

var errType = reflect.TypeOf(new(error)).Elem()

var (
	errNotFound   = hterror.WithStatusCode(http.StatusNotFound, errors.New("not found"))
	errBadRequest = hterror.WithStatusCode(http.StatusBadRequest, errors.New("bad request"))
)

func Func(f interface{}, requestBody bool, bind func(context.Context) (context.Context, error)) http.Handler {
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
		body:   requestBody,
		bind:   bind,
	}
}

type funcHandler struct {
	invoke         reflect.Value
	bind           func(context.Context) (context.Context, error)
	res, err, body bool
}

func (h *funcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	in, err := h.interpretRequest(r)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	out := h.invoke.Call(in)
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

func (h *funcHandler) interpretRequest(r *http.Request) ([]reflect.Value, error) {
	ctx, err := h.bind(r.Context())
	if err != nil {
		return nil, errNotFound
	}
	ft := h.invoke.Type()
	if h.body {
		rt := ft.In(1)
		req := reflect.New(rt)
		err = di.Require(ctx, func(e Encoding) error {
			return e.Decode(r, req.Interface())
		})
		if err != nil {
			return nil, errBadRequest
		}
		ctx = di.Insert(ctx, rt, req.Elem())
	}
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
