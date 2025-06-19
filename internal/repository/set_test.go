package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewSet_RepositoriesNotNil(t *testing.T) {
	db := testutils.SetupTestDB(t)
	s := NewSet(db)
	assert.NotNil(t, s.Content)
	assert.NotNil(t, s.Field)
	assert.NotNil(t, s.Collection)
	assert.NotNil(t, s.ContentValue)
	assert.NotNil(t, s.Asset)
	assert.NotNil(t, s.User)
	assert.NotNil(t, s.Apikey)
	assert.NotNil(t, s.Webhook)
}
