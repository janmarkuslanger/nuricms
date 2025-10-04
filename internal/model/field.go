package model

import (
	"gorm.io/gorm"
)

type FieldType string

const (
	FieldTypeText        FieldType = "Text"
	FieldTypeNumber      FieldType = "Number"
	FieldTypeBoolean     FieldType = "Boolean"
	FieldTypeDate        FieldType = "Date"
	FieldTypeAsset       FieldType = "Asset"
	FieldTypeCollection  FieldType = "Collection"
	FieldTypeTextarea    FieldType = "Textarea"
	FieldTypeRichText    FieldType = "RichText"
	FieldTypeMultiSelect FieldType = "MultiSelect"
)

type Field struct {
	gorm.Model
	Name         string     `gorm:"size:80;not null"`
	Alias        string     `gorm:"size:80;not null"`
	FieldType    FieldType  `gorm:"type:varchar(20);not null"`
	CollectionID uint       `gorm:"not null"`
	Collection   Collection `gorm:"foreignKey:CollectionID"`
	IsList       bool       `gorm:"not null;default:false"`
	IsRequired   bool       `gorm:"not null;default:false"`
	DisplayField bool       `gorm:"not null;default:false"`
}
