package application

import (
	"context"
	"testing"
)

type testRunner struct {
	hasRun bool
}

func (r *testRunner) Run(ctx context.Context) error {
	r.hasRun = true
	return nil
}

type testApp struct {
	r *testRunner
}

func (a testApp) Runner() Runner {
	return a.r
}

func TestRun(t *testing.T) {
	r := &testRunner{}
	err := Run(testApp{r})
	if err != nil {
		t.Error(err)
	}
	if !r.hasRun {
		t.Fail()
	}
}
