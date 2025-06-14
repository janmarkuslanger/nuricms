package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestContentValueRepository_Create_Success(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewContentValueRepository(db)

	content := &model.Content{}
	assert.NoError(t, db.Create(content).Error)
	field := &model.Field{Name: "F", Alias: "f", FieldType: "text", CollectionID: 1}
	assert.NoError(t, db.Create(field).Error)

	cv := &model.ContentValue{ContentID: content.ID, FieldID: field.ID, Value: "initial"}
	err := repo.Create(cv)
	assert.NoError(t, err)
	assert.NotZero(t, cv.ID)
}

func TestContentValueRepository_Save_UpdatesValue(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewContentValueRepository(db)

	content := &model.Content{}
	assert.NoError(t, db.Create(content).Error)
	field := &model.Field{Name: "F2", Alias: "f2", FieldType: "text", CollectionID: 1}
	assert.NoError(t, db.Create(field).Error)

	cv := &model.ContentValue{ContentID: content.ID, FieldID: field.ID, Value: "old"}
	assert.NoError(t, repo.Create(cv))

	cv.Value = "updated"
	err := repo.Save(cv)
	assert.NoError(t, err)

	var fetched model.ContentValue
	assert.NoError(t, db.First(&fetched, cv.ID).Error)
	assert.Equal(t, "updated", fetched.Value)
}

func TestContentValueRepository_Delete_RemovesRecord(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewContentValueRepository(db)

	content := &model.Content{}
	assert.NoError(t, db.Create(content).Error)
	field := &model.Field{Name: "F3", Alias: "f3", FieldType: "text", CollectionID: 1}
	assert.NoError(t, db.Create(field).Error)

	cv := &model.ContentValue{ContentID: content.ID, FieldID: field.ID, Value: "toDelete"}
	assert.NoError(t, repo.Create(cv))

	err := repo.Delete(cv)
	assert.NoError(t, err)

	_, err = db.First(&model.ContentValue{}, cv.ID).Error, gorm.ErrRecordNotFound
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
