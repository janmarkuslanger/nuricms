package model

import (
	"gorm.io/gorm"
)

type Content struct {
	gorm.Model
	CollectionID  uint64         `gorm:"not null"`
	Collection    Collection     `gorm:"foreignKey:CollectionID"`
	ContentValues []ContentValue `gorm:"foreignKey:ContentID;references:ID"`
}
