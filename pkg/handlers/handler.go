package handlers

import (
	"encoding/json"
	"net/http"

	store "github.com/popcore/verisart/pkg/store"
)

type handler func(s store.Storer, w http.ResponseWriter, r *http.Request) *HTTPError

// Handler is responsible for handling http requests. It holds the actual function
// that will be invoked as a request is received and any other configuration or
// services required by the handlers.
type Handler struct {
	S store.Storer
	H handler
}

// HTTPError is the error type returned by handlers when requests cannot
// be successfully completed. It contains the error HTTP status code and
// an error message.
type HTTPError struct {
	Code int    `json:"httpStatus"`
	Msg  string `json:"error"`
}

func newHTTPError(code int, msg string) *HTTPError {
	return &HTTPError{
		Code: code,
		Msg:  msg,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.S, w, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Code)

		errResp, innerErr := json.Marshal(err)
		if innerErr != nil {
			http.Error(w, innerErr.Error(), http.StatusInternalServerError)
		}

		w.Write(errResp)
	}
}
