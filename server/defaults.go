package server

import (
	"github.com/bobappleyard/anathema/server/entity"
	"github.com/bobappleyard/anathema/server/hterror"
)

type WebApplicationDefaults struct {
}

func (s *WebApplicationDefaults) GetEncoding() entity.Encoding {
	return entity.JSONEncoding
}

func (s *WebApplicationDefaults) GetErrorHandler() hterror.Handler {
	return hterror.DefaultHandler
}


