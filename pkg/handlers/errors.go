package handlers

import (
	"encoding/json"
	"net/http"
)

// JSONError is the type used to render JSON errors to clients
type JSONError struct {
	Message string `json:"error"`
	Code    int    `json:"httStatus"`
}

// renderInternalError is the fail safe response returned when we are dealing
// with an error generated by the application itself (i.e. somewhere something
// went very wrong).
// in prodcution we should attempt to return a sensible message with the
// error correlation id and log (and the id) the error for internal inspection
func renderInternalError(err error, w http.ResponseWriter) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// renderJSONError attempts to respond to the client by setting the
// header error code and disppays a message describing what the error was.
func renderJSONError(httpStatus int, msg string, w http.ResponseWriter) {
	w.WriteHeader(httpStatus)

	out := JSONError{
		Message: msg,
		Code:    httpStatus,
	}

	jsonMsg, err := json.Marshal(out)
	if err != nil {
		renderInternalError(err, w)
	}

	w.Write(jsonMsg)
}