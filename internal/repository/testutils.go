package repository

import (
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
