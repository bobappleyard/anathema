package resource

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Encoding interface {
	Decode(r *http.Request, entity interface{}) error
	Encode(r *http.Request, entity interface{}) ([]byte, error)
}

var JSONEncoding Encoding = jsonEncoding{}

type jsonEncoding struct{}

func (jsonEncoding) Decode(r *http.Request, entity interface{}) error {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, entity)
}

func (jsonEncoding) Encode(r *http.Request, entity interface{}) ([]byte, error) {
	return json.Marshal(entity)
}