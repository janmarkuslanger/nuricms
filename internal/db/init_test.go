package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
)

func TestInit_WithSQLiteMemory(t *testing.T) {
	dialector := sqlite.Open(":memory:")

	db, err := Init(dialector)

	assert.NotNil(t, db)
	assert.Equal(t, db, DB)
	assert.NoError(t, err)

	sqlDB, err := db.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())
}

func TestInit_InvalidDriver_Panics(t *testing.T) {
	db, err := Init(nil)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestInit_SQLite_InvalidFilePath(t *testing.T) {
	dialector := sqlite.Open("/invalid/path/to/file.db")
	db, err := Init(dialector)

	assert.Error(t, err)
	assert.Nil(t, db)
}
