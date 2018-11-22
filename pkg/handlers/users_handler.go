package handlers

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"

	cert "github.com/popcore/verisart_exercise/pkg/certificate"
	store "github.com/popcore/verisart_exercise/pkg/store"
)

// ListUserCertsHandler accepts requests dealing with the removal of
// existing certificates.
func ListUserCertsHandler(s store.Storer, w http.ResponseWriter, r *http.Request) *HTTPError {
	userID := pat.Param(r, "userId")

	certs := []cert.Certificate{}

	certs, err := s.GetCerts(userID)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp, err := json.Marshal(certs)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

// NewUserHandler accepts requests dealing with the creation of
// new users.
func NewUserHandler(s store.Storer, w http.ResponseWriter, r *http.Request) *HTTPError {
	return nil
}
