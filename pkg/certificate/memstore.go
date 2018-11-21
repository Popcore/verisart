package certificate

import (
	"errors"
)

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

	// ensure certficate does not already exist
	if _, ok := m.Certs[c.ID]; ok {
		return nil, errors.New("The certificate already exists")
	}

	m.Certs[c.ID] = c

	return &c, nil
}

// Update modifies an existing certificate in the MemStore
func (m *memStore) Update(c Certificate) (*Certificate, error) {

	if _, ok := m.Certs[c.ID]; !ok {
		return nil, errors.New("Certificate not found")
	}

	m.Certs[c.ID] = c

	return &c, nil
}

// Delete modifies an existing certificate in the MemStore
func (m *memStore) Delete(id CertID) error {

	if _, ok := m.Certs[id]; !ok {
		return errors.New("Certificate not found")
	}

	delete(m.Certs, id)

	return nil
}
