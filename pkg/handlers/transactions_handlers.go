package handlers

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"

	cert "github.com/popcore/verisart_exercise/pkg/certificate"
	store "github.com/popcore/verisart_exercise/pkg/store"
)

// PostTransferHandler deals with requests that attempt to
// create a new certificate transfer.
func PostTransferHandler(s store.Storer, w http.ResponseWriter, r *http.Request) *HTTPError {
	certID := pat.Param(r, "id")

	// parse transfer payload
	txInfo := cert.Transaction{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&txInfo)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err.Error())
	}

	// attemp to update certificate transfer
	trx, err := s.CreateTx(certID, txInfo)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := json.Marshal(trx)
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

// PatchTransferHandler deals with requests that attempt to
// finalized (i.e complete or reject) a certificate transfer.
func PatchTransferHandler(s store.Storer, w http.ResponseWriter, r *http.Request) *HTTPError {
	certID := pat.Param(r, "id")

	// parse transfer payload
	txInfo := cert.Transaction{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&txInfo)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err.Error())
	}

	if txInfo.Status != cert.Accepted {
		return newHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	trx, err := s.CreateTx(certID, txInfo)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := json.Marshal(trx)
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
