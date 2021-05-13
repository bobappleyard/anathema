package main

import (
	"context"
	"github.com/bobappleyard/anathema/cmd/anathema/internal/application"
	"github.com/bobappleyard/anathema/component/di"
	"github.com/bobappleyard/anathema/component/load"
	"os"
)

//go:generate anathema application

func main() {
	scope := di.EnterScope(context.Background(), "application")
	if err := load.Services(di.GetScope(scope)); err != nil {
		panic(err)
	}

	var runner application.Runner
	if err := di.Furnish(scope, &runner); err != nil {
		panic(err)
	}

	if err := runner.Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
