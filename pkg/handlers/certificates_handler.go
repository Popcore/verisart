package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"goji.io/pat"

	cert "github.com/popcore/verisart_exercise/pkg/certificate"
	store "github.com/popcore/verisart_exercise/pkg/store"
)

// PostCertHandler accepts requests dealing with the creation of
// new certificates
func PostCertHandler(s store.Storer, w http.ResponseWriter, r *http.Request) *HTTPError {

	w.Header().Set("Content-Type", "application/json")

	// parse payload
	decoder := json.NewDecoder(r.Body)
	newCert := cert.Certificate{}

	err := decoder.Decode(&newCert)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, "invalid json payload")
	}

	// get user from header.
	// Here we simply read it from the Authorization header.
	userID := r.Header.Get("Authorization")
	userID = strings.TrimPrefix(userID, "Bearer ")

	if userID == "" {
		return newHTTPError(http.StatusUnprocessableEntity, "user must be set in the Authorization header")
	}

	newCert.OwnerID = userID

	// update storer
	savedCert, err := s.CreateCert(newCert)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, err.Error())
	}

	// return new cert
	resp, err := json.Marshal(savedCert)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(resp)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

// PatchCertHandler accepts requests dealing with the updating of
// exisitng certificates
func PatchCertHandler(s store.Storer, w http.ResponseWriter, r *http.Request) *HTTPError {
	w.Header().Set("Content-Type", "application/json")

	certID := pat.Param(r, "id")

	// parse payload
	decoder := json.NewDecoder(r.Body)
	toUpdate := cert.Certificate{}

	err := decoder.Decode(&toUpdate)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, "invalid json payload")
	}

	// update storer
	updatedCert, err := s.UpdateCert(certID, toUpdate)
	if err != nil {
		return newHTTPError(http.StatusNotFound, err.Error())
	}

	// return new cert
	resp, err := json.Marshal(updatedCert)
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

// DeleteCertHandler accepts requests dealing with the removal of
// exisitng certificates
func DeleteCertHandler(s store.Storer, w http.ResponseWriter, r *http.Request) *HTTPError {
	w.Header().Set("Content-Type", "application/json")

	certID := pat.Param(r, "id")

	// update storer
	err := s.DeleteCert(certID)
	if err != nil {
		return newHTTPError(http.StatusNotFound, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
