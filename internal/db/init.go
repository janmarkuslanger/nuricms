package db

import (
	"errors"

	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(dl gorm.Dialector) (*gorm.DB, error) {
	if dl == nil {
		return nil, errors.New("no database dialector provided")
	}

	db, err := gorm.Open(dl, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}
