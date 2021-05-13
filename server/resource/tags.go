package resource

import (
	"github.com/bobappleyard/anathema/server/a"
	"github.com/bobappleyard/anathema/server/router"
	"reflect"
)

type path struct {
	a.Service
}

func (*path) NameFromTag(tag reflect.StructTag) string {
	return tag.Get("path")
}

func (*path) ValueFromName(m router.Match, name string) string {
	return m.GetValue(name)
}

type get struct {
	a.Service
}

func (*get) NameFromTag(tag reflect.StructTag) string {
	return tag.Get("get")
}

func (*get) ValueFromName(m router.Match, name string) string {
	return m.Request.URL.Query().Get(name)
}

type head struct {
	a.Service
}

func (*head) NameFromTag(tag reflect.StructTag) string {
	return tag.Get("head")
}

func (*head) ValueFromName(m router.Match, name string) string {
	return m.Request.Header.Get(name)
}
