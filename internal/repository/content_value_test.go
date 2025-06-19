package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
)

func TestContentValueRepository_FindByContentID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewContentValueRepository(db)
	col := &model.Collection{Name: "col"}
	db.Create(col)
	f := &model.Field{Name: "F", Alias: "f", FieldType: "text", CollectionID: col.ID}
	db.Create(f)
	c1 := &model.Content{CollectionID: col.ID}
	db.Create(c1)
	c2 := &model.Content{CollectionID: col.ID}
	db.Create(c2)
	for i := 0; i < 3; i++ {
		cv := &model.ContentValue{ContentID: c1.ID, FieldID: f.ID, Value: "v1"}
		repo.Create(cv)
	}
	for i := 0; i < 2; i++ {
		cv := &model.ContentValue{ContentID: c2.ID, FieldID: f.ID, Value: "v2"}
		repo.Create(cv)
	}
	list1, err1 := repo.FindByContentID(c1.ID)
	assert.NoError(t, err1)
	assert.Len(t, list1, 3)
	list2, err2 := repo.FindByContentID(c2.ID)
	assert.NoError(t, err2)
	assert.Len(t, list2, 2)
}
