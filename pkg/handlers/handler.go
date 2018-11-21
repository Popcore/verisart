package handlers

import (
	"net/http"

	cert "github.com/popcore/verisart_exercise/pkg/certificate"
)

type handler func(s cert.Storer, w http.ResponseWriter, r *http.Request)

// Handler responsible for handling http requests. It holds the actual function
// that will be invoked as a request is receive and any other configuration
// required by the handlers to perform the required actions.
type Handler struct {
	S cert.Storer
	H handler
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.H(h.S, w, r)
}
