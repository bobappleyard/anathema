package app

import (
	"github.com/bobappleyard/anathema/component"
	"github.com/bobappleyard/anathema/component/registry"
	"github.com/bobappleyard/anathema/server/a"
	"reflect"
)

var resourceType = reflect.TypeOf(new(a.Resource)).Elem()

type GroupProvider interface {
	Group(path string) (a.Group, error)
}

type ResourceImplementation interface {
	Init(group a.Group)
}

type ResourceSet struct {
	resourceTypes []reflect.Type
}

func loadResources(scan string) ResourceSet {
	return ResourceSet{registry.ListTypes(
		registry.InPackage(scan),
		registry.AssignableTo(resourceType),
	)}
}

func (s *ResourceSet) Visit(groups GroupProvider) error {
	for _, t := range s.resourceTypes {
		resource := reflect.New(t).Interface().(a.Resource)
		resTag := component.Tag(resource)
		group, err := groups.Group(resTag.Get("path"))
		if err != nil {
			return err
		}
		if resource, ok := resource.(ResourceImplementation); ok {
			resource.Init(group)
		}
	}
	return nil
}
