// Wrapper for http request handlers
package srvwrapper

import (
	"bytes"
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
		rw.WriteHeader(http.StatusBadRequest)
		writeErrorText(rw, "parse request", err)
	}

	reqValidation, ok := any(req).(Validator)
	if ok {
		errValidation := reqValidation.Validate()
		if errValidation != nil {
			rw.WriteHeader(http.StatusBadRequest)
			writeErrorText(rw, "bad request", errValidation)
			return
		}
	}

	resp, err := w.handler(req.Context(), reqBody)
	if err != nil {
		log.Printf("executor fail: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		writeErrorText(rw, "exec handler", err)
		return
	}

	rawData, err := json.Marshal(&resp)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		writeErrorText(rw, "decode response", err)
		return
	}

	_, _ = rw.Write(rawData)

}

// Saves an error in the specified format to the response object
func writeErrorText(w http.ResponseWriter, text string, err error) {
	buf := bytes.NewBufferString(text)

	buf.WriteString(": ")
	buf.WriteString(err.Error())
	buf.WriteByte('\n')

	w.Write(buf.Bytes())
}