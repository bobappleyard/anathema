package di

import (
	"context"
	"testing"
)

type testService struct {
	value int
}

type testModule struct{}

func (testModule) CreateInt() int {
	return 10
}

func (testModule) CreateService(x int) (*testService, error) {
	return &testService{x + 5}, nil
}

func TestModuleIntegration(t *testing.T) {
	var x struct {
		Int     int
		Service *testService
	}
	ctx := Install(context.Background(), testModule{})
	err := Furnish(ctx, &x)
	if err != nil {
		t.Error(err)
	}
	if x.Int != 10 {
		t.Fail()
	}
	if x.Service.value != 15 {
		t.Fail()
	}
}
