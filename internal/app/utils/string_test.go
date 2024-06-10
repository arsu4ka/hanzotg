package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomString(t *testing.T) {
	length := 16

	s1 := RandomString(length)
	s2 := RandomString(length)

	assert.NotEqual(t, s1, s2)
}
