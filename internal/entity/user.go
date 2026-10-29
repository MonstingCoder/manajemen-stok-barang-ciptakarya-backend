package entity

import "gorm.io/gorm"

type User struct {
	Name         *string `gorm:"not null"`
	Email        *string `gorm:"not null;unique"`
	PasswordHash *string `gorm:"not null"`
	RoleID       *uint   `gorm:"not null"`

	Role *Role
	gorm.Model
}
