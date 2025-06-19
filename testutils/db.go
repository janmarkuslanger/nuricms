package testutils

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateTestDB() (*gorm.DB, error) {
	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN: ":memory:",
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = gormDB.AutoMigrate(
		&model.Collection{},
		&model.Field{},
		&model.Content{},
		&model.ContentValue{},
		&model.Asset{},
		&model.Webhook{},
		&model.Apikey{},
		&model.User{},
	)
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create test DB: %v", err)
	}

	t.Cleanup(func() {
		models := []interface{}{
			&model.ContentValue{},
			&model.Content{},
			&model.Asset{},
			&model.Webhook{},
			&model.Apikey{},
			&model.Field{},
			&model.Collection{},
			&model.User{},
		}

		for _, m := range models {
			if err := db.
				Session(&gorm.Session{AllowGlobalUpdate: true}).
				Delete(m).Error; err != nil {
				t.Fatalf("failed to clear table %T: %v", m, err)
			}
		}
	})

	return db
}
