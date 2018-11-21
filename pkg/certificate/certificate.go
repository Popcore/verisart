package certificate

import (
	"time"
)

// Certificate is a type that represents an artwork certificate.
// It contains information about its name, provenance, status etc.
type Certificate struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	OwnerID   string    `json:"ownerId"`
	Year      int       `json:"year"`
	Note      string    `json:"note,omitempty"`
	Transfer  *transfer `json:"transfer"`
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
	Update(id string, c Certificate) (*Certificate, error)

	// Delete removes a Certificate from the store. It returns an error if
	// the operation could not be completed.
	Delete(id string) error
}
