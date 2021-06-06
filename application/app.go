package application

import (
	"context"

	"github.com/bobappleyard/anathema/di"
)

type Runner interface {
	Run(ctx context.Context) error
}

func Run(app interface{}) error {
	return RunContext(context.Background(), app)
}

func RunContext(ctx context.Context, app interface{}) error {
	ctx = di.Install(ctx, app)
	var runner Runner
	err := di.Furnish(ctx, &runner)
	if err != nil {
		return err
	}
	return runner.Run(ctx)
}
