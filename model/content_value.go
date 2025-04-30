package model

import "gorm.io/gorm"

type ContentValue struct {
	gorm.Model

	SortIndex int `gorm:"not null"`

	ContentID uint
	Content   Content `gorm:"constraint:OnDelete:CASCADE"`

	FieldID uint
	Field   Field `gorm:"foreignKey:FieldID;references:ID"`

	Value string `gorm:"type:text"`
}
