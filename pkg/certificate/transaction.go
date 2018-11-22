package certificate

// Transaction is a type that represents a certificate transaction
// from one uer to another.
type Transaction struct {
	To     string         `json:"email"`
	Status transferStatus `json:"status"`
}

type transferStatus string

const (
	// Pending is a status that can be applied to a transaction
	// waiting for approval.
	Pending transferStatus = "pending"

	// Accepted is a status that can be applied to a transaction
	// that has been agreed.
	Accepted transferStatus = "accepted"

	// Rejected is a status that can be applied to a transaction
	// that has been declined. Currently unused.
	Rejected transferStatus = "rejected"
)

// Transferer s the interface that represents operations on certificate
// transactions.
type Transferer interface {

	// CreateTx returns a new peding transaction for a certificate
	// idnetified by its id.
	CreateTx(certID string, trx Transaction) (*Transaction, error)

	// AcceptTx finalizes a certificate transaction to a new user.
	AcceptTx(certID string) error
}
