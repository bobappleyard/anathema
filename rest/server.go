package rest

import (
	"github.com/bobappleyard/anathema/application"

	"github.com/bobappleyard/anathema/rest/runner"
)

type Server struct{}

func (Server) Runner(cfgs []runner.ServerConfig) application.Runner {
	return &runner.Runner{Config: cfgs}
}
