package a

// A Group provides the methods required in registering resource methods.
type Group interface {
	Sub(name string) Group

	HEAD(f interface{})
	OPTIONS(f interface{})
	GET(f interface{})
	POST(f interface{})
	PUT(f interface{})
	PATCH(f interface{})
	DELETE(f interface{})
}
