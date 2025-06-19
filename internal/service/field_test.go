package service

import (
	"fmt"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFieldService_FindByCollectionID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	f1 := &model.Field{Name: "FieldA", Alias: "aliasA", FieldType: "text", CollectionID: col.ID}
	f2 := &model.Field{Name: "FieldB", Alias: "aliasB", FieldType: "text", CollectionID: col.ID}
	repos.Field.Create(f1)
	repos.Field.Create(f2)
	list, err := s.FindByCollectionID(col.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestFieldService_FindDisplayFieldsByCollectionID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	fd := &model.Field{Name: "DisplayField", Alias: "display", FieldType: "text", CollectionID: col.ID, DisplayField: true}
	fh := &model.Field{Name: "HiddenField", Alias: "hidden", FieldType: "text", CollectionID: col.ID, DisplayField: false}
	repos.Field.Create(fd)
	repos.Field.Create(fh)
	list, err := s.FindDisplayFieldsByCollectionID(col.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, fd.ID, list[0].ID)
}

func TestFieldService_FindByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	f := &model.Field{Name: "SomeField", Alias: "some", FieldType: "text", CollectionID: col.ID}
	repos.Field.Create(f)
	got, err := s.FindByID(f.ID)
	assert.NoError(t, err)
	assert.Equal(t, f.ID, got.ID)
}

func TestFieldService_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	for i := 0; i < 5; i++ {
		repos.Field.Create(&model.Field{Name: "SameName", Alias: "samealias", FieldType: "text", CollectionID: col.ID})
	}
	page1, total, err := s.List(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, page1, 2)
	page3, total3, err := s.List(3, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total3)
	assert.Len(t, page3, 1)
}

func TestFieldService_Create_InvalidCollectionID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	data := dto.FieldData{CollectionID: "invalid", Name: "Name", Alias: "alias", FieldType: "text"}
	_, err := s.Create(data)
	assert.EqualError(t, err, "cannot convert collection id")
}

func TestFieldService_Create_NoName(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	data := dto.FieldData{CollectionID: fmt.Sprint(col.ID), Name: "", Alias: "alias", FieldType: "text"}
	_, err := s.Create(data)
	assert.EqualError(t, err, "no name given")
}

func TestFieldService_Create_NoAlias(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	data := dto.FieldData{CollectionID: fmt.Sprint(col.ID), Name: "Name", Alias: "", FieldType: "text"}
	_, err := s.Create(data)
	assert.EqualError(t, err, "no alias given")
}

func TestFieldService_Create_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	data := dto.FieldData{
		CollectionID: fmt.Sprint(col.ID),
		Name:         "FieldName",
		Alias:        "fieldalias",
		FieldType:    "text",
		IsList:       "on",
		IsRequired:   "on",
		DisplayField: "on",
	}
	field, err := s.Create(data)
	assert.NoError(t, err)
	assert.Equal(t, data.Name, field.Name)
	assert.Equal(t, data.Alias, field.Alias)
	assert.Equal(t, col.ID, field.CollectionID)
	assert.True(t, field.IsList)
	assert.True(t, field.IsRequired)
	assert.True(t, field.DisplayField)
}

func TestFieldService_UpdateByID_NotFound(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	_, err := s.UpdateByID(999, dto.FieldData{CollectionID: "1", Name: "Name", Alias: "alias", FieldType: "text"})
	assert.Error(t, err)
}

func TestFieldService_UpdateByID_InvalidCollectionID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	f := &model.Field{Name: "OldField", Alias: "oldalias", FieldType: "text", CollectionID: col.ID}
	repos.Field.Create(f)
	data := dto.FieldData{CollectionID: "bad", Name: "NewName", Alias: "newalias", FieldType: "text"}
	_, err := s.UpdateByID(f.ID, data)
	assert.EqualError(t, err, "cannot convert collection id")
}

func TestFieldService_UpdateByID_NoName(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	f := &model.Field{Name: "OldField", Alias: "oldalias", FieldType: "text", CollectionID: col.ID}
	repos.Field.Create(f)
	data := dto.FieldData{CollectionID: fmt.Sprint(col.ID), Name: "", Alias: "newalias", FieldType: "text"}
	_, err := s.UpdateByID(f.ID, data)
	assert.EqualError(t, err, "no name given")
}

func TestFieldService_UpdateByID_NoAlias(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	f := &model.Field{Name: "OldField", Alias: "oldalias", FieldType: "text", CollectionID: col.ID}
	repos.Field.Create(f)
	data := dto.FieldData{CollectionID: fmt.Sprint(col.ID), Name: "NewName", Alias: "", FieldType: "text"}
	_, err := s.UpdateByID(f.ID, data)
	assert.EqualError(t, err, "no alias given")
}

func TestFieldService_UpdateByID_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	f := &model.Field{Name: "OldField", Alias: "oldalias", FieldType: "text", CollectionID: col.ID, IsList: false, IsRequired: false, DisplayField: false}
	repos.Field.Create(f)
	data := dto.FieldData{
		CollectionID: fmt.Sprint(col.ID),
		Name:         "NewName",
		Alias:        "newalias",
		FieldType:    "text",
		IsList:       "on",
		IsRequired:   "",
		DisplayField: "on",
	}
	updated, err := s.UpdateByID(f.ID, data)
	assert.NoError(t, err)
	assert.Equal(t, "NewName", updated.Name)
	assert.Equal(t, "newalias", updated.Alias)
	assert.Equal(t, col.ID, updated.CollectionID)
	assert.True(t, updated.IsList)
	assert.False(t, updated.IsRequired)
	assert.True(t, updated.DisplayField)
}

func TestFieldService_DeleteByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	s := NewFieldService(repos)
	col := &model.Collection{Name: "TestCollection", Alias: "test"}
	repos.Collection.Create(col)
	f := &model.Field{Name: "FieldToDelete", Alias: "todelete", FieldType: "text", CollectionID: col.ID}
	repos.Field.Create(f)
	err := s.DeleteByID(f.ID)
	assert.NoError(t, err)
	_, err = repos.Field.FindByID(f.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
