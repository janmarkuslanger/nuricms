package repository

import (
	"testing"

	"fmt"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCollectionFindByID_Success(t *testing.T) {
	db, err := CreateTestDB()
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
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewCollectionRepository(db)

	_, err = repo.FindByID(999)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCollectionFindByAlias_Success(t *testing.T) {
	db, err := CreateTestDB()
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
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewCollectionRepository(db)

	_, err = repo.FindByAlias("doesnotexist")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCollectionCreate_And_Save(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewCollectionRepository(db)

	col := &model.Collection{Name: "Create Test", Alias: "create-alias"}
	assert.NoError(t, repo.Create(col))
	assert.NotZero(t, col.ID)

	col.Name = "Updated Name"
	assert.NoError(t, repo.Save(col))

	fetched, err := repo.FindByID(col.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", fetched.Name)
}

func TestCollectionDelete(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewCollectionRepository(db)

	col := &model.Collection{Name: "To Be Deleted"}
	assert.NoError(t, repo.Create(col))

	assert.NoError(t, repo.Delete(col))

	_, err = repo.FindByID(col.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCollectionList_Pagination(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewCollectionRepository(db)

	for i := 1; i <= 5; i++ {
		assert.NoError(t, repo.Create(&model.Collection{Name: fmt.Sprintf("Item %d", i)}))
	}

	page1, total, err := repo.List(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, page1, 2)
	assert.Equal(t, "Item 1", page1[0].Name)

	page3, total, err := repo.List(3, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, page3, 1)
	assert.Equal(t, "Item 5", page3[0].Name)
}

func TestCollectionList_Empty(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewCollectionRepository(db)

	list, total, err := repo.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, list)
}
