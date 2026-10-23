package model

import (
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/helper"
	"gorm.io/gorm"
)

type User struct {
	Name         *string `gorm:"not null"`
	Email        *string `gorm:"not null;unique"`
	PasswordHash *string `gorm:"not null"`
	RoleID       *uint   `gorm:"not null"`

	Role *Role
	gorm.Model
}

type UserInput struct {
	Name *string
	*helper.Email
	*helper.Password
	RoleID *uint
}
