package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	cert "github.com/popcore/verisart_exercise/pkg/certificate"
	"github.com/popcore/verisart_exercise/pkg/users"
)

func TestCreateNewCert(t *testing.T) {
	mc := memStore{
		Certs:     map[string]cert.Certificate{},
		userStore: newUserStore(),
	}

	mc.Users["owner@email.com"] = users.User{
		ID:    "the-user-id",
		Email: "owner@email.com",
		Name:  "foo bar",
	}

	mockCert := cert.Certificate{
		Title:   "the-title",
		OwnerID: "owner@email.com",
		Year:    2018,
		Note:    "some-notes",
	}

	got, err := mc.CreateCert(mockCert)
	assert.Nil(t, err)

	// sanity check
	assert.NotNil(t, got.ID)
	assert.Equal(t, got.Title, "the-title")
	assert.NotNil(t, got.CreatedAt)
	assert.Equal(t, got.OwnerID, "owner@email.com")
	assert.Equal(t, got.Year, 2018)
	assert.Equal(t, got.Note, "some-notes")
	assert.Nil(t, got.Transfer)

	// attempting to create the same certificate - or a certificate
	// with an id - should return an error
	_, err = mc.CreateCert(*got)
	assert.NotNil(t, err)
	assert.Len(t, mc.Certs, 1)
}

