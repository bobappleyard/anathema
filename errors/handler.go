package errors

import (
	"net/http"
)

var DefaultHandler = defaultHandler{}

type defaultHandler struct{}

func (defaultHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	statusCode := 500
	if err, ok := err.(interface{ StatusCode() int }); ok {
		statusCode = err.StatusCode()
	}
	w.WriteHeader(statusCode)
}
