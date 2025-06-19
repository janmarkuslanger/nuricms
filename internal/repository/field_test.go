package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
)

func TestFieldRepository_FindByCollectionID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewFieldRepository(db)
	col := &model.Collection{Name: "C1"}
	assert.NoError(t, db.Create(col).Error)
	f1 := &model.Field{Name: "F1", Alias: "a1", FieldType: "text", CollectionID: col.ID}
	f2 := &model.Field{Name: "F2", Alias: "a2", FieldType: "text", CollectionID: col.ID}
	assert.NoError(t, repo.Create(f1))
	assert.NoError(t, repo.Create(f2))
	list, err := repo.FindByCollectionID(col.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	aliases := []string{list[0].Alias, list[1].Alias}
	assert.Contains(t, aliases, "a1")
	assert.Contains(t, aliases, "a2")
}

func TestFieldRepository_FindDisplayFieldsByCollectionID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewFieldRepository(db)
	col := &model.Collection{Name: "C2"}
	assert.NoError(t, db.Create(col).Error)
	fd := &model.Field{Name: "Show", Alias: "show", FieldType: "text", DisplayField: true, CollectionID: col.ID}
	fh := &model.Field{Name: "Hide", Alias: "hide", FieldType: "text", DisplayField: false, CollectionID: col.ID}
	assert.NoError(t, repo.Create(fd))
	assert.NoError(t, repo.Create(fh))
	list, err := repo.FindDisplayFieldsByCollectionID(col.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "show", list[0].Alias)
}
