package di_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/bobappleyard/anathema/di"
)

type Item struct {
	ID int
}

type Repo interface {
	Get(id int) (Item, error)
}

type testRepo struct {
	item Item
}

func (r *testRepo) Get(id int) (Item, error) {
	return r.item, nil
}

func Example() {
	ctx := context.Background()

	ctx = di.Provide(ctx, func(r Repo) (Item, error) { return r.Get(1) })

	err := di.Require(ctx, func(item Item) { fmt.Println("item:", item) })
	fmt.Println("err:", err)
	ctx = di.ProvideValue(ctx, &testRepo{Item{1}})
	di.Require(ctx, func(item Item) { fmt.Println("item:", item) })
	err = di.Require(ctx, func(item Item) error {
		return nil
	})
	fmt.Println("err:", err)
	err = di.Require(ctx, func(item Item) error {
		return errors.New("test")
	})
	fmt.Println("err:", err)

	//Output: err: di_test.Repo not found
	// item: {1}
	// err: <nil>
	// err: test
}
