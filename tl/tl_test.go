package tl

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	writerType := reflect.TypeOf(new(io.Writer)).Elem()
	bufioType := reflect.TypeOf(new(bufio.Writer))

	fmt.Println(bufioType.Elem().Name())

	for _, t := range ListTypes() {
		if t.AssignableTo(writerType) {
			pkg, name := idInfo(t)
			fmt.Printf("%s.%s", pkg, name)
			for _, m := range methodsFor(t) {
				fmt.Print(" ", m.Name)
			}
			fmt.Println()
		}
	}
}

func methodsFor(t reflect.Type) []reflect.Method {
	var res []reflect.Method
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.PkgPath == "" {
			res = append(res, m)
		}
	}
	t = t.Elem()
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.PkgPath == "" {
			res = append(res, m)
		}
	}
	return res
}

func idInfo(t reflect.Type) (string, string) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath(), t.Name()
}
