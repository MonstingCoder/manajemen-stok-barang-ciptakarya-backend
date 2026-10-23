package model

import "gorm.io/gorm"

type Role struct {
	Name *string `gorm:"not null;unique"`

	Users []User
	gorm.Model
}

type RoleInput struct {
	Name  *string
	Users []UserInput
}
