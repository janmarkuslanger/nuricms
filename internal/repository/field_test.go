package repository

import (
	"fmt"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFieldFindByCollectionIDSuccess(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewFieldRepository(db)

	col := &model.Collection{Name: "CollectionA"}
	assert.NoError(t, db.Create(col).Error)

	f1 := &model.Field{Name: "Name1", Alias: "alias1", FieldType: "text", CollectionID: col.ID}
	f2 := &model.Field{Name: "Name2", Alias: "alias2", FieldType: "number", CollectionID: col.ID}
	assert.NoError(t, repo.Create(f1))
	assert.NoError(t, repo.Create(f2))

	fields, err := repo.FindByCollectionID(col.ID)
	assert.NoError(t, err)
	assert.Len(t, fields, 2)
	aliases := []string{fields[0].Alias, fields[1].Alias}
	assert.Contains(t, aliases, "alias1")
	assert.Contains(t, aliases, "alias2")
}

func TestFieldFindByCollectionIDEmpty(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewFieldRepository(db)

	fields, err := repo.FindByCollectionID(999)
	assert.NoError(t, err)
	assert.Empty(t, fields)
}

func TestFieldFindDisplayByCollectionID(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewFieldRepository(db)

	col := &model.Collection{Name: "CollectionB"}
	assert.NoError(t, db.Create(col).Error)

	fd := &model.Field{Name: "Visible", Alias: "vis-alias", FieldType: "text", DisplayField: true, CollectionID: col.ID}
	fh := &model.Field{Name: "Hidden", Alias: "hid-alias", FieldType: "text", DisplayField: false, CollectionID: col.ID}
	assert.NoError(t, repo.Create(fd))
	assert.NoError(t, repo.Create(fh))

	fields, err := repo.FindDisplayFieldsByCollectionID(col.ID)
	assert.NoError(t, err)
	assert.Len(t, fields, 1)
	assert.Equal(t, "vis-alias", fields[0].Alias)
}

func TestFieldFindByIDSuccess(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewFieldRepository(db)

	col := &model.Collection{Name: "CollectionC"}
	assert.NoError(t, db.Create(col).Error)

	f := &model.Field{Name: "FindMe", Alias: "find-alias", FieldType: "text", CollectionID: col.ID}
	assert.NoError(t, repo.Create(f))

	found, err := repo.FindByID(f.ID)
	assert.NoError(t, err)
	assert.Equal(t, "find-alias", found.Alias)
}

func TestFieldFindByIDNotFound(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewFieldRepository(db)

	_, err = repo.FindByID(1234)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFieldSave(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewFieldRepository(db)

	col := &model.Collection{Name: "CollectionD"}
	assert.NoError(t, db.Create(col).Error)

	f := &model.Field{Name: "OldName", Alias: "old-alias", FieldType: "text", CollectionID: col.ID}
	assert.NoError(t, repo.Create(f))

	f.IsRequired = true
	f.IsList = true
	assert.NoError(t, repo.Save(f))

	updated, err := repo.FindByID(f.ID)
	assert.NoError(t, err)
	assert.True(t, updated.IsRequired)
	assert.True(t, updated.IsList)
}

func TestFieldListPagination(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewFieldRepository(db)

	for ci := 1; ci <= 2; ci++ {
		col := &model.Collection{Name: fmt.Sprintf("Col%d", ci)}
		assert.NoError(t, db.Create(col).Error)
		for fi := 1; fi <= 3; fi++ {
			f := &model.Field{Name: fmt.Sprintf("Name%d_%d", ci, fi), Alias: fmt.Sprintf("a%d_%d", ci, fi), FieldType: "text", CollectionID: col.ID}
			assert.NoError(t, repo.Create(f))
		}
	}

	fields, total, err := repo.List(1, 4)
	assert.NoError(t, err)
	assert.Equal(t, int64(6), total)
	assert.Len(t, fields, 4)

	fields2, total2, err := repo.List(2, 4)
	assert.NoError(t, err)
	assert.Equal(t, int64(6), total2)
	assert.Len(t, fields2, 2)
}

func TestFieldListEmpty(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewFieldRepository(db)

	fields, total, err := repo.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, fields)
}

func TestFieldCreate(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewFieldRepository(db)

	col := &model.Collection{Name: "CollectionE"}
	assert.NoError(t, db.Create(col).Error)

	f := &model.Field{Name: "NewName", Alias: "new-alias", FieldType: "text", IsRequired: false, IsList: false, DisplayField: false, CollectionID: col.ID}
	assert.NoError(t, repo.Create(f))
	assert.NotZero(t, f.ID)

	fields, err := repo.FindByCollectionID(col.ID)
	assert.NoError(t, err)
	assert.Equal(t, "new-alias", fields[0].Alias)
}
