package hterror

type Error interface {
	StatusCode() int
}

type withStatusCode struct {
	code int
	err  error
}

func WithStatusCode(code int, err error) error {
	return &withStatusCode{code, err}
}

func (e *withStatusCode) Error() string {
	return e.err.Error()
}

func (e *withStatusCode) Unwrap() error {
	return e.err
}

func (e *withStatusCode) StatusCode() int {
	return e.code
}
