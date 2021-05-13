package a

import "github.com/bobappleyard/anathema/component"

type Service interface {
	component.Marker
	service()
}

type Provider interface {
	component.Marker
	provider()
}

