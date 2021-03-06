package store

import (
	"errors"

	"github.com/satori/go.uuid"

	"github.com/Popcore/verisart/pkg/users"
)

type userStore struct {
	Users map[string]users.User
}

func newUserStore() userStore {
	return userStore{
		Users: make(map[string]users.User),
	}
}

// NewUser adds a new user to the Store
func (s *userStore) NewUser(email string, name string) (*users.User, error) {
	if _, ok := s.Users[email]; ok {
		return nil, errors.New("a user with the same email address already exists")
	}

	newUser := users.User{
		ID:    uuid.NewV4().String(),
		Email: email,
		Name:  name,
	}

	s.Users[email] = newUser

	return &newUser, nil
}
