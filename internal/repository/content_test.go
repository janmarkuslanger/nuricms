package repository

import (
	"fmt"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestContentRepository_DeleteByID(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewContentRepository(db)
	col := &model.Collection{Name: "col"}
	db.Create(col)
	c := &model.Content{CollectionID: col.ID}
	repo.Create(c)
	assert.NoError(t, repo.DeleteByID(c.ID))
	_, err := repo.FindByID(c.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestContentRepository_FindByCollectionID(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewContentRepository(db)
	col := &model.Collection{Name: "col"}
	db.Create(col)
	f := &model.Field{Name: "Field", Alias: "field", FieldType: "text", CollectionID: col.ID}
	db.Create(f)
	for i := 0; i < 3; i++ {
		c := &model.Content{CollectionID: col.ID}
		repo.Create(c)
		cv := &model.ContentValue{ContentID: c.ID, FieldID: f.ID, Value: fmt.Sprintf("v%d", i)}
		db.Create(cv)
	}
	list, err := repo.FindByCollectionID(col.ID, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Len(t, list[0].ContentValues, 1)
	assert.Equal(t, f.ID, list[0].ContentValues[0].Field.ID)
}

func TestContentRepository_FindDisplayValueByCollectionID(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewContentRepository(db)
	col := &model.Collection{Name: "col"}
	db.Create(col)
	fShow := &model.Field{Name: "Show", Alias: "show", FieldType: "text", CollectionID: col.ID, DisplayField: true}
	fHide := &model.Field{Name: "Hide", Alias: "hide", FieldType: "text", CollectionID: col.ID, DisplayField: false}
	db.Create(fShow)
	db.Create(fHide)
	for i := 0; i < 3; i++ {
		c := &model.Content{CollectionID: col.ID}
		repo.Create(c)
		db.Create(&model.ContentValue{ContentID: c.ID, FieldID: fShow.ID, Value: "x"})
		db.Create(&model.ContentValue{ContentID: c.ID, FieldID: fHide.ID, Value: "y"})
	}

	list1, total1, err1 := repo.FindDisplayValueByCollectionID(col.ID, 1, 2)
	assert.NoError(t, err1)
	assert.Equal(t, int64(3), total1)
	assert.Len(t, list1, 2)
	for _, c := range list1 {
		assert.Len(t, c.ContentValues, 1)
		assert.Equal(t, fShow.ID, c.ContentValues[0].Field.ID)
	}

	list2, total2, err2 := repo.FindDisplayValueByCollectionID(col.ID, 2, 2)
	assert.NoError(t, err2)
	assert.Equal(t, int64(3), total2)
	assert.Len(t, list2, 1)
	assert.Len(t, list2[0].ContentValues, 1)
	assert.Equal(t, fShow.ID, list2[0].ContentValues[0].Field.ID)
}

func TestContentRepository_ListWithDisplayContentValue(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewContentRepository(db)
	col := &model.Collection{Name: "col"}
	db.Create(col)
	fShow := &model.Field{Name: "Show", Alias: "show", FieldType: "text", CollectionID: col.ID, DisplayField: true}
	fHide := &model.Field{Name: "Hide", Alias: "hide", FieldType: "text", CollectionID: col.ID, DisplayField: false}
	db.Create(fShow)
	db.Create(fHide)
	for i := 0; i < 3; i++ {
		c := &model.Content{CollectionID: col.ID}
		repo.Create(c)
		db.Create(&model.ContentValue{ContentID: c.ID, FieldID: fShow.ID, Value: "x"})
		db.Create(&model.ContentValue{ContentID: c.ID, FieldID: fHide.ID, Value: "y"})
	}
	list, err := repo.ListWithDisplayContentValue()
	assert.NoError(t, err)
	assert.Len(t, list, 3)
	for _, c := range list {
		assert.Len(t, c.ContentValues, 1)
		assert.Equal(t, fShow.ID, c.ContentValues[0].Field.ID)
		assert.Equal(t, col.ID, c.Collection.ID)
	}
}

func TestContentRepository_FindByCollectionAndFieldValue(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewContentRepository(db)
	col := &model.Collection{Name: "col"}
	db.Create(col)
	f := &model.Field{Name: "F", Alias: "alias", FieldType: "text", CollectionID: col.ID}
	db.Create(f)
	c1 := &model.Content{CollectionID: col.ID}
	repo.Create(c1)
	db.Create(&model.ContentValue{ContentID: c1.ID, FieldID: f.ID, Value: "match"})
	c2 := &model.Content{CollectionID: col.ID}
	repo.Create(c2)
	db.Create(&model.ContentValue{ContentID: c2.ID, FieldID: f.ID, Value: "other"})
	list, total, err := repo.FindByCollectionAndFieldValue(col.ID, "alias", "match", 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, list, 1)
	assert.Equal(t, c1.ID, list[0].ID)
}
