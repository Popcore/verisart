package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	port := ":1234"
	s := New(port)

	assert.Equal(t, s.Address, port)
}
