package certificate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewCert(t *testing.T) {
	mc := memStore{
		Certs: map[CertID]Certificate{},
	}

	mockTime := time.Now()
	mockCert := Certificate{
		ID:        "the-id",
		Title:     "the-title",
		CreatedAt: mockTime,
		OwnerID:   "the-owner-id",
		Year:      2018,
		Note:      "some-notes",
	}

	got, err := mc.Create(mockCert)
	assert.Nil(t, err)

	// sanity check
	assert.Equal(t, got.ID, CertID("the-id"))
	assert.Equal(t, got.Title, "the-title")
	assert.Equal(t, got.CreatedAt, mockTime)
	assert.Equal(t, got.OwnerID, "the-owner-id")
	assert.Equal(t, got.Year, 2018)
	assert.Equal(t, got.Note, "some-notes")
	assert.Nil(t, got.Transfer)

	// attempting to create the same certificate should return an error
	_, err = mc.Create(mockCert)
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
		Certs: map[CertID]Certificate{
			"the-id": mockCert,
		},
	}

	mockCert.OwnerID = "new-owner-id"
	mockCert.Transfer = &transfer{
		To:     "another-user",
		Status: "in-progress",
	}

	got, err := mc.Update(mockCert)
	assert.Nil(t, err)
	assert.Equal(t, got.OwnerID, "new-owner-id")
	assert.Equal(t, got.OwnerID, mockCert.OwnerID)
	assert.Equal(t, got.Transfer, &transfer{
		To:     "another-user",
		Status: "in-progress",
	})
	assert.Equal(t, got.Transfer, mockCert.Transfer)

	// attempting to update a non existing certificate should return an error
	mockCert.ID = "i-dont-exists"
	got, err = mc.Update(mockCert)
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
		Certs: map[CertID]Certificate{
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
