package entity

import "gorm.io/gorm"

type Role struct {
	Name *string `gorm:"not null;unique"`

	Users []User
	gorm.Model
}
