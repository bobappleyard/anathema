package resource

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/bobappleyard/anathema/di"
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
	}, false, func(context.Context) error {
		return nil
	})
	r := httptest.NewRequest("GET", "/", nil)
	ctx := r.Context()
	scope := di.EnterScope(ctx)
	scope.AddProvider(di.Instance(Resource{1}))
	scope.AddProvider(di.Instance(Request{"example body"}))
	scope.AddProvider(di.Instance(JSONEncoding))
	r = r.WithContext(scope.Install(ctx))
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
