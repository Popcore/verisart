package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	goji "goji.io"
	"goji.io/pat"

	cert "github.com/popcore/verisart_exercise/pkg/certificate"
	store "github.com/popcore/verisart_exercise/pkg/store"
)

func TestPostCertHandlerOK(t *testing.T) {
	memStore := store.NewMemStore()
	memStore.NewUser("user@email.com", "joe blog")

	mux := goji.NewMux()
	mux.Handle(pat.Post("/certificates"), Handler{S: memStore, H: PostCertHandler})

	input := `{
    "title": "my-thing",
    "year": 1998,
    "note": "some notes"
	}`

	req, err := http.NewRequest("POST", "/certificates", strings.NewReader(input))
	req.Header.Set("X-User-Email", "user@email.com")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestPostCertHandlerInvalidJSON(t *testing.T) {
	memStore := store.NewMemStore()
	memStore.NewUser("user@email.com", "joe blog")

	mux := goji.NewMux()
	mux.Handle(pat.Post("/certificates"), Handler{S: memStore, H: PostCertHandler})

	input := `{
    "title": "my-thing",
    "year": 19"
	}`

	req, err := http.NewRequest("POST", "/certificates", strings.NewReader(input))
	req.Header.Set("X-User-Email", "user@email.com")

	assert.Nil(t, err)

	expected := `{
	  "error": "invalid json payload",
	  "httpStatus": 400
	}`

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.JSONEq(t, expected, recorder.Body.String())
}

func TestPostCertHandlerInvalidCert(t *testing.T) {
	memStore := store.NewMemStore()

	mux := goji.NewMux()
	mux.Handle(pat.Post("/certificates"), Handler{S: memStore, H: PostCertHandler})

	input := `{
		"id": "my id",
    "title": "my-thing",
    "year": 1998,
    "note": "some notes"
	}`

	req, err := http.NewRequest("POST", "/certificates", strings.NewReader(input))
	req.Header.Set("X-User-Email", "abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestPostCertHandlerErrorNoUser(t *testing.T) {
	memStore := store.NewMemStore()

	mux := goji.NewMux()
	mux.Handle(pat.Post("/certificates"), Handler{S: memStore, H: PostCertHandler})

	input := `{
    "title": "my-thing",
    "year": 1998,
    "note": "some notes"
	}`

	expected := `{
	  "error": "user must be set in the X-User-Email header",
	  "httpStatus": 422
	}
	`

	req, err := http.NewRequest("POST", "/certificates", strings.NewReader(input))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.JSONEq(t, expected, recorder.Body.String())
}

func TestPatchCertHandlerOK(t *testing.T) {
	memStore := store.NewMemStore()
	_, err := memStore.NewUser("user@email.com", "joe blog")
	assert.Nil(t, err)

	toUpdate, err := memStore.CreateCert(cert.Certificate{
		OwnerID: "user@email.com",
		Title:   "my cert",
		Year:    2018,
	})
	assert.Nil(t, err)

	mux := goji.NewMux()
	mux.Handle(pat.Patch("/certificates/:id"), Handler{S: memStore, H: PatchCertHandler})

	input := `{
    "title": "my new thing",
    "year": 2018,
    "note": "some notes"
	}`

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/certificates/%s", toUpdate.ID), strings.NewReader(input))
	req.Header.Set("X-User-Email", "user@email.com")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestPatchCertHandlerInvalidJSON(t *testing.T) {
	memStore := store.NewMemStore()
	memStore.NewUser("user@email.com", "joe blog")

	toUpdate, err := memStore.CreateCert(cert.Certificate{
		OwnerID: "user@email.com",
		Title:   "my cert",
		Year:    2018,
	})
	assert.Nil(t, err)

	mux := goji.NewMux()
	mux.Handle(pat.Patch("/certificates/:id"), Handler{S: memStore, H: PatchCertHandler})

	input := `this-is-not-valid-json`

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/certificates/%s", toUpdate.ID), strings.NewReader(input))
	req.Header.Set("X-User-Email", "abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestPatchCertHandlerInvalidCertID(t *testing.T) {
	memStore := store.NewMemStore()

	mux := goji.NewMux()
	mux.Handle(pat.Patch("/certificates/:id"), Handler{S: memStore, H: PatchCertHandler})

	input := `{
		"title": "my new thing",
		"year": 2018,
		"note": "some notes"
	}`

	req, err := http.NewRequest("PATCH", "/certificates/i-dont-exists", strings.NewReader(input))
	req.Header.Set("X-User-Email", "abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestDeleteCertHandlerOK(t *testing.T) {
	memStore := store.NewMemStore()
	memStore.NewUser("user@email.com", "joe blog")

	toDelete, err := memStore.CreateCert(cert.Certificate{
		OwnerID: "user@email.com",
		Title:   "my cert",
		Year:    2018,
	})
	assert.Nil(t, err)

	mux := goji.NewMux()
	mux.Handle(pat.Delete("/certificates/:id"), Handler{S: memStore, H: DeleteCertHandler})

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/certificates/%s", toDelete.ID), nil)
	req.Header.Set("X-User-Email", "abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNoContent, recorder.Code)
}
