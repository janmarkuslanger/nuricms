package model

import "gorm.io/gorm"

type ContentValue struct {
	gorm.Model
	ContentID uint
	Content   Content `gorm:"constraint:OnDelete:CASCADE"`

	FieldID uint
	Field   Field

	Value string `gorm:"type:text"`
}
