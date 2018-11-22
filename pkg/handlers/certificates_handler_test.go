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
	mux := goji.NewMux()
	memStore := store.NewMemStore()
	mux.Handle(pat.Post("/certificates"), Handler{S: memStore, H: PostCertHandler})

	input := `{
    "title": "my-thing",
    "year": 1998,
    "note": "some notes"
	}`

	req, err := http.NewRequest("POST", "/certificates", strings.NewReader(input))
	req.Header.Set("Authorization", "Bearer abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestPostCertHandlerInvalidJSON(t *testing.T) {
	mux := goji.NewMux()
	memStore := store.NewMemStore()
	mux.Handle(pat.Post("/certificates"), Handler{S: memStore, H: PostCertHandler})

	input := `{
    "title": "my-thing",
    "year": 19"
	}`

	req, err := http.NewRequest("POST", "/certificates", strings.NewReader(input))
	req.Header.Set("Authorization", "Bearer abc")

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
	mux := goji.NewMux()
	memStore := store.NewMemStore()
	mux.Handle(pat.Post("/certificates"), Handler{S: memStore, H: PostCertHandler})

	input := `{
		"id": "my id",
    "title": "my-thing",
    "year": 1998,
    "note": "some notes"
	}`

	req, err := http.NewRequest("POST", "/certificates", strings.NewReader(input))
	req.Header.Set("Authorization", "Bearer abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestPostCertHandlerErrorNoUser(t *testing.T) {
	mux := goji.NewMux()
	memStore := store.NewMemStore()
	mux.Handle(pat.Post("/certificates"), Handler{S: memStore, H: PostCertHandler})

	input := `{
    "title": "my-thing",
    "year": 1998,
    "note": "some notes"
	}`

	expected := `{
	  "error": "user must be set in the Authorization header",
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
	mux := goji.NewMux()
	memStore := store.NewMemStore()

	toUpdate, err := memStore.CreateCert(cert.Certificate{
		Title: "my cert",
		Year:  2018,
	})
	assert.Nil(t, err)

	mux.Handle(pat.Patch("/certificates/:id"), Handler{S: memStore, H: PatchCertHandler})

	input := `{
    "title": "my new thing",
    "year": 2018,
    "note": "some notes"
	}`

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/certificates/%s", toUpdate.ID), strings.NewReader(input))
	req.Header.Set("Authorization", "Bearer abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestPatchCertHandlerInvalidJSON(t *testing.T) {
	mux := goji.NewMux()
	memStore := store.NewMemStore()

	toUpdate, err := memStore.CreateCert(cert.Certificate{
		Title: "my cert",
		Year:  2018,
	})
	assert.Nil(t, err)

	mux.Handle(pat.Patch("/certificates/:id"), Handler{S: memStore, H: PatchCertHandler})

	input := `this-is-not-valid-json`

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/certificates/%s", toUpdate.ID), strings.NewReader(input))
	req.Header.Set("Authorization", "Bearer abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestPatchCertHandlerInvalidCertID(t *testing.T) {
	mux := goji.NewMux()
	memStore := store.NewMemStore()

	mux.Handle(pat.Patch("/certificates/:id"), Handler{S: memStore, H: PatchCertHandler})

	input := `{
		"title": "my new thing",
		"year": 2018,
		"note": "some notes"
	}`

	req, err := http.NewRequest("PATCH", "/certificates/i-dont-exists", strings.NewReader(input))
	req.Header.Set("Authorization", "Bearer abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestDeleteCertHandlerOK(t *testing.T) {
	mux := goji.NewMux()
	memStore := store.NewMemStore()
	toDelete, err := memStore.CreateCert(cert.Certificate{
		Title: "my cert",
		Year:  2018,
	})
	assert.Nil(t, err)

	mux.Handle(pat.Delete("/certificates/:id"), Handler{S: memStore, H: DeleteCertHandler})

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/certificates/%s", toDelete.ID), nil)
	req.Header.Set("Authorization", "Bearer abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNoContent, recorder.Code)
}
