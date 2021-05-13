package hterror

import (
	"fmt"
	"net/http"
)

type Handler interface {
	HandleError(w http.ResponseWriter, r *http.Request, e error)
}

var DefaultHandler = defaultHandler{}

type defaultHandler struct{}

func (defaultHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Println(err)
	statusCode := 500
	if err, ok := err.(Error); ok {
		statusCode = err.StatusCode()
	}
	w.WriteHeader(statusCode)
}
