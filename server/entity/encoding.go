package entity

import (
	"encoding/json"
	"errors"
	"github.com/bobappleyard/anathema/server/hterror"
	"io/ioutil"
	"net/http"
)

var ErrBadRequest = hterror.WithStatusCode(http.StatusBadRequest, errors.New("bad request"))

type Encoding interface {
	Decode(r *http.Request, entity interface{}) error
	Encode(r *http.Request, entity interface{}) ([]byte, string, error)
}

var JSONEncoding Encoding = jsonEncoding{}

type jsonEncoding struct{}

func (jsonEncoding) Decode(r *http.Request, entity interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return ErrBadRequest
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, entity)
}

func (jsonEncoding) Encode(r *http.Request, entity interface{}) ([]byte, string, error) {
	buf, err := json.Marshal(entity)
	return buf, "application/json", err
}
