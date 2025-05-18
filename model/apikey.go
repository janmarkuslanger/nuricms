// model/api_key.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type Apikey struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt *time.Time
}
