package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/testutils"
)

func TestContentService_CreateAndFindByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewContentService(repos)

	col := &model.Collection{Name: "Col1", Alias: "col1"}
	repos.Collection.Create(col)

	c := &model.Content{CollectionID: col.ID}
	created, err := s.Create(c)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	found, err := s.FindByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
}

func TestContentService_ListByCollectionAliasAndByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewContentService(repos)

	col := &model.Collection{Name: "Col2", Alias: "col2"}
	repos.Collection.Create(col)

	for i := 0; i < 3; i++ {
		repos.Content.Create(&model.Content{CollectionID: col.ID})
	}

	listByID, err := s.FindByCollectionID(col.ID)
	assert.NoError(t, err)
	assert.Len(t, listByID, 3)

	listByAlias, err := s.ListByCollectionAlias("col2", 0, 0)
	assert.NoError(t, err)
	assert.Len(t, listByAlias, 3)

	_, err = s.ListByCollectionAlias("missing", 0, 0)
	assert.Error(t, err)
}

func TestContentService_FindDisplayAndListWithDisplay(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewContentService(repos)

	col := &model.Collection{Name: "Col3", Alias: "col3"}
	repos.Collection.Create(col)
	f := &model.Field{Name: "F", Alias: "f", FieldType: "text", CollectionID: col.ID, DisplayField: true}
	repos.Field.Create(f)

	repos.Field.Create(&model.Field{Name: "X", Alias: "x", FieldType: "text", CollectionID: col.ID})

	for i := 0; i < 2; i++ {
		c, _ := s.Create(&model.Content{CollectionID: col.ID})
		repos.ContentValue.Create(&model.ContentValue{ContentID: c.ID, FieldID: f.ID, Value: "v"})
		repos.ContentValue.Create(&model.ContentValue{ContentID: c.ID, FieldID: f.ID, Value: "w"})
	}

	disp, total, err := s.FindDisplayValueByCollectionID(col.ID, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, disp, 1)
	for _, c := range disp {
		assert.Len(t, c.ContentValues, 2)
	}

	allDisp, err := s.FindContentsWithDisplayContentValue()
	assert.NoError(t, err)
	assert.Len(t, allDisp, 2)
}
