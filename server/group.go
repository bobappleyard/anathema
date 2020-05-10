package server

import (
	"context"
	"github.com/bobappleyard/anathema/binding"
	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/resource"
	"github.com/bobappleyard/anathema/router"
	"net/http"
	"reflect"
)

// A Group provides the methods required in registering resource methods.
type Group interface {
	Sub(name string) Group

	HEAD(f interface{})
	OPTIONS(f interface{})
	GET(f interface{})
	POST(f interface{})
	PUT(f interface{})
	PATCH(f interface{})
	DELETE(f interface{})
}

type groupImpl interface {
	Sub(name string) Group
	addRoute(method string, requestBody bool, f interface{})
}

type groupWrapper struct {
	groupImpl
}

func (g *groupWrapper) HEAD(f interface{})    { g.addRoute("HEAD", false, f) }
func (g *groupWrapper) OPTIONS(f interface{}) { g.addRoute("OPTIONS", false, f) }
func (g *groupWrapper) POST(f interface{})    { g.addRoute("POST", true, f) }
func (g *groupWrapper) PATCH(f interface{})   { g.addRoute("PATCH", true, f) }
func (g *groupWrapper) DELETE(f interface{})  { g.addRoute("DELETE", false, f) }
func (g *groupWrapper) GET(f interface{})     { g.addRoute("GET", false, f) }
func (g *groupWrapper) PUT(f interface{})     { g.addRoute("PUT", true, f) }

type resourceGroup struct {
	server       *Server
	path         string
	route        *router.Route
	pathB, getB  binding.Binding
	resourceType reflect.Type
}

func (g *resourceGroup) Sub(name string) Group {
	r, err := router.ParseRoute(g.path + "/" + name)
	if err != nil {
		panic(err)
	}
	return &groupWrapper{&subResourceGroup{g.server, g, r, name}}
}

func (g *resourceGroup) bind(ctx context.Context) error {
	reg := di.GetRegistry(ctx)
	return di.Require(ctx, func(m router.Match, r *http.Request) error {
		rv := reflect.New(g.resourceType)
		err := g.pathB.FromStrings(m.Values, rv)
		if err != nil {
			return err
		}
		err = g.getB.FromFunc(func(name string) (string, bool) {
			vs, ok := r.URL.Query()[name]
			if !ok || len(vs) == 0 {
				return "", false
			}
			return vs[0], true
		}, rv)
		if err != nil {
			return err
		}
		reg.Insert(g.resourceType, rv.Elem())
		return nil
	})
}

func (g *resourceGroup) handler(method string, requestBody bool, f interface{}) http.Handler {
	return resource.Func(f, requestBody, g.bind)
}

func (g *resourceGroup) addRoute(method string, requestBody bool, f interface{}) {
	r := g.route.WithHandler(g.handler(method, requestBody, f))
	err := g.server.router.AddRoute(method, r)
	if err != nil {
		panic(err)
	}
}

type subResourceGroup struct {
	server *Server
	parent *resourceGroup
	route  *router.Route
	suffix string
}

func (g *subResourceGroup) Sub(name string) Group {
	suffix := g.suffix + "/" + name
	r, err := router.ParseRoute(g.parent.path + "/" + suffix)
	if err != nil {
		panic(err)
	}
	return &groupWrapper{&subResourceGroup{g.server, g.parent, r, suffix}}
}

func (g *subResourceGroup) addRoute(method string, requestBody bool, f interface{}) {
	r := g.route.WithHandler(g.parent.handler(method, requestBody, f))
	err := g.server.router.AddRoute(method, r)
	if err != nil {
		panic(err)
	}
}
