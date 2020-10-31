package router

import (
	"errors"
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

func (rt *Router) AddRoute(method string, r *Route) error {
	routesPtr, ok := rt.matchMethod(method)
	if !ok {
		return ErrInvalidMethod
	}
	n := len(r.segments)
	routes := *routesPtr
	if n >= len(*routesPtr) {
		newRoutes := make([][]*Route, n+1)
		copy(newRoutes, routes)
		*routesPtr = newRoutes
		routes = newRoutes
	}
	routes[n] = append(routes[n], r)
	return nil
}

func (rt *Router) Match(r *http.Request) (Match, error) {
	segments := splitPath(r.URL.Path)
	routesPtr, ok := rt.matchMethod(r.Method)
	if !ok {
		return Match{}, ErrInvalidMethod
	}
	routes := *routesPtr
	if len(segments) < len(routes) {
		for _, route := range routes[len(segments)] {
			m, ok := route.match(segments)
			if ok {
				return m, nil
			}
		}
	}
	return Match{}, ErrNoRoute
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

func splitPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}
