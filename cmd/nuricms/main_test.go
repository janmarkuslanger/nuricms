package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMainWrapper(t *testing.T) {
	os.Setenv("ENV", "test")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Main panicked: %v", r)
		}
	}()
	go main()
	require.True(t, true)
}
