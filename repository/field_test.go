package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func createFields(gormDB *gorm.DB) {
	var testCollection *model.Collection
	testCollection = &model.Collection{Name: "Cars"}

	gormDB.Create(&testCollection)

	gormDB.Create(&model.Field{
		Name:         "Name",
		Alias:        "name",
		FieldType:    model.FieldTypeText,
		CollectionID: testCollection.ID,
		IsList:       false,
		IsRequired:   false,
		DisplayField: false,
	})

	gormDB.Create(&model.Field{
		Name:         "Nice",
		Alias:        "nice",
		FieldType:    model.FieldTypeText,
		CollectionID: testCollection.ID,
		IsList:       false,
		IsRequired:   false,
		DisplayField: true,
	})

	gormDB.Create(&model.Field{
		Name:         "Foo",
		Alias:        "bar",
		FieldType:    model.FieldTypeText,
		CollectionID: testCollection.ID,
		IsList:       false,
		IsRequired:   false,
		DisplayField: false,
	})
}

func TestFindByCollectionID(t *testing.T) {
	gormDB, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	repo := NewFieldRepository(gormDB)

	createFields(gormDB)

	fields, err := repo.FindByCollectionID(uint(1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Len(t, fields, 3)
}

func TestFindDisplayFieldsByCollectionID(t *testing.T) {
	gormDB, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	repo := NewFieldRepository(gormDB)
	createFields(gormDB)

	fields, err := repo.FindDisplayFieldsByCollectionID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Len(t, fields, 1)
	assert.Equal(t, fields[0].Name, "Nice")
}
