package user

// User is a type that represents an artworks owner or dealer.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
