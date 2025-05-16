package model

import "gorm.io/gorm"

type Role string

const (
	RoleAdmin  Role = "Admin"
	RoleEditor Role = "Editor"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     Role   `gorm:"type:varchar(20);not null"`
}

func GetUserRoles() []Role {
	return []Role{RoleEditor, RoleAdmin}
}
