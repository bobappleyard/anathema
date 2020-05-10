// Package server affords the description of HTTP servers in terms of resource
// types.
//
// Resource Types
//
// A resource type is a struct type that embeds Resource. The fields of the
// struct correspond to aspects of the resource, and tags on those fields
// instruct the framework on what that correspondence is.
//
// The most important tag is on the embedded Resource. This is the "path" tag,
// and informs the framework about the set of URLs that this resource concerns
// itself with:
//
//	type ExampleResource struct {
//		server.Resource `path:"/example"`
//	}
//
// Enclose parts of the path in curly braces to bind that portion of the URL to
// a field on the resource of the same name:
//
//	type ResourceWithID struct {
//		server.Resource `path:"/another/{ID}"`
//
//		ID int
//	}
//
// Resource types also support binding fields to query parameters by tagging
// them with "get":
//
//	type ResourceWithQuery struct {
//		server.Resource `path:"/collection"`
//
//		Search string `get:"search"`
//	}
//
package server

import (
	"context"
	"github.com/bobappleyard/anathema/di"
	"github.com/bobappleyard/anathema/hterror"
	"github.com/bobappleyard/anathema/resource"
	"github.com/bobappleyard/anathema/router"
	"net/http"
)

// Resource should be embedded in your resource types to signal that they are
// so.
type Resource interface {
	rimpl()
}

type Server struct {
	router   router.Router
	services *di.Registry
}

// New creates a new server with some sensible defaults for the common case of
// a RESTful web service consuming and producing JSON.
func New() *Server {
	s := &Server{}
	s.AddService(func() hterror.Handler { return hterror.DefaultHandler })
	s.AddService(func() resource.Encoding { return resource.JSONEncoding })
	return s
}

// AddService registers the provided factory with the server so that the type
// that the factory creates is available during the processing of requests.
func (s *Server) AddService(f interface{}) {
	if s.services == nil {
		s.services, _ = di.New(context.Background())
	}
	s.services.Provide(f)
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reg := s.services.Extend()
	reg.Provide(func() *http.Request {
		return r
	})
	ctx = reg.Bind(ctx)
	s.router.ServeHTTP(w, r.WithContext(ctx))
}
