package hterror

import (
	"fmt"
	"net/http"
)

var (
	ErrNotFound = &WithStatusCode{404}
)

type WithStatusCode struct {
	Status int
}

func (e *WithStatusCode) Error() string {
	return fmt.Sprintf("status %d", e.Status)
}

func (e *WithStatusCode) StatusCode() int {
	return e.Status
}

type Handler interface {
	HandleError(w http.ResponseWriter, r *http.Request, e error)
}
