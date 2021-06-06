package runner

import (
	"context"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type ServerConfig interface {
	ConfigureServer(ctx context.Context, server *http.Server) error
}

type Runner struct {
	Config []ServerConfig
}

func (r *Runner) Run(ctx context.Context) error {
	server, err := r.createServer(ctx)
	if err != nil {
		return err
	}
	return r.runServer(ctx, server)
}

func (r *Runner) createServer(ctx context.Context) (*http.Server, error) {
	server := new(http.Server)

	for _, c := range r.Config {
		err := c.ConfigureServer(ctx, server)
		if err != nil {
			return nil, err
		}
	}

	return server, nil
}

func (r *Runner) runServer(ctx context.Context, server *http.Server) error {
	g, ictx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ictx.Done()
		server.Shutdown(context.Background())
		return nil
	})

	g.Go(func() error {
		return server.ListenAndServe()
	})

	return g.Wait()
}
