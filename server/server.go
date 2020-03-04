package server

import (
	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/errors"
	"github.com/bobappleyard/anathema/resource"
	"github.com/bobappleyard/anathema/router"
	"net/http"
)

type Resource interface {
	rimpl()
}

type Server struct {
	router   router.Router
	services []interface{}
}

func New() *Server {
	s := &Server{}
	s.AddService(func() errors.Handler { return errors.DefaultHandler })
	s.AddService(func() resource.Encoding { return resource.JSONEncoding })
	return s
}

func (s *Server) AddService(f interface{}) {
	s.services = append(s.services, f)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	for _, service := range s.services {
		ctx = di.Provide(ctx, service)
	}
	s.router.ServeHTTP(w, r.WithContext(ctx))
}
