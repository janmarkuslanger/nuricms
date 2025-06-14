package service

import (
	"os"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/testutils"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

func TestNewSet_Success(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	repos := repository.NewSet(testDB)
	hr := plugin.NewHookRegistry()
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	s, err := NewSet(repos, hr)
	assert.NoError(t, err)
	assert.NotNil(t, s.Collection)
	assert.NotNil(t, s.Field)
	assert.NotNil(t, s.Content)
	assert.NotNil(t, s.ContentValue)
	assert.NotNil(t, s.Asset)
	assert.NotNil(t, s.User)
	assert.NotNil(t, s.Apikey)
	assert.NotNil(t, s.Webhook)
	assert.NotNil(t, s.Api)
}

func TestNewSet_MissingJWTSecret(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	repos := repository.NewSet(testDB)
	hr := plugin.NewHookRegistry()
	os.Unsetenv("JWT_SECRET")

	s, err := NewSet(repos, hr)
	assert.Nil(t, s)
	assert.EqualError(t, err, "JWT_SECRET must be set")
}
