package base

import "gorm.io/gorm"

func Preload(field string, args ...interface{}) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(field, args...)
	}
}
