package base_test

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"github.com/stretchr/testify/assert"
)

type TestEntity struct {
	gorm.Model
	Name string
}

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	if err := db.AutoMigrate(&TestEntity{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestBaseRepository_Create(t *testing.T) {
	db := setupDB(t)
	repo := base.NewBaseRepository[TestEntity](db)
	e := &TestEntity{Name: "Alice"}
	assert.NoError(t, repo.Create(e))
	assert.NotZero(t, e.ID)
}

func TestBaseRepository_FindByID(t *testing.T) {
	db := setupDB(t)
	repo := base.NewBaseRepository[TestEntity](db)
	e := &TestEntity{Name: "Bob"}
	repo.Create(e)
	got, err := repo.FindByID(e.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", got.Name)
}

func TestBaseRepository_FindByID_WithQueryOption(t *testing.T) {
	db := setupDB(t)
	repo := base.NewBaseRepository[TestEntity](db)
	e := &TestEntity{Name: "QueryOpt"}
	assert.NoError(t, repo.Create(e))

	opt := func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", "QueryOpt")
	}

	got, err := repo.FindByID(e.ID, opt)
	assert.NoError(t, err)
	assert.Equal(t, "QueryOpt", got.Name)
}

func TestBaseRepository_Save(t *testing.T) {
	db := setupDB(t)
	repo := base.NewBaseRepository[TestEntity](db)
	e := &TestEntity{Name: "Carol"}
	repo.Create(e)
	e.Name = "CarolUpdated"
	assert.NoError(t, repo.Save(e))
	got, err := repo.FindByID(e.ID)
	assert.NoError(t, err)
	assert.Equal(t, "CarolUpdated", got.Name)
}

func TestBaseRepository_Delete(t *testing.T) {
	db := setupDB(t)
	repo := base.NewBaseRepository[TestEntity](db)
	e := &TestEntity{Name: "Dave"}
	repo.Create(e)
	assert.NoError(t, repo.Delete(e))
	_, err := repo.FindByID(e.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestBaseRepository_ListPagination(t *testing.T) {
	db := setupDB(t)
	repo := base.NewBaseRepository[TestEntity](db)
	for i := 1; i <= 5; i++ {
		e := &TestEntity{Name: fmt.Sprintf("N%d", i)}
		repo.Create(e)
	}
	page1, total1, err1 := repo.List(1, 2)
	assert.NoError(t, err1)
	assert.Equal(t, int64(5), total1)
	assert.Len(t, page1, 2)
	assert.Equal(t, "N1", page1[0].Name)
	page3, total3, err3 := repo.List(3, 2)
	assert.NoError(t, err3)
	assert.Equal(t, int64(5), total3)
	assert.Len(t, page3, 1)
	assert.Equal(t, "N5", page3[0].Name)
}

func TestBaseRepository_List_CountError(t *testing.T) {
	db := setupDB(t)
	repo := base.NewBaseRepository[TestEntity](db)

	badOption := func(db *gorm.DB) *gorm.DB {
		return db.Where("INVALID SQL SYNTAX ???")
	}

	_, _, err := repo.List(1, 10, badOption)
	assert.Error(t, err)
}
