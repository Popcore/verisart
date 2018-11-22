package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"goji.io"
	"goji.io/pat"

	store "github.com/popcore/verisart_exercise/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestNewUserHandlerOK(t *testing.T) {
	mux := goji.NewMux()
	memStore := store.NewMemStore()

	mux.Handle(pat.Post("/users"), Handler{S: memStore, H: NewUserHandler})

	input := `{
    "email": "test@email.com",
    "name": "test-user"
	}`

	req, err := http.NewRequest("POST", "/users", strings.NewReader(input))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestNewUserHandlerErrInvalidPayload(t *testing.T) {
	mux := goji.NewMux()
	memStore := store.NewMemStore()

	mux.Handle(pat.Post("/users"), Handler{S: memStore, H: NewUserHandler})

	input := `{invalid-json`
	expected := `{
		"httpStatus": 400,
		"error": "invalid json payload"
	}`

	req, err := http.NewRequest("POST", "/users", strings.NewReader(input))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.JSONEq(t, expected, recorder.Body.String())
}

func TestNewUserHandlerErrMissingFields(t *testing.T) {
	mux := goji.NewMux()
	memStore := store.NewMemStore()

	mux.Handle(pat.Post("/users"), Handler{S: memStore, H: NewUserHandler})

	input := `{"email": "test@email.com"}`
	expected := `{
		"httpStatus": 422,
		"error": "user email and name must be set in the request body"
	}`

	req, err := http.NewRequest("POST", "/users", strings.NewReader(input))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.JSONEq(t, expected, recorder.Body.String())
}
