package model

import (
	"gorm.io/gorm"
)

type Asset struct {
	gorm.Model
	Name string `gorm:"size:80;not null"`
	Path string `gorm:"size:255;notnull"`
}
