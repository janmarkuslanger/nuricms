package setup

import (
	"errors"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"gorm.io/gorm"
)

func InitDatabase(dl gorm.Dialector) (*gorm.DB, error) {
	if dl == nil {
		return nil, errors.New("no database dialector provided")
	}

	db, err := gorm.Open(dl, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(
		&model.Collection{},
		&model.Field{},
		&model.FieldOption{},
		&model.Content{},
		&model.ContentValue{},
		&model.Asset{},
		&model.User{},
		&model.Apikey{},
		&model.Webhook{},
	)
}
