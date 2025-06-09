package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateContent(t *testing.T) {
	gormDB, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	repo := NewContentRepository(gormDB)

	var testContent model.Content
	repo.Create(&testContent)

	assert.Equal(t, testContent.ID, uint(1))
}