func TestUpdateCert(t *testing.T) {
	mockCert := cert.Certificate{
		ID:        "the-id",
		Title:     "the-title",
		CreatedAt: time.Now(),
		OwnerID:   "the-owner-id",
		Year:      2018,
		Note:      "some-notes",
	}

	mc := memStore{
		Certs: map[string]cert.Certificate{
			"the-id": mockCert,
		},
	}

	toUpdate := cert.Certificate{
		Title: "the-new-title",
		Note:  "some-new-notes",
	}

	got, err := mc.UpdateCert("the-id", toUpdate)
	assert.Nil(t, err)
	assert.Equal(t, got.Title, "the-new-title")
	assert.Equal(t, got.Note, "some-new-notes")

	// attempting to update a non existing certificate should return an error
	got, err = mc.UpdateCert("i-dont-exists", mockCert)
	assert.Nil(t, got)
	assert.NotNil(t, err)

	// attempting to change ownership should return an error
	got, err = mc.UpdateCert("the-id", cert.Certificate{
		OwnerID: "new-owner",
	})
	assert.NotNil(t, err)
	assert.Equal(t, "ownership can only be changed with a transfer", err.Error())

	// attempting to update the transaction should return an error
	got, err = mc.UpdateCert("the-id", cert.Certificate{
		Transfer: &cert.Transaction{
			To:     "another-user@email.com",
			Status: cert.Accepted,
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, "ownership can only be changed with a transfer", err.Error())
}

func TestDeleteCert(t *testing.T) {
	mockCert := cert.Certificate{
		ID:        "the-id",
		Title:     "the-title",
		CreatedAt: time.Now(),
		OwnerID:   "the-owner-id",
		Year:      2018,
		Note:      "some-notes",
	}

	mc := memStore{
		Certs: map[string]cert.Certificate{
			"the-id": mockCert,
		},
	}

	err := mc.DeleteCert(mockCert.ID)
	assert.Nil(t, err)
	assert.Len(t, mc.Certs, 0)

	// attempting to delete a non existing certificate should return an error
	err = mc.DeleteCert("i-dont-exists")
	assert.NotNil(t, err)
}

func TestGetCert(t *testing.T) {
	mockCert1 := cert.Certificate{
		ID:        "id1",
		Title:     "the-title1",
		CreatedAt: time.Now(),
		OwnerID:   "owner-id1",
		Year:      2018,
		Note:      "some-notes",
	}

	mockCert2 := cert.Certificate{
		ID:        "id2",
		Title:     "the-title2",
		CreatedAt: time.Now(),
		OwnerID:   "owner-id1",
		Year:      2018,
		Note:      "some-notes",
	}

	mockCert3 := cert.Certificate{
		ID:        "id3",
		Title:     "title3",
		CreatedAt: time.Now(),
		OwnerID:   "owner-id2",
		Year:      2018,
		Note:      "some-notes",
	}

	mc := memStore{
		Certs: map[string]cert.Certificate{
			"id1": mockCert1,
			"id2": mockCert2,
			"id3": mockCert3,
		},
	}

	certs, err := mc.GetCerts("owner-id1")
	assert.Nil(t, err)
	assert.Len(t, certs, 2)
}

func TestCreateTxOK(t *testing.T) {
	mockCert := cert.Certificate{
		ID:      "key1",
		Title:   "the-title",
		OwnerID: "owner1@email.com",
		Year:    2018,
		Note:    "some-notes",
	}

	mc := memStore{
		Certs: map[string]cert.Certificate{
			"key1": mockCert,
		},
		Txs:       map[string][]cert.Transaction{},
		userStore: newUserStore(),
	}

	mc.NewUser("owner1@email.com", "joe blog")
	mc.NewUser("owner2@email.com", "miss smith")

	tx := cert.Transaction{
		To: "owner2@email.com",
	}

	_, err := mc.CreateTx("i-dond-exist", tx)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "certificate not found. Please use a valid ID")

	got, err := mc.CreateTx("key1", tx)
	expected := &cert.Transaction{
		To:     "owner2@email.com",
		Status: cert.Pending,
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Len(t, mc.Txs["key1"], 1)
	assert.Equal(t, mc.Certs["key1"].Transfer, &mc.Txs["key1"][0])
}

func TestCreateTxErrorNoPendingTx(t *testing.T) {
	mockCert := cert.Certificate{
		ID:      "key1",
		Title:   "the-title",
		OwnerID: "owner1@email.com",
		Year:    2018,
		Note:    "some-notes",
		Transfer: &cert.Transaction{
			To:     "owner2@email.com",
			Status: cert.Pending,
		},
	}

	mc := memStore{
		Certs: map[string]cert.Certificate{
			"key1": mockCert,
		},
		Txs: map[string][]cert.Transaction{
			"key1": []cert.Transaction{
				{
					To:     "owner2@email.com",
					Status: cert.Pending,
				},
			},
		},
		userStore: newUserStore(),
	}

	mc.NewUser("owner1@email.com", "joe blog")
	mc.NewUser("owner2@email.com", "miss smith")
	mc.NewUser("owner3@email.com", "luke")

	tx := cert.Transaction{
		To: "owner3@email.com",
	}

	_, err := mc.CreateTx("key1", tx)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "A pending transaction for certificate key1 already exist")

	// ensure old values are unchanged
	assert.Equal(t, mc.Certs["key1"].Transfer.Status, cert.Pending)
	assert.Equal(t, mc.Certs["key1"].Transfer.To, "owner2@email.com")
	assert.Equal(t, mc.Txs["key1"][0].Status, cert.Pending)
	assert.Equal(t, mc.Txs["key1"][0].To, "owner2@email.com")
}

func TestAcceptTx(t *testing.T) {
	certKey := "key1"
	mockCert := cert.Certificate{
		ID:      certKey,
		Title:   "the-title",
		OwnerID: "the-owner-id",
		Year:    2018,
		Note:    "some-notes",
		Transfer: &cert.Transaction{
			To:     "another-user@email.com",
			Status: cert.Pending,
		},
	}

	mc := memStore{
		Certs: map[string]cert.Certificate{
			certKey: mockCert,
		},
		Txs: map[string][]cert.Transaction{
			certKey: []cert.Transaction{
				{
					To:     "another-user@email.com",
					Status: cert.Pending,
				},
			},
		},
	}

	_, err := mc.AcceptTx("i-don't-exist")
	assert.NotNil(t, err)
	assert.Equal(t, "certificate not found. Please use a valid ID", err.Error())

	got, err := mc.AcceptTx(certKey)
	assert.Nil(t, err)
	assert.Equal(t, *got, mc.Certs[certKey])
	assert.Equal(t, string(cert.Accepted), string(mc.Txs[certKey][0].Status))
	assert.Equal(t, string(cert.Accepted), string(mc.Certs[certKey].Transfer.Status))

	assert.Equal(t, "another-user@email.com", mc.Certs[certKey].OwnerID)
}

func TestAcceptTxErrorEmptyTx(t *testing.T) {
	certKey := "key1"
	mockCert := cert.Certificate{
		ID:       certKey,
		Title:    "the-title",
		OwnerID:  "the-owner-id",
		Year:     2018,
		Note:     "some-notes",
		Transfer: nil,
	}

	mc := memStore{
		Certs: map[string]cert.Certificate{
			certKey: mockCert,
		},
		Txs: map[string][]cert.Transaction{},
	}

	_, err := mc.AcceptTx(certKey)
	assert.NotNil(t, err)
	assert.Equal(t, "no transactions found", err.Error())
}

func TestAcceptTxErrorNoPendingTx(t *testing.T) {
	certKey := "key1"
	mockCert := cert.Certificate{
		ID:      certKey,
		Title:   "the-title",
		OwnerID: "the-owner-id",
		Year:    2018,
		Note:    "some-notes",
		Transfer: &cert.Transaction{
			To:     "another-user@email.com",
			Status: cert.Accepted,
		},
	}

	mc := memStore{
		Certs: map[string]cert.Certificate{
			certKey: mockCert,
		},
		Txs: map[string][]cert.Transaction{
			certKey: []cert.Transaction{
				{
					To:     "another-user@email.com",
					Status: cert.Accepted,
				},
			},
		},
	}

	_, err := mc.AcceptTx(certKey)
	assert.NotNil(t, err)
	assert.Equal(t, "no pending transactions found", err.Error())
}
