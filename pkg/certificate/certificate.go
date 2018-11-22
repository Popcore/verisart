package certificate

import (
	"time"
)

// Certificate is a type that represents an artwork certificate.
// It contains information about its name, provenance, status etc.
type Certificate struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	CreatedAt time.Time    `json:"createdAt"`
	OwnerID   string       `json:"ownerId"`
	Year      int          `json:"year"`
	Note      string       `json:"note,omitempty"`
	Transfer  *Transaction `json:"transfer"`
}

type CertManger interface {
	// CreateCert adds a new Certificate to the store. It returns the generated
	// certificate or an error if anything goes wrong.
	CreateCert(c Certificate) (*Certificate, error)

	// UpdateCert modifies an existing Certificate. It returns the updated certificate
	// or an error if anything goes wrong.
	UpdateCert(id string, c Certificate) (*Certificate, error)

	// DeleteCert removes a Certificate from the store. It returns an error if
	// the operation could not be completed.
	DeleteCert(id string) error

	// GetCerts returns the certificates belonging to the user identified by
	// the ownerID.
	GetCerts(ownerID string) ([]Certificate, error)
}
