package app

import (
	"context"
	"github.com/bobappleyard/anathema/component"
	"github.com/bobappleyard/anathema/component/di"
	"github.com/bobappleyard/anathema/component/load"
	"github.com/bobappleyard/anathema/server/a"
	"github.com/bobappleyard/anathema/server/entity"
	"github.com/bobappleyard/anathema/server/hterror"
	"github.com/bobappleyard/anathema/server/router"
	"net"
	"net/http"
	"reflect"
)

type Config struct {
	base      a.WebApplication
	ctx       context.Context
	resources ResourceSet
	router    *router.Router

	encoding entity.Encoding
	handler  hterror.Handler
}

func LoadConfig(ctx context.Context, app a.WebApplication) (*Config, error) {
	ctx = di.EnterScope(ctx, "application")
	if err := load.Services(di.GetScope(ctx)); err != nil {
		return nil, err
	}

	var conf Config
	if err := di.Furnish(ctx, &conf); err != nil {
		return nil, err
	}

	scan := component.Tag(app).Get("scan")
	if scan == "" {
		scan = reflect.TypeOf(app).PkgPath() + "/..."
	}
	resources := loadResources(scan)

	conf.base = app
	conf.ctx = ctx
	conf.router = new(router.Router)
	conf.resources = resources

	if err := resources.Visit(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *Config) Inject(encoding entity.Encoding, handler hterror.Handler) {
	c.encoding = encoding
	c.handler = handler
}

func (c *Config) Resources() ResourceSet {
	return c.resources
}

func (c *Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	match, err := c.router.Match(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ctx := di.EnterScope(r.Context(), "request")
	scope := di.GetScope(ctx)
	defer scope.Close()

	scope.AddRule(di.Instance(match))
	r = r.WithContext(ctx)

	match.Route.Handler().ServeHTTP(w, r)
}

func (c *Config) Install(s *http.Server) {
	s.Addr = component.Tag(c.base).Get("listen")
	s.Handler = c
	s.BaseContext = func(net.Listener) context.Context { return c.ctx }
}
