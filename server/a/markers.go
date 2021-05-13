package a

import (
	"github.com/bobappleyard/anathema/component"
	"github.com/bobappleyard/anathema/component/a"
)

type Service = a.Service

type Provider = a.Provider

type WebApplication interface {
	Provider
	webApplication()
}

type Resource interface {
	component.Marker
	resource()
}
