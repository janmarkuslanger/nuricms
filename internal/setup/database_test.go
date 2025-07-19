package setup_test

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/setup"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
)

func TestInit_WithSQLiteMemory(t *testing.T) {
	dialector := sqlite.Open(":memory:")

	db, err := setup.InitDatabase(dialector)

	assert.NotNil(t, db)
	assert.NoError(t, err)

	sqlDB, err := db.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())
}

func TestInit_InvalidDriver_Panics(t *testing.T) {
	db, err := setup.InitDatabase(nil)

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestInit_SQLite_InvalidFilePath(t *testing.T) {
	dialector := sqlite.Open("/invalid/path/to/file.db")
	db, err := setup.InitDatabase(dialector)

	assert.Error(t, err)
	assert.Nil(t, db)
}
