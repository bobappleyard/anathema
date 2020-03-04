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

	g := &resourceGroup{s, route, bdg, rt}
	if r, ok := r.(resourceInit); ok {
		r.Init(g)
	}

	return g
}

type resourceGroup struct {
	server       *Server
	route        *router.Route
	bindings     binding.Binding
	resourceType reflect.Type
}

func (g *resourceGroup) Sub(name string) Group { return nil }

func (g *resourceGroup) HEAD(f interface{})    {}
func (g *resourceGroup) OPTIONS(f interface{}) {}
func (g *resourceGroup) POST(f interface{})    {}
func (g *resourceGroup) PATCH(f interface{})   {}
func (g *resourceGroup) DELETE(f interface{})  {}

func (g *resourceGroup) GET(f interface{}) {
	rh := resource.Func(f)
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx, err := g.bind(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}
		rh.ServeHTTP(w, r.WithContext(ctx))
	}
	r := g.route.WithHandler(http.HandlerFunc(h))
	err := g.server.router.AddRoute("GET", r)
	if err != nil {
		panic(err)
	}
}

func (g *resourceGroup) PUT(f interface{}) {
	rh := resource.Func(f)
	rt := reflect.TypeOf(f).In(1)
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx, err := g.bind(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}
		req := reflect.New(rt)
		err = di.Require(ctx, func(e resource.Encoding) error {
			return e.Decode(r, req.Interface())
		})
		if err != nil {
			w.WriteHeader(400)
			return
		}
		ctx = di.Insert(ctx, rt, req.Elem())
		rh.ServeHTTP(w, r.WithContext(ctx))
	}
	r := g.route.WithHandler(http.HandlerFunc(h))
	err := g.server.router.AddRoute("PUT", r)
	if err != nil {
		panic(err)
	}
}

func (g *resourceGroup) bind(ctx context.Context) (context.Context, error) {
	err := di.Require(ctx, func(m router.Match) error {
		rv, err := g.bindings.FromStrings(m.Values)
		if err != nil {
			return err
		}
		ctx = di.Insert(ctx, g.resourceType, rv)
		return nil
	})
	return ctx, err
}
