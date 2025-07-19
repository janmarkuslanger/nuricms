package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type FailingEnv struct{}

func (te FailingEnv) Getenv(v string) string {
	return ""
}

type SuccessEnv struct{}

func (te SuccessEnv) Getenv(v string) string {
	envs := map[string]string{
		"JWT_SECRET": "hiitsme",
	}

	return envs[v]
}

func TestLoadEnv_Failing(t *testing.T) {
	tenv := FailingEnv{}
	_, err := LoadEnv(tenv)

	assert.EqualError(t, err, "JWT_SECRET must be set")
}

func TestLoadEnv_Success(t *testing.T) {
	tenv := SuccessEnv{}
	env, err := LoadEnv(tenv)

	assert.Equal(t, err, nil)
	assert.Equal(t, env.Secret, "hiitsme")
}
