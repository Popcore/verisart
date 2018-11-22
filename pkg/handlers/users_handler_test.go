package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"goji.io"
	"goji.io/pat"

	cert "github.com/popcore/verisart_exercise/pkg/certificate"
	store "github.com/popcore/verisart_exercise/pkg/store"
)

func TestListUserCertsHandlerOK(t *testing.T) {
	mux := goji.NewMux()
	memStore := store.NewMemStore()
	memStore.NewUser("owner1@email.com", "joe blog")
	memStore.NewUser("owner2@email.com", "miss smith")

	_, err := memStore.CreateCert(cert.Certificate{
		Title:   "my cert1",
		OwnerID: "owner1@email.com",
		Year:    2018,
	})

	_, err = memStore.CreateCert(cert.Certificate{
		Title:   "my cert2",
		OwnerID: "owner1@email.com",
		Year:    2018,
	})

	_, err = memStore.CreateCert(cert.Certificate{
		Title:   "my cert3",
		OwnerID: "owner2@email.com",
		Year:    2018,
	})
	assert.Nil(t, err)

	mux.Handle(pat.Get("/users/:userId/certificates"), Handler{S: memStore, H: ListUserCertsHandler})

	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%s/certificates", "owner1@email.com"), nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

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
