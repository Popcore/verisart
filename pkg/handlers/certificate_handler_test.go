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
)

func TestPostCertHandlerOK(t *testing.T) {
	mux := goji.NewMux()
	memStore := cert.NewMemStore()
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
	memStore := cert.NewMemStore()
	mux.Handle(pat.Post("/certificates"), Handler{S: memStore, H: PostCertHandler})

	input := `{
    "title": "my-thing",
    "year": 19"
	}`

	req, err := http.NewRequest("POST", "/certificates", strings.NewReader(input))
	req.Header.Set("Authorization", "Bearer abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestPostCertHandlerInvalidCert(t *testing.T) {
	mux := goji.NewMux()
	memStore := cert.NewMemStore()
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

func TestPatchCertHandlerOK(t *testing.T) {
	mux := goji.NewMux()
	memStore := cert.NewMemStore()

	toUpdate, err := memStore.Create(cert.Certificate{
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
	memStore := cert.NewMemStore()

	toUpdate, err := memStore.Create(cert.Certificate{
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
	memStore := cert.NewMemStore()

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
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestDeleteCertHandlerOK(t *testing.T) {
	mux := goji.NewMux()
	memStore := cert.NewMemStore()
	toDelete, err := memStore.Create(cert.Certificate{
		Title: "my cert",
		Year:  2018,
	})
	assert.Nil(t, err)

	mux.Handle(pat.Delete("/certificates/:id"), Handler{S: memStore, H: DeleteCertHandler})

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/certificates/%s", toDelete.ID), strings.NewReader(""))
	req.Header.Set("Authorization", "Bearer abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteCertHandlerBadRequest(t *testing.T) {
	mux := goji.NewMux()
	memStore := cert.NewMemStore()

	mux.Handle(pat.Delete("/certificates/:id"), Handler{S: memStore, H: DeleteCertHandler})

	req, err := http.NewRequest("DELETE", "/certificates/i-dont'exist", strings.NewReader(""))
	req.Header.Set("Authorization", "Bearer abc")

	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
