package main

//go:generate anathema application

import (
	"github.com/bobappleyard/anathema/server"
	"github.com/bobappleyard/anathema/server/a"
)

type application struct {
	a.WebApplication `listen:":8888"`
	server.WebApplicationDefaults
}

func main() {
	panic(server.Run(application{}))
}
