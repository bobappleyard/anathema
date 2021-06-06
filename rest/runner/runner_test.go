package runner

//go:generate mockgen -package runner -destination mock_runner_test.go . ServerConfig

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestInstall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	cfg := NewMockServerConfig(ctrl)
	cfg.EXPECT().ConfigureServer(ctx, gomock.Any()).Return(nil)

	r := Runner{[]ServerConfig{cfg}}
	_, err := r.createServer(ctx)
	if err != nil {
		t.Error(err)
	}
}

func TestInstallError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	e := errors.New("error")

	cfg := NewMockServerConfig(ctrl)
	cfg.EXPECT().ConfigureServer(ctx, gomock.Any()).Return(e)

	r := Runner{[]ServerConfig{cfg}}
	_, err := r.createServer(ctx)
	if err != e {
		t.Error(err)
	}
}

func TestRunError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	e := errors.New("error")

	cfg := NewMockServerConfig(ctrl)
	cfg.EXPECT().ConfigureServer(ctx, gomock.Any()).Return(e)

	r := Runner{[]ServerConfig{cfg}}
	err := r.Run(ctx)
	if err != e {
		t.Error(err)
	}
}

func TestRunCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := Runner{}
	s := new(http.Server)
	s.Addr = ":0"
	r.runServer(ctx, s)
}

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())

	cfg := NewMockServerConfig(ctrl)
	cfg.EXPECT().ConfigureServer(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, s *http.Server) error {
			s.Addr = ":0"
			return nil
		})
	r := Runner{[]ServerConfig{cfg}}

	cancel()
	r.Run(ctx)
}
