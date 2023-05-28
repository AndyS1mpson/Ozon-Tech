// Wrapper for http request handlers
package srvwrapper

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// Describe methods for validating the request body
type Validator interface {
	Validate() error
}

// Describe a wrapper for request handlers
type Wrapper[Req any, Res any] struct {
	handler func(ctx context.Context, request Req) (Res, error)
}

// Creates a new wrapper instance
func New[Req any, Res any](handler func(ctx context.Context, request Req) (Res, error)) *Wrapper[Req, Res] {
	return &Wrapper[Req, Res] {
		handler: handler,
	}
}

// Standard request handler
func (w *Wrapper[Req, Res]) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var reqBody Req

	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		http.Error(rw, err.Error(),http.StatusBadRequest)
	}

	reqValidation, ok := any(req).(Validator)
	if ok {
		errValidation := reqValidation.Validate()
		if errValidation != nil {
			http.Error(rw, err.Error(),http.StatusBadRequest)
			return
		}
	}

	resp, err := w.handler(req.Context(), reqBody)
	if err != nil {
		log.Printf("executor fail: %s", err)
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		return
	}

	rawData, err := json.Marshal(&resp)
	if err != nil {
		http.Error(rw, err.Error(),http.StatusInternalServerError)
		return
	}

	_, _ = rw.Write(rawData)

}
