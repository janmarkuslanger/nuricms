package model

import (
	"gorm.io/gorm"
)

type Collection struct {
	gorm.Model
	Name        string  `gorm:"size:80;not null"`
	Alias       string  `gorm:"size:80;not null"`
	Description string  `gorm:"size:255"`
	Fields      []Field `gorm:"foreignKey:CollectionID"`
}
