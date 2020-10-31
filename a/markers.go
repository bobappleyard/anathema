package a

import "github.com/bobappleyard/anathema/component"

type Service interface {
	component.Marker
	service()
}

type Provider interface {
	Service
	provider()
}

type WebApplication interface {
	Provider
	service()
}

type Resource interface {
	component.Marker
	resource()
}
