package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateContentValue(t *testing.T) {
	gormDB, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	repo := NewContentValueRepository(gormDB)

	var testContentValue model.ContentValue
	repo.Create(&testContentValue)

	var testAnotherContentValue model.ContentValue
	repo.Create(&testAnotherContentValue)

	assert.Equal(t, testContentValue.ID, uint(1))
	assert.Equal(t, testAnotherContentValue.ID, uint(2))
}
