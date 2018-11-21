package certificate

import (
	"errors"
	"time"

	"github.com/satori/go.uuid"
)

// MemStore is the in-memory certificates store.
// It is a concrete implementation of the storer interface
type memStore struct {
	Certs map[string]Certificate
}

// NewMemStore returns a memStore instance.
func NewMemStore() Storer {
	return &memStore{
		Certs: make(map[string]Certificate),
	}
}

// Create adds a new certificate to the MemStore.
func (m *memStore) Create(c Certificate) (*Certificate, error) {

	// return error if the Certificate already includes and id.
	// Ensure user knows what he/she is doing
	if c.ID != "" {
		return nil, errors.New("The certificate cannot contain an ID before it is created")
	}

	c.ID = uuid.NewV4().String()
	c.CreatedAt = time.Now().UTC()
	m.Certs[c.ID] = c

	return &c, nil
}

// Update modifies an existing certificate in the MemStore
func (m *memStore) Update(id string, c Certificate) (*Certificate, error) {

	if _, ok := m.Certs[id]; !ok {
		return nil, errors.New("Certificate not found")
	}

	m.Certs[id] = c

	return &c, nil
}

// Delete modifies an existing certificate in the MemStore
func (m *memStore) Delete(id string) error {

	if _, ok := m.Certs[id]; !ok {
		return errors.New("Certificate not found")
	}

	delete(m.Certs, id)

	return nil
}

// Delete modifies an existing certificate in the MemStore
func (m *memStore) Get(ownerID string) ([]Certificate, error) {
	certs := []Certificate{}

	for _, v := range m.Certs {
		if v.OwnerID == ownerID {
			certs = append(certs, v)
		}
	}

	return certs, nil
}
