package web

import (
	"encoding/json"
	"maps"
	"net/http"
)

type DtoModel interface {
	Repr() string
}

type Response struct {
	w http.ResponseWriter
	r *http.Request
}

func NewResponse(w http.ResponseWriter, r *http.Request) *Response {
	return &Response{
		w: w,
		r: r,
	}
}

func (response *Response) JsonResponseOk(data DtoModel, headers http.Header) error {
	maps.Copy(response.w.Header(), headers)
	response.w.Header().Set("Content-Type", "application/json")
	response.w.WriteHeader(http.StatusOK)

	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n') // more for pretty formatting
	_, err = response.w.Write(js)
	if err != nil {
		return err
	}
	return nil
}
