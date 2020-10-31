package router

import (
	"net/http"
	"reflect"
	"testing"
)

func TestRouter(t *testing.T) {
	r := new(Router)
	addRoute := func(path string) *Route {
		rt, _ := ParseRoute(path)
		_ = r.AddRoute("GET", rt)
		return rt
	}
	user := addRoute("user/{Name}")
	division := addRoute("/division/{CompanyName}/{DivisionID}")
	divUsers := addRoute("/division/{CompanyName}/{DivisionID}/users")

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
			rq, _ := http.NewRequest("GET", test.path, nil)
			m, _ := r.Match(rq)
			if m.Route != test.route {
				t.Errorf("wrong route")
			}
			if !reflect.DeepEqual(m.Values, test.values) {
				t.Errorf("got %v, expecting %v", m.Values, test.values)
			}
		})
	}
}
