package certificate

import (
	"time"
)

// CertID is a type alias representing certificates IDs
type CertID string

// Certificate is a type that represents an artwork certificate.
// It contains information about its name, provenance, status etc.
type Certificate struct {
	ID        CertID    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	OwnerID   string    `json:"ownerId"`
	Year      int64     `json:"year"`
	Note      string    `json:"note,omitempty"`
	Transfer  transfer  `json:"transfer"`
}

type transfer struct {
	To     string `json:"email"`
	Status string `json:"status"`
}

// Storer is the interface that represents CRUD operations a certificate
// store must implement. Strictly speaking not required but good to have
// in case we decide to replace out in-memory sotre with a real database.
type Storer interface {
	// Create adds a new Certificate to the store. It returns the generated
	// certificate or an error if anything goes wrong.
	Create(c Certificate) (*Certificate, error)

	// Update modifies an existing Certificate. It returns the updated certificate
	// or an error if anything goes wrong.
	Update(c Certificate) (*Certificate, error)

	// Delete removes a Certificate from the store. It returns an error if
	// the operation could not be completed.
	Delete(id CertID) error
}

// MemStore is the in-memory certificates store.
// It is a concrete implementation of the storer interface
type memStore struct {
	Certs map[CertID]Certificate
}

// NewMemStore returns a memStore instance.
func NewMemStore() Storer {
	return &memStore{
		Certs: make(map[CertID]Certificate),
	}
}

// Create adds a new certificate to the MemStore
func (m *memStore) Create(c Certificate) (*Certificate, error) {
	return nil, nil
}

// Update modifies an existing certificate in the MemStore
func (m *memStore) Update(c Certificate) (*Certificate, error) {
	return nil, nil
}

// Delete modifies an existing certificate in the MemStore
func (m *memStore) Delete(id CertID) error {
	return nil
}
