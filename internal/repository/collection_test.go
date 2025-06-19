package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCollectionFindByID_Success(t *testing.T) {
	db, err := testutils.CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewCollectionRepository(db)

	created := &model.Collection{Name: "Test Collection"}
	assert.NoError(t, repo.Create(created))

	result, err := repo.FindByID(created.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "Test Collection", result.Name)
}

func TestCollectionFindByID_NotFound(t *testing.T) {
	db, err := testutils.CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewCollectionRepository(db)

	_, err = repo.FindByID(999)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCollectionFindByAlias_Success(t *testing.T) {
	db, err := testutils.CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewCollectionRepository(db)

	created := &model.Collection{Name: "Alias Collection", Alias: "alias123"}
	assert.NoError(t, repo.Create(created))

	result, err := repo.FindByAlias("alias123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "Alias Collection", result.Name)
}

func TestCollectionFindByAlias_NotFound(t *testing.T) {
	db, err := testutils.CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewCollectionRepository(db)

	_, err = repo.FindByAlias("doesnotexist")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
