package router

import (
	"context"
	"github.com/bobappleyard/anathema/di"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRouter(t *testing.T) {
	var foundRoute *Route
	var matchedSegments []string

	testFn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		di.Require(r.Context(), func(m Match) {
			foundRoute = m.Route
			matchedSegments = m.Values
		})
	})
	r := new(Router)
	addRoute := func(path string) *Route {
		rt, _ := ParseRoute(path)
		rt = rt.WithHandler(testFn)
		r.AddRoute("GET", rt)
		return rt
	}
	user := addRoute("user/{Name}")
	division := addRoute("/division/{CompanyName}/{DivisionID}")
	divUsers := addRoute("/division/{CompanyName}/{DivisionID}/users")
	ctx := context.Background()

	for _, test := range []struct {
		name   string
		path   string
		route  *Route
		values []string
	}{
		{
			name:   "user",
			path:   "/user/bob",
			route:  user,
			values: []string{"bob"},
		},
		{
			name:   "division",
			path:   "/division/banana/1",
			route:  division,
			values: []string{"banana", "1"},
		},
		{
			name:   "divusers",
			path:   "/division/banana/1/users",
			route:  divUsers,
			values: []string{"banana", "1"},
		},
		{
			name:   "missing subroute",
			path:   "/division/banana/1/timelines",
			route:  nil,
			values: nil,
		},
		{
			name:   "nowhere",
			path:   "nowhere",
			route:  nil,
			values: nil,
		},
		{
			name:   "long nowhere",
			path:   "nowhere/in/particular/just/to/trip/the/test",
			route:  nil,
			values: nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			foundRoute, matchedSegments = nil, nil
			rc := httptest.NewRecorder()
			rq, _ := http.NewRequestWithContext(ctx, "GET", test.path, nil)
			r.ServeHTTP(rc, rq)
			if foundRoute != test.route {
				t.Errorf("wrong route")
			}
			if !reflect.DeepEqual(matchedSegments, test.values) {
				t.Errorf("got %v, expecting %v", matchedSegments, test.values)
			}
		})
	}
}
