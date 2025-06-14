package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewUserRepository(db)
	u := &model.User{Email: "test@example.com"}
	assert.NoError(t, repo.Create(u))
	found, err := repo.FindByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewUserRepository(db)
	_, err := repo.FindByEmail("missing@example.com")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
