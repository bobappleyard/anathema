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

type resourceInit interface {
	Init(Group)
}

func (s *Server) Resource(r Resource) Group {
	rt := reflect.TypeOf(r)
	bdg := binding.ForStruct(rt)

	f, ok := rt.FieldByName("Resource")
	if !ok {
		panic("could not find embedded interface")
	}
	path := f.Tag.Get("path")
	if path == "" {
		panic("empty path")
	}
	route, err := router.ParseRoute(path)
	if err != nil {
		panic(err)
	}

	bdg = bdg.Slice(route.Names())
	if !bdg.Defined() {
		panic("missing field defintions")
	}

	g := &resourceGroup{s, path, route, bdg, rt}
	if r, ok := r.(resourceInit); ok {
		r.Init(g)
	}

	return g
}

type resourceGroup struct {
	server       *Server
	path         string
	route        *router.Route
	bindings     binding.Binding
	resourceType reflect.Type
}

func (g *resourceGroup) Sub(name string) Group {
	r, err := router.ParseRoute(g.path + "/" + name)
	if err != nil {
		panic(err)
	}
	return &subResourceGroup{g.server, g, r, name}
}

func (g *resourceGroup) HEAD(f interface{})    { g.addRoute("HEAD", false, f) }
func (g *resourceGroup) OPTIONS(f interface{}) { g.addRoute("OPTIONS", false, f) }
func (g *resourceGroup) POST(f interface{})    { g.addRoute("POST", true, f) }
func (g *resourceGroup) PATCH(f interface{})   { g.addRoute("PATCH", true, f) }
func (g *resourceGroup) DELETE(f interface{})  { g.addRoute("DELETE", false, f) }
func (g *resourceGroup) GET(f interface{})     { g.addRoute("GET", false, f) }
func (g *resourceGroup) PUT(f interface{})     { g.addRoute("PUT", true, f) }

func (g *resourceGroup) bind(ctx context.Context) error {
	reg := di.GetRegistry(ctx)
	return di.Require(ctx, func(m router.Match) error {
		rv, err := g.bindings.FromStrings(m.Values)
		if err != nil {
			return err
		}
		reg.Insert(g.resourceType, rv)
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
	return &subResourceGroup{g.server, g.parent, r, suffix}
}

func (g *subResourceGroup) HEAD(f interface{})    { g.addRoute("HEAD", false, f) }
func (g *subResourceGroup) OPTIONS(f interface{}) { g.addRoute("OPTIONS", false, f) }
func (g *subResourceGroup) POST(f interface{})    { g.addRoute("POST", true, f) }
func (g *subResourceGroup) PATCH(f interface{})   { g.addRoute("PATCH", true, f) }
func (g *subResourceGroup) DELETE(f interface{})  { g.addRoute("DELETE", false, f) }
func (g *subResourceGroup) GET(f interface{})     { g.addRoute("GET", false, f) }
func (g *subResourceGroup) PUT(f interface{})     { g.addRoute("PUT", true, f) }

func (g *subResourceGroup) addRoute(method string, requestBody bool, f interface{}) {
	r := g.route.WithHandler(g.parent.handler(method, requestBody, f))
	err := g.server.router.AddRoute(method, r)
	if err != nil {
		panic(err)
	}
}
