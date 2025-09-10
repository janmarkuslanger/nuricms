package model

import (
	"gorm.io/gorm"
)

type FieldOptionType string

const (
	FieldOptionTypeSelectOption FieldType = "SelectOption"
)

type FieldOption struct {
	gorm.Model
	OptionType FieldOptionType `gorm:"type:varchar(20);not null"`
	FieldID    uint            `gorm:"not null"`
	Field      Field           `gorm:"foreignKey:FieldID"`

	Value string `gorm:"type:text"`
}
