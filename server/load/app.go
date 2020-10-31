package load

import (
	"errors"
	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/component"
	"reflect"
)

var ErrInvalidResource = errors.New("invalid resource")

var resourceType = reflect.TypeOf(new(a.Resource)).Elem()

type GroupProvider interface {
	Group(path string) (a.Group, error)
}

type ResourceImplementation interface {
	Init(group a.Group)
}

func Resources(app a.WebApplication, groups GroupProvider) error {
	tag := getAppTag(app)

	for _, resource := range component.ListTypes(
		component.InPackage(tag.Get("scan")),
		component.AssignableTo(resourceType),
	) {
		r, ok := reflect.New(resource).Interface().(ResourceImplementation)
		if !ok {
			return ErrInvalidResource
		}
		resTag := getResourceTag(resource)
		group, err := groups.Group(resTag.Get("path"))
		if err != nil {
			return err
		}
		r.Init(group)
	}

	return nil
}

func getAppTag(app a.WebApplication) reflect.StructTag {
	at := reflect.TypeOf(app)
	f, _ := at.FieldByName("WebApplication")
	return f.Tag
}

func getResourceTag(res reflect.Type) reflect.StructTag {
	f, _ := res.FieldByName("Resource")
	return f.Tag
}
