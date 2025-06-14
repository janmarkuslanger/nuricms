package repository

import (
	"fmt"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFieldRepository_FindByCollectionID_WithResults(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewFieldRepository(db)

	col := &model.Collection{Name: "ColA"}
	assert.NoError(t, db.Create(col).Error)

	f1 := &model.Field{Name: "Name1", Alias: "alias1", FieldType: "text", CollectionID: col.ID}
	f2 := &model.Field{Name: "Name2", Alias: "alias2", FieldType: "number", CollectionID: col.ID}
	assert.NoError(t, repo.Create(f1))
	assert.NoError(t, repo.Create(f2))

	fields, err := repo.FindByCollectionID(col.ID)
	assert.NoError(t, err)
	assert.Len(t, fields, 2)
}

func TestFieldRepository_FindByCollectionID_NoResults(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewFieldRepository(db)

	fields, err := repo.FindByCollectionID(999)
	assert.NoError(t, err)
	assert.Empty(t, fields)
}

func TestFieldRepository_FindDisplayFieldsByCollectionID_ReturnsOnlyDisplay(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewFieldRepository(db)

	col := &model.Collection{Name: "ColB"}
	assert.NoError(t, db.Create(col).Error)

	disp := &model.Field{Name: "D", Alias: "d-alias", FieldType: "text", DisplayField: true, CollectionID: col.ID}
	hide := &model.Field{Name: "H", Alias: "h-alias", FieldType: "text", DisplayField: false, CollectionID: col.ID}
	assert.NoError(t, repo.Create(disp))
	assert.NoError(t, repo.Create(hide))

	fields, err := repo.FindDisplayFieldsByCollectionID(col.ID)
	assert.NoError(t, err)
	assert.Len(t, fields, 1)
	assert.Equal(t, disp.ID, fields[0].ID)
}

func TestFieldRepository_FindByID_Found(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewFieldRepository(db)

	col := &model.Collection{Name: "ColC"}
	assert.NoError(t, db.Create(col).Error)
	f := &model.Field{Name: "F", Alias: "f-alias", FieldType: "text", CollectionID: col.ID}
	assert.NoError(t, repo.Create(f))

	got, err := repo.FindByID(f.ID)
	assert.NoError(t, err)
	assert.Equal(t, f.ID, got.ID)
	assert.Equal(t, col.ID, got.CollectionID)
}

func TestFieldRepository_FindByID_NotFound(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewFieldRepository(db)

	_, err := repo.FindByID(12345)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFieldRepository_Save_UpdatesFields(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewFieldRepository(db)

	col := &model.Collection{Name: "ColD"}
	assert.NoError(t, db.Create(col).Error)
	f := &model.Field{Name: "Old", Alias: "old-alias", FieldType: "text", CollectionID: col.ID}
	assert.NoError(t, repo.Create(f))

	f.IsRequired = true
	f.IsList = true
	assert.NoError(t, repo.Save(f))

	updated, err := repo.FindByID(f.ID)
	assert.NoError(t, err)
	assert.True(t, updated.IsRequired)
	assert.True(t, updated.IsList)
}

func TestFieldRepository_List_Pagination(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewFieldRepository(db)

	for ci := 1; ci <= 2; ci++ {
		col := &model.Collection{Name: fmt.Sprintf("Col%d", ci)}
		assert.NoError(t, db.Create(col).Error)
		for fi := 1; fi <= 3; fi++ {
			f := &model.Field{Name: fmt.Sprintf("N%d_%d", ci, fi), Alias: fmt.Sprintf("a%d_%d", ci, fi), FieldType: "text", CollectionID: col.ID}
			assert.NoError(t, repo.Create(f))
		}
	}

	page1, total, err := repo.List(1, 4)
	assert.NoError(t, err)
	assert.Equal(t, int64(6), total)
	assert.Len(t, page1, 4)

	page2, total2, err := repo.List(2, 4)
	assert.NoError(t, err)
	assert.Equal(t, int64(6), total2)
	assert.Len(t, page2, 2)
}

func TestFieldRepository_List_Empty(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewFieldRepository(db)

	list, total, err := repo.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, list)
}

func TestFieldRepository_CreateAndDelete(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewFieldRepository(db)

	col := &model.Collection{Name: "ColE"}
	assert.NoError(t, db.Create(col).Error)
	f := &model.Field{Name: "Temp", Alias: "temp-alias", FieldType: "text", CollectionID: col.ID}
	assert.NoError(t, repo.Create(f))
	assert.NotZero(t, f.ID)

	assert.NoError(t, repo.Delete(f))
	_, err := repo.FindByID(f.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
