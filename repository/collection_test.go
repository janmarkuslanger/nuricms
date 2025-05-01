package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/model"
	"github.com/stretchr/testify/assert"
)

func TestFindByID(t *testing.T) {
	gormDB, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	repo := NewCollectionRepository(gormDB)

	gormDB.Create(&model.Collection{Name: "Test Collection"})

	collection, err := repo.FindByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.NotNil(t, collection)
	assert.Equal(t, uint(1), collection.ID)
	assert.Equal(t, "Test Collection", collection.Name)
}

func TestGetAll(t *testing.T) {
	gormDB, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	repo := NewCollectionRepository(gormDB)

	gormDB.Create(&model.Collection{Name: "Collection 1"})
	gormDB.Create(&model.Collection{Name: "Collection 2"})

	collections, err := repo.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Len(t, collections, 2)
	assert.Equal(t, uint(1), collections[0].ID)
	assert.Equal(t, "Collection 1", collections[0].Name)
	assert.Equal(t, uint(2), collections[1].ID)
	assert.Equal(t, "Collection 2", collections[1].Name)
}
