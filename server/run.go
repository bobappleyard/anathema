package server

import (
	"context"
	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/server/load"
	"net"
	"net/http"
	"reflect"
)

func Run() error {
	ctx := context.Background()

	s, err := New(ctx)
	if err != nil {
		return err
	}

	return s.ListenAndServe()
}

func New(ctx context.Context) (*http.Server, error) {
	scope, err := load.ServerScope(ctx)
	if err != nil {
		return nil, err
	}
	ctx = scope.Install(ctx)

	var app a.WebApplication
	if err := scope.Require(&app); err != nil {
		return nil, err
	}
	tag := getAppTag(app)

	resources := &serverResources{}
	if err := load.Resources(app, resources); err != nil {
		return nil, err
	}

	s := new(http.Server)
	s.Addr = tag.Get("addr")
	s.Handler = resources
	s.BaseContext = func(listener net.Listener) context.Context { return ctx }

	return s, nil
}

func getAppTag(app a.WebApplication) reflect.StructTag {
	at := reflect.TypeOf(app)
	f, _ := at.FieldByName("WebApplication")
	return f.Tag
}
