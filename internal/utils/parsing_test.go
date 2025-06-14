package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToUint_OkFalse(t *testing.T) {
	to, ok := StringToUint("id")

	assert.Equal(t, to, uint(0))
	assert.False(t, ok)
}

func TestStringToUint_OkTrue(t *testing.T) {
	to, ok := StringToUint("12")
	assert.Equal(t, to, uint(12))
	assert.True(t, ok)
}
