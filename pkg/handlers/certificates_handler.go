package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"goji.io/pat"

	cert "github.com/popcore/verisart_exercise/pkg/certificate"
)

// PostCertHandler accept requests dealing with the creation of
// new certificates
func PostCertHandler(store cert.Storer, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// parse payload
	decoder := json.NewDecoder(r.Body)
	newCert := cert.Certificate{}

	err := decoder.Decode(&newCert)
	if err != nil {
		renderJSONError(http.StatusBadRequest, "invalid json payload", w)
		return
	}

	// get user from header.
	// Here we simply read it from the Authorization header.
	userID := r.Header.Get("Authorization")
	userID = strings.TrimPrefix(userID, "Bearer ")

	if userID == "" {
		renderJSONError(http.StatusUnprocessableEntity, "user must be set in the Authorization header", w)
		return
	}

	newCert.OwnerID = userID

	// update storer
	savedCert, err := store.Create(newCert)
	if err != nil {
		renderJSONError(http.StatusBadRequest, err.Error(), w)
		return
	}

	// return new cert
	resp, err := json.Marshal(savedCert)
	if err != nil {
		renderInternalError(err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(resp)
	if err != nil {
		renderInternalError(err, w)
		return
	}
}

// PatchCertHandler accept requests dealing with the updating of
// exisitng certificates
func PatchCertHandler(store cert.Storer, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	certID := pat.Param(r, "id")

	// parse payload
	decoder := json.NewDecoder(r.Body)
	toUpdate := cert.Certificate{}

	err := decoder.Decode(&toUpdate)
	if err != nil {
		renderJSONError(http.StatusBadRequest, "invalid json payload", w)
	}

	// update storer
	updatedCert, err := store.Update(certID, toUpdate)
	if err != nil {
		renderJSONError(http.StatusBadRequest, err.Error(), w)
	}

	// return new cert
	resp, err := json.Marshal(updatedCert)
	if err != nil {
		renderInternalError(err, w)
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		renderInternalError(err, w)
	}
}

// DeleteCertHandler accept requests dealing with the removal of
// exisitng certificates
func DeleteCertHandler(store cert.Storer, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	certID := pat.Param(r, "id")

	// update storer
	err := store.Delete(certID)
	if err != nil {
		renderJSONError(http.StatusBadRequest, err.Error(), w)
	}

	w.WriteHeader(http.StatusNoContent)
}
