package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
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
