package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToUint(t *testing.T) {
	to, ok := StringToUint("id")

	assert.Equal(t, to, uint(0))
	assert.False(t, ok)

	to, ok = StringToUint("12")
	assert.Equal(t, to, uint(12))
	assert.True(t, ok)
}
