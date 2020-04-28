package di

import (
	"context"
	"fmt"
	"reflect"
)

// NotFoundErr signals that a type is missing from a context.
type NotFoundErr struct {
	Context context.Context
	Type    reflect.Type
}

func (e *NotFoundErr) Error() string {
	return fmt.Sprintf("%s not found", e.Type)
}
