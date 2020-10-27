package component

import (
	"reflect"
	"strings"
)

var registry []reflect.Type

func RegisterType(t reflect.Type) {
	registry = append(registry, t)
}

type Option func(reflect.Type) bool

func ListTypes(options ...Option) []reflect.Type {
	var res []reflect.Type
	for _, t := range registry {
		if testOptions(t, options) {
			res = append(res, t)
		}
		res = append(res, t)
	}
	return res
}

func testOptions(t reflect.Type, options []Option) bool {
	for _, option := range options {
		if !option(t) {
			return false
		}
	}
	return true
}

func InPackage(pkg string) Option {
	return func(t reflect.Type) bool {
		p := t.PkgPath()
		if pkg == p {
			return true
		}
		if !strings.HasPrefix(p, pkg) {
			return false
		}
		// This is safe because if HasPrefix == true and pkg != p then there
		// must be more characters in p
		return p[len(pkg)] == '/'
	}
}

func AssignableTo(intf reflect.Type) Option {
	return func(t reflect.Type) bool {
		return t.AssignableTo(intf)
	}
}
