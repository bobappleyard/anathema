package server

import (
	"context"
	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/hterror"
	"github.com/bobappleyard/anathema/resource"
	"github.com/bobappleyard/anathema/router"
	"net/http"
)

type Resource interface {
	rimpl()
}

type Server struct {
	router   router.Router
	services *di.Registry
}

func New() *Server {
	s := &Server{}
	s.AddService(func() hterror.Handler { return hterror.DefaultHandler })
	s.AddService(func() resource.Encoding { return resource.JSONEncoding })
	return s
}

func (s *Server) AddService(f interface{}) {
	if s.services == nil {
		s.services, _ = di.New(context.Background())
	}
	s.services.Provide(f)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = s.services.Extend().Bind(ctx)
	s.router.ServeHTTP(w, r.WithContext(ctx))
}
