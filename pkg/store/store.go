package store

import (
	"errors"
	"fmt"
	"time"

	"github.com/satori/go.uuid"

	cert "github.com/Popcore/verisart/pkg/certificate"
	"github.com/Popcore/verisart/pkg/users"
)

// Storer is the interface that defines CRUD operations allowed
// on certificates, transactions and users.
type Storer interface {
	users.UserManager
	cert.CertManager
	cert.Transferer
}

// MemStore is the in-memory concrete implementation of the storer interface.
// Internally it holds three maps: one for storing certificates, one for storing a
// list of transactions associated to certificates and a map for users.
type memStore struct {
	Certs map[string]cert.Certificate
	Txs   map[string][]cert.Transaction
	userStore
}

// NewMemStore returns a memStore instance.
func NewMemStore() Storer {
	return &memStore{
		Certs:     make(map[string]cert.Certificate),
		Txs:       make(map[string][]cert.Transaction),
		userStore: newUserStore(),
	}
}

// Create adds a new certificate to the MemStore.
func (m *memStore) CreateCert(c cert.Certificate) (*cert.Certificate, error) {

	// return error if the Certificate already includes and id since id are created by
	// the applcation
	if c.ID != "" {
		return nil, errors.New("The certificate cannot contain an ID before it is created")
	}

	// ensure user exists
	if _, ok := m.Users[c.OwnerID]; !ok {
		return nil, errors.New("The certificate must contain a valid user ID (aka email address). The email supplied did not match any user")
	}

	c.ID = uuid.NewV4().String()
	c.CreatedAt = time.Now().UTC()
	m.Certs[c.ID] = c

	return &c, nil
}

// Update modifies an existing certificate in the MemStore
func (m *memStore) UpdateCert(id string, c cert.Certificate) (*cert.Certificate, error) {

	toUpdate, ok := m.Certs[id]
	if !ok {
		return nil, errors.New("Certificate not found")
	}

	// reject changes to ownership or transactions
	if (c.OwnerID != "" && c.OwnerID != toUpdate.OwnerID) || c.Transfer != toUpdate.Transfer {
		return nil, errors.New("ownership can only be changed with a transfer")
	}

	// updatable fields are title, year and notes.
	// Id and createdAt should not be updated as are generated as internal metadata
	toUpdate.Title = c.Title
	toUpdate.Year = c.Year
	toUpdate.Note = c.Note

	m.Certs[id] = toUpdate

	return &toUpdate, nil
}

// Delete modifies an existing certificate in the MemStore.
func (m *memStore) DeleteCert(id string) error {

	if _, ok := m.Certs[id]; !ok {
		return errors.New("Certificate not found")
	}

	delete(m.Certs, id)

	return nil
}

// Delete modifies an existing certificate in the MemStore.
func (m *memStore) GetCerts(ownerID string) ([]cert.Certificate, error) {
	certs := []cert.Certificate{}

	for _, v := range m.Certs {
		if v.OwnerID == ownerID {
			certs = append(certs, v)
		}
	}

	return certs, nil
}

// CreateTx appends a new transaction and sets its status to "pending"
// to the list of the existing transaction associated to a certificate
// and updates the corresponding certificate information.
// It returns an error in case of failure.
func (m *memStore) CreateTx(certID string, tx cert.Transaction) (*cert.Transaction, error) {

	// ensure certificate exists before updating transactions
	// this will stop the transaction slice from growing indefinitely
	// if a certificate is deleted
	selectedCert, ok := m.Certs[certID]
	if !ok {
		return nil, errors.New("certificate not found. Please use a valid ID")
	}

	// ensure the transaction recipient exists
	if _, ok := m.Users[tx.To]; !ok {
		return nil, errors.New("invalid transaction recipient. The email address did not match any known user")
	}

	// update certificate transfer status and add transaction to the list
	// of existing ones and
	if canCreateTransaction(m.Txs[certID]) {
		tx.Status = cert.Pending

		selectedCert.Transfer = &tx

		m.Certs[certID] = selectedCert
		m.Txs[certID] = append([]cert.Transaction{tx}, m.Txs[certID]...)

		return &tx, nil
	}

	return nil, fmt.Errorf("A pending transaction for certificate %s already exist", certID)
}

// canCreateTransaction returns true if txs is empty or if the most
// recent transaction is not pending
func canCreateTransaction(txs []cert.Transaction) bool {
	if len(txs) == 0 {
		return true
	}

	return txs[0].Status != cert.Pending
}

// AcceptTx sets a transaction status to "accepted" and updates the
// corresponding certificate information. It returns an error in case
// of failure.
func (m *memStore) AcceptTx(certID string) (*cert.Certificate, error) {
	// ensure certificate exists
	selectedCert, ok := m.Certs[certID]
	if !ok {
		return nil, errors.New("certificate not found. Please use a valid ID")
	}

	lastTx, err := getLastPendingTx(m.Txs[certID])
	if err != nil {
		return nil, err
	}

	lastTx.Status = cert.Accepted
	selectedCert.Transfer = lastTx
	selectedCert.OwnerID = lastTx.To

	//"we must also set the new user id now"
	m.Certs[certID] = selectedCert
	m.Txs[certID][0] = *lastTx

	return &selectedCert, nil
}

// getLastPendingTx returns the last transaction if it exists and
// is not pending
func getLastPendingTx(txs []cert.Transaction) (*cert.Transaction, error) {
	if len(txs) == 0 {
		return nil, errors.New("no transactions found")
	}

	lastTx := txs[0]
	if lastTx.Status != cert.Pending {
		return nil, errors.New("no pending transactions found")
	}

	return &lastTx, nil
}
