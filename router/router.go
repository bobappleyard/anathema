package router

import (
	"errors"
	"github.com/bobappleyard/anathema/di"
	"net/http"
	"strings"
)

type Router struct {
	routes [5][][]*Route
}

var (
	ErrNoRoute       = errors.New("no route")
	ErrInvalidMethod = errors.New("invalid method")
)

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := splitPath(r.URL.Path)
	routesp, _ := rt.matchMethod(r.Method)
	routes := *routesp
	if len(segments) < len(routes) {
		for _, route := range routes[len(segments)] {
			m, ok := route.match(segments)
			if !ok {
				continue
			}
			r = r.Clone(di.Provide(r.Context(), func() Match { return m }))
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func (rt *Router) matchMethod(method string) (*[][]*Route, bool) {
	var routes *[][]*Route
	switch method {
	case "GET":
		routes = &rt.routes[0]
	case "HEAD":
		routes = &rt.routes[1]
	case "POST":
		routes = &rt.routes[2]
	case "PUT":
		routes = &rt.routes[3]
	case "DELETE":
		routes = &rt.routes[4]
	default:
		return nil, false
	}
	return routes, true
}

func (rt *Router) AddRoute(method string, r *Route) error {
	routesp, ok := rt.matchMethod(method)
	if !ok {
		return ErrInvalidMethod
	}
	n := len(r.segments)
	routes := *routesp
	if n >= len(*routesp) {
		newRoutes := make([][]*Route, n+1)
		copy(newRoutes, routes)
		*routesp = newRoutes
		routes = newRoutes
	}
	routes[n] = append(routes[n], r)
	return nil
}

func splitPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}
