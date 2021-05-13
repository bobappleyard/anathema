package server

import (
	"context"
	"github.com/bobappleyard/anathema/server/a"
	"github.com/bobappleyard/anathema/server/app"
	"net/http"
)

type WebApplicationContext interface {
	WebApplicationContext() context.Context
}

func Run(application a.WebApplication) error {
	s, err := New(application)
	if err != nil {
		return err
	}

	return s.ListenAndServe()
}

func New(application a.WebApplication) (*http.Server, error) {
	conf, err := app.LoadConfig(context.Background(), application)
	if err != nil {
		return nil, err
	}
	s := new(http.Server)
	conf.Install(s)
	return s, nil
}
