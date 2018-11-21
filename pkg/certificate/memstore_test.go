package certificate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewCert(t *testing.T) {
	mc := memStore{
		Certs: map[string]Certificate{},
	}

	mockCert := Certificate{
		Title:   "the-title",
		OwnerID: "the-owner-id",
		Year:    2018,
		Note:    "some-notes",
	}

	got, err := mc.Create(mockCert)
	assert.Nil(t, err)

	// sanity check
	assert.NotNil(t, got.ID)
	assert.Equal(t, got.Title, "the-title")
	assert.NotNil(t, got.CreatedAt)
	assert.Equal(t, got.OwnerID, "the-owner-id")
	assert.Equal(t, got.Year, 2018)
	assert.Equal(t, got.Note, "some-notes")
	assert.Nil(t, got.Transfer)

	// attempting to create the same certificate - or a certificate
	// with an id - should return an error
	_, err = mc.Create(*got)
	assert.NotNil(t, err)
	assert.Len(t, mc.Certs, 1)
}

func TestUpdateCert(t *testing.T) {
	mockCert := Certificate{
		ID:        "the-id",
		Title:     "the-title",
		CreatedAt: time.Now(),
		OwnerID:   "the-owner-id",
		Year:      2018,
		Note:      "some-notes",
	}

	mc := memStore{
		Certs: map[string]Certificate{
			"the-id": mockCert,
		},
	}

	mockCert.OwnerID = "new-owner-id"
	mockCert.Transfer = &transfer{
		To:     "another-user",
		Status: "in-progress",
	}

	got, err := mc.Update("the-id", mockCert)
	assert.Nil(t, err)
	assert.Equal(t, got.OwnerID, "new-owner-id")
	assert.Equal(t, got.OwnerID, mockCert.OwnerID)
	assert.Equal(t, got.Transfer, &transfer{
		To:     "another-user",
		Status: "in-progress",
	})
	assert.Equal(t, got.Transfer, mockCert.Transfer)

	// attempting to update a non existing certificate should return an error
	got, err = mc.Update("i-dont-exists", mockCert)
	assert.Nil(t, got)
	assert.NotNil(t, err)
}

func TestDeleteCert(t *testing.T) {
	mockCert := Certificate{
		ID:        "the-id",
		Title:     "the-title",
		CreatedAt: time.Now(),
		OwnerID:   "the-owner-id",
		Year:      2018,
		Note:      "some-notes",
	}

	mc := memStore{
		Certs: map[string]Certificate{
			"the-id": mockCert,
		},
	}

	err := mc.Delete(mockCert.ID)
	assert.Nil(t, err)
	assert.Len(t, mc.Certs, 0)

	// attempting to delete a non existing certificate should return an error
	err = mc.Delete("i-dont-exists")
	assert.NotNil(t, err)
}

func TestGetCert(t *testing.T) {
	mockCert1 := Certificate{
		ID:        "id1",
		Title:     "the-title1",
		CreatedAt: time.Now(),
		OwnerID:   "owner-id1",
		Year:      2018,
		Note:      "some-notes",
	}

	mockCert2 := Certificate{
		ID:        "id2",
		Title:     "the-title2",
		CreatedAt: time.Now(),
		OwnerID:   "owner-id1",
		Year:      2018,
		Note:      "some-notes",
	}

	mockCert3 := Certificate{
		ID:        "id3",
		Title:     "title3",
		CreatedAt: time.Now(),
		OwnerID:   "owner-id2",
		Year:      2018,
		Note:      "some-notes",
	}

	mc := memStore{
		Certs: map[string]Certificate{
			"id1": mockCert1,
			"id2": mockCert2,
			"id3": mockCert3,
		},
	}

	certs, err := mc.Get("owner-id1")
	assert.Nil(t, err)
	assert.Len(t, certs, 2)
}
