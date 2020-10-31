package binding

import (
	"github.com/bobappleyard/anathema/a"
	"reflect"
)

type EncodingService struct {
	a.Service

	encodings []Encoding
}

func (s *EncodingService) Inject(encodings []Encoding) {
	s.encodings = encodings
}

func (s *EncodingService) GetEncoding(t reflect.Type) (Encoding, error) {
	for _, e := range s.encodings {
		if e.Accept(t) {
			return e, nil
		}
	}
	return nil, ErrUnknownEncoding
}
