package base

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPreload_QueryOption(t *testing.T) {
	db := testutils.SetupTestDB(t)
	session := db.Session(&gorm.Session{})
	opt := Preload("TestField", 1, "two")
	newDB := opt(session)
	preloads := newDB.Statement.Preloads
	assert.Contains(t, preloads, "TestField")
	assert.Equal(t, []interface{}{1, "two"}, preloads["TestField"])
}
