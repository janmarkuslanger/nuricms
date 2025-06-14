package repository

import (
	"fmt"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_Create_Success(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewUserRepository(db)

	u := &model.User{Email: "create@example.com"}
	err := repo.Create(u)
	assert.NoError(t, err)
	assert.NotZero(t, u.ID)
}

func TestUserRepository_Save_UpdateEmail(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewUserRepository(db)

	u := &model.User{Email: "old@example.com"}
	assert.NoError(t, repo.Create(u))

	u.Email = "new@example.com"
	err := repo.Save(u)
	assert.NoError(t, err)

	found, err := repo.FindByEmail("new@example.com")
	assert.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)

	_, err = repo.FindByEmail("old@example.com")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUserRepository_Delete_RemovesUser(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewUserRepository(db)

	u := &model.User{Email: "todelete@example.com"}
	assert.NoError(t, repo.Create(u))

	assert.NoError(t, repo.Delete(u))
	_, err := repo.FindByID(u.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUserRepository_FindByEmail_FoundAndNotFound(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewUserRepository(db)

	u := &model.User{Email: "findme@example.com"}
	assert.NoError(t, repo.Create(u))

	found, err := repo.FindByEmail("findme@example.com")
	assert.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)

	_, err = repo.FindByEmail("doesnotexist@example.com")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUserRepository_FindByID_FoundAndNotFound(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewUserRepository(db)

	u := &model.User{Email: "byid@example.com"}
	assert.NoError(t, repo.Create(u))

	found, err := repo.FindByID(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.Email, found.Email)

	_, err = repo.FindByID(9999)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUserRepository_List_Pagination(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewUserRepository(db)

	for i := 1; i <= 5; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		assert.NoError(t, repo.Create(&model.User{Email: email}))
	}

	page1, total, err := repo.List(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, page1, 2)

	page3, total3, err := repo.List(3, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total3)
	assert.Len(t, page3, 1)
}

func TestUserRepository_List_Empty(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewUserRepository(db)

	list, total, err := repo.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, list)
}
