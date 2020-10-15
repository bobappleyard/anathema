package server

import (
	"errors"
	"reflect"

	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/binding"
	"github.com/bobappleyard/anathema/router"
)

var (
	errMissingField = errors.New("missing resource field")
	errMissingTag   = errors.New("missing resourfce path")
)

// Resource registers a resource type with the server.
func (s *Server) Resource(r a.Resource) Group {
	rt := reflect.TypeOf(r)

	path, err := resourcePath(rt)
	if err != nil {
		panic(err)
	}

	route, err := router.ParseRoute(path)
	if err != nil {
		panic(err)
	}

	pathB, err := resourceBinding(rt, route.Names())
	if err != nil {
		panic(err)
	}
	getB := binding.Tag("get").ForStruct(rt)

	g := &groupWrapper{&resourceGroup{s, path, route, pathB, getB, rt}}
	if r, ok := r.(interface{ Init(Group) }); ok {
		r.Init(g)
	}

	return g
}

func resourcePath(rt reflect.Type) (string, error) {
	f, ok := rt.FieldByName("Resource")
	if !ok {
		return "", errMissingField
	}

	path := f.Tag.Get("path")
	if path == "" {
		return "", errMissingTag
	}

	return path, nil
}

func resourceBinding(rt reflect.Type, names []string) (binding.Binding, error) {
	bdg := binding.Fields().ForStruct(rt)

	bdg = bdg.Slice(names)
	if !bdg.Defined() {
		return binding.Binding{}, errMissingField
	}

	return bdg, nil
}
