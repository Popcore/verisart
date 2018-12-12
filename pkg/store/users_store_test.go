package store

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Popcore/verisart/pkg/users"
)

func TestNewUserStore(t *testing.T) {
	got := newUserStore()
	expected := userStore{
		Users: make(map[string]users.User),
	}

	assert.Equal(t, expected, got)
}

func TestNewUserOK(t *testing.T) {
	u := newUserStore()
	user, err := u.NewUser("test@email.com", "test-user")

	assert.Nil(t, err)
	assert.Len(t, u.Users, 1)
	assert.Equal(t, "test@email.com", user.Email)
	assert.Equal(t, "test-user", user.Name)
	assert.NotNil(t, user.ID)
}

func TestNewUserError(t *testing.T) {
	u := newUserStore()
	_, err := u.NewUser("test@email.com", "test-user")
	assert.Nil(t, err)

	_, err = u.NewUser("test@email.com", "test-user")
	assert.NotNil(t, err)
	assert.Len(t, u.Users, 1)
	assert.Equal(t, "a user with the same email address already exists", err.Error())
}
