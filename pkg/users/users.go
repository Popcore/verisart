package users

// User is a type that represents a certificate owner or dealer.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UserManager is the interface that defines CRUD operations allowed
// on users.
type UserManager interface {

	// New generates a new user. Email address and name must be provided
	// while ID should be generated internally by the application.
	NewUser(email string, name string) (*User, error)
}
