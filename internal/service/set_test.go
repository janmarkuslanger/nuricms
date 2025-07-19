package service

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/env"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewSet_Success(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	repos := repository.NewSet(testDB)
	hr := plugin.NewHookRegistry()

	env := env.Env{
		Secret: "testsecret",
	}

	s, err := NewSet(repos, hr, testDB, &env)
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
