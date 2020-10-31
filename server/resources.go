package server

import (
	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/router"
	"net/http"
)

type serverResources struct {
	router *router.Router
}

type serverGroup struct {
	resources *serverResources
	route     *router.Route
}

func (s *serverResources) Group(path string) (a.Group, error) {
	route, err := router.ParseRoute(path)
	if err != nil {
		return nil, err
	}
	return &serverGroup{s, route}, nil
}

func (s *serverResources) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (g *serverGroup) Sub(name string) a.Group {
	return &serverGroup{g.resources, g.route.SubRoute(name)}
}

func (g *serverGroup) HEAD(f interface{}) {
	panic("implement me")
}

func (g *serverGroup) OPTIONS(f interface{}) {
	panic("implement me")
}

func (g *serverGroup) GET(f interface{}) {
	panic("implement me")
}

func (g *serverGroup) POST(f interface{}) {
	panic("implement me")
}

func (g *serverGroup) PUT(f interface{}) {
	panic("implement me")
}

func (g *serverGroup) PATCH(f interface{}) {
	panic("implement me")
}

func (g *serverGroup) DELETE(f interface{}) {
	panic("implement me")
}
