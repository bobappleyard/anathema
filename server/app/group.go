package app

import (
	"github.com/bobappleyard/anathema/server/a"
	"github.com/bobappleyard/anathema/server/router"
	"reflect"
)

type serverGroup struct {
	app   *Config
	route *router.Route
}

func (c *Config) Group(path string) (a.Group, error) {
	route, err := router.ParseRoute(path)
	if err != nil {
		return nil, err
	}
	return &serverGroup{c, route}, nil
}

func (g *serverGroup) Sub(name string) a.Group {
	return &serverGroup{g.app, g.route.SubRoute(name)}
}

func (g *serverGroup) HEAD(f interface{}) {
	g.addEndpoint("HEAD", reflect.ValueOf(f))
}

func (g *serverGroup) OPTIONS(f interface{}) {
	g.addEndpoint("OPTIONS", reflect.ValueOf(f))
}

func (g *serverGroup) GET(f interface{}) {
	g.addEndpoint("GET", reflect.ValueOf(f))
}

func (g *serverGroup) POST(f interface{}) {
	g.addEndpoint("POST", reflect.ValueOf(f))
}

func (g *serverGroup) PUT(f interface{}) {
	g.addEndpoint("PUT", reflect.ValueOf(f))
}

func (g *serverGroup) PATCH(f interface{}) {
	g.addEndpoint("PATCH", reflect.ValueOf(f))
}

func (g *serverGroup) DELETE(f interface{}) {
	g.addEndpoint("DELETE", reflect.ValueOf(f))
}
