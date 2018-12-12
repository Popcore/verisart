package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"goji.io"
	"goji.io/pat"

	cert "github.com/Popcore/verisart/pkg/certificate"
	mocks "github.com/Popcore/verisart/pkg/mocks"
)

func TestPostTransferHandlerOK(t *testing.T) {
	mux := goji.NewMux()
	memStore := mocks.MockStore{
		Err: nil,
		Tx: cert.Transaction{
			To:     "user@email.com",
			Status: "pending",
		},
	}
	mux.Handle(pat.Post("/certificates/:id/transfers"), Handler{S: memStore, H: PostTransferHandler})

	input := `{
		"email": "user@email.com",
		"status": "pending"
	}`

	expected := `{
		"email": "user@email.com",
		"status": "pending"
	}`

	req, err := http.NewRequest("POST", fmt.Sprintf("/certificates/mock-id/transfers"), strings.NewReader(input))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.JSONEq(t, expected, recorder.Body.String())
}

func TestPostTransferHandlerErrorInvalidJSON(t *testing.T) {
	mux := goji.NewMux()
	memStore := mocks.MockStore{
		Err: nil,
		Tx:  cert.Transaction{},
	}
	mux.Handle(pat.Post("/certificates/:id/transfers"), Handler{S: memStore, H: PostTransferHandler})

	input := `{
		"email": "user@email.com",
		`

	expected := `{
		"httpStatus": 400,
		"error": "invalid json payload"
	}`

	req, err := http.NewRequest("POST", fmt.Sprintf("/certificates/mock-id/transfers"), strings.NewReader(input))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.JSONEq(t, expected, recorder.Body.String())
}

func TestPostTransferHandlerStoreError(t *testing.T) {
	mux := goji.NewMux()
	memStore := mocks.MockStore{
		Err: errors.New("some error"),
		Tx:  cert.Transaction{},
	}
	mux.Handle(pat.Post("/certificates/:id/transfers"), Handler{S: memStore, H: PostTransferHandler})

	input := `{
		"email": "user@email.com",
		"status": "pending"
	}`

	expected := `{
		"httpStatus": 422,
		"error": "some error"
	}`

	req, err := http.NewRequest("POST", fmt.Sprintf("/certificates/mock-id/transfers"), strings.NewReader(input))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.JSONEq(t, expected, recorder.Body.String())
}

func TestPatchTransferHandlerOK(t *testing.T) {
	mux := goji.NewMux()
	createdAt, _ := time.Parse(time.RFC3339, "2018-11-21T12:00:00Z")

	memStore := mocks.MockStore{
		Err: nil,
		Tx: cert.Transaction{
			To:     "user@email.com",
			Status: "accepted",
		},
		Cert: cert.Certificate{
			ID:        "123abc",
			Title:     "the-cert-title",
			OwnerID:   "user@email.com",
			CreatedAt: createdAt,
			Year:      2001,
		},
	}
	mux.Handle(pat.Patch("/certificates/:id/transfers"), Handler{S: memStore, H: PatchTransferHandler})

	input := `{
		"email": "user@email.com",
		"status": "accepted"
	}`

	expected := `{
		"id": "123abc",
		"title": "the-cert-title",
		"ownerId": "user@email.com",
		"year" : 2001,
		"createdAt": "2018-11-21T12:00:00Z",
		"transfer": null
	}`

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/certificates/mock-id/transfers"), strings.NewReader(input))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, expected, recorder.Body.String())
}

func TestPatchTransferHandlerErrorInvalidStatus(t *testing.T) {
	mux := goji.NewMux()
	memStore := mocks.MockStore{
		Err: nil,
		Tx: cert.Transaction{
			To:     "user@email.com",
			Status: "accepted",
		},
	}
	mux.Handle(pat.Patch("/certificates/:id/transfers"), Handler{S: memStore, H: PatchTransferHandler})

	input := `{
		"email": "user@email.com",
		"status": "something-else"
	}`

	expected := `{
		"httpStatus": 422,
		"error": "transaction status can only be set to 'accepted' for now"
	}`

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/certificates/mock-id/transfers"), strings.NewReader(input))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.JSONEq(t, expected, recorder.Body.String())
}

func TestPatchTransferHandlerErrorInvalidJSON(t *testing.T) {
	mux := goji.NewMux()
	memStore := mocks.MockStore{
		Err: errors.New("some error"),
		Tx:  cert.Transaction{},
	}
	mux.Handle(pat.Patch("/certificates/:id/transfers"), Handler{S: memStore, H: PatchTransferHandler})

	input := `{
		"email": "user@email.com",
		"status": "accepted"
	}`

	expected := `{
		"httpStatus": 422,
		"error": "some error"
	}`

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/certificates/mock-id/transfers"), strings.NewReader(input))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.JSONEq(t, expected, recorder.Body.String())
}
