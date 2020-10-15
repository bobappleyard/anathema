package a

import "github.com/bobappleyard/anathema/component"

type Provider interface {
	component.Marker
	provider()
}

type Resource interface {
	component.Marker
	resource()
}

type Service interface {
	component.Marker
	service()
}
