package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	var err error
	db, err := gorm.Open(sqlite.Open("nuricms.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	DB = db
	return db
}
