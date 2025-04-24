package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Content struct {
	gorm.Model
	CollectionID uint           `gorm:"not null"`
	Collection   Collection     `gorm:"foreignKey:CollectionID"`
	Data         datatypes.JSON `gorm:"type:json;not null"`
}
