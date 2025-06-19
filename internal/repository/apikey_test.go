package repository

import (
	"testing"
	"time"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestApikeyRepository_FindByToken_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewApikeyRepository(db)
	future := time.Now().Add(time.Hour)
	exp := &model.Apikey{Name: "KeyA", Token: "tokenA", ExpiresAt: &future}
	repo.Create(exp)
	got, err := repo.FindByToken("tokenA")
	assert.NoError(t, err)
	assert.Equal(t, exp.ID, got.ID)
}

func TestApikeyRepository_FindByToken_NotFound(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewApikeyRepository(db)
	_, err := repo.FindByToken("doesnotexist")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
