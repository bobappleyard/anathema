package resource

import (
	"github.com/bobappleyard/anathema/di"
	"net/http/httptest"
	"testing"
)

func TestFunc(t *testing.T) {
	type Resource struct {
		ID int
	}
	type Request struct {
		Content string
	}
	type Response struct {
		Code int
	}
	h := Func(func(res Resource, req Request) Response {
		return Response{1}
	})
	r := httptest.NewRequest("GET", "/", nil)
	ctx := r.Context()
	ctx = di.Provide(ctx, func() Resource { return Resource{1} })
	ctx = di.Provide(ctx, func() Request { return Request{"example body"} })
	ctx = di.Provide(ctx, func() Encoding { return JSONEncoding })
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Errorf("got status %d, expecting %d", w.Code, 200)
	}
	b := w.Body.String()
	if b != `{"Code":1}` {
		t.Errorf("got response body %q, expecting %q", b, `{"Code":1}`)
	}
}
