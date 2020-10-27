package a

import "github.com/bobappleyard/anathema/component"

type Provider interface {
	component.Marker
	provider()
}

type Service interface {
	component.Marker
	service()
}

type Resource interface {
	component.Marker
	resource()
}

type WebApplication interface {
	component.Marker
	service()
}
