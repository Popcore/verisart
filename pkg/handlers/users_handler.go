package handlers

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"

	cert "github.com/popcore/verisart/pkg/certificate"
	store "github.com/popcore/verisart/pkg/store"
	users "github.com/popcore/verisart/pkg/users"
)

// ListUserCertsHandler accepts requests dealing with the listing of
// certificates that belong to the user ID specified in the URL.
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

	w.Header().Set("Content-Type", "application/json")
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
	// parse payload
	decoder := json.NewDecoder(r.Body)
	newUser := users.User{}

	err := decoder.Decode(&newUser)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, "invalid json payload")
	}

	if newUser.Email == "" || newUser.Name == "" {
		return newHTTPError(http.StatusUnprocessableEntity, "user email and name must be set in the request body")
	}

	resp, err := s.NewUser(newUser.Email, newUser.Name)
	if err != nil {
		return newHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(jsonResp)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
