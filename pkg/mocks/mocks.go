package mocks

import (
	cert "github.com/popcore/verisart/pkg/certificate"
	"github.com/popcore/verisart/pkg/users"
)

// MockStore is a mock implementation of the Storer interface.MockStore
// It must be used for testing purposes only.
type MockStore struct {
	Err   error
	Certs []cert.Certificate
	Cert  cert.Certificate
	Txs   []cert.Transaction
	Tx    cert.Transaction
	User  users.User
}

// CreateCert mock
func (m MockStore) CreateCert(c cert.Certificate) (*cert.Certificate, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &c, nil
}

// UpdateCert mock
func (m MockStore) UpdateCert(id string, c cert.Certificate) (*cert.Certificate, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &c, nil
}

// DeleteCert mock
func (m MockStore) DeleteCert(id string) error {
	if m.Err != nil {
		return m.Err
	}

	return nil
}

// GetCerts mock
func (m MockStore) GetCerts(ownerID string) ([]cert.Certificate, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return []cert.Certificate{}, nil
}

// CreateTx mock
func (m MockStore) CreateTx(certID string, tx cert.Transaction) (*cert.Transaction, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &m.Tx, nil
}

// AcceptTx mock
func (m MockStore) AcceptTx(certID string) (*cert.Certificate, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &m.Cert, nil
}

// NewUser mock
func (m MockStore) NewUser(email string, name string) (*users.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &m.User, nil
}
