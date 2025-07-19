package env_test

import (
	"os"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/env"
	"github.com/stretchr/testify/assert"
)

func TestOsEnv(t *testing.T) {
	k := "Foo"
	v := "Bar"

	os.Setenv(k, v)

	t.Cleanup(func() {
		os.Unsetenv(k)
	})

	src := env.OsEnv{}
	res := src.Getenv(k)

	assert.Equal(t, res, v)
}
