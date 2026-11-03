package model

import "gorm.io/gorm"

type User struct {
	Name         *string `gorm:"not null" json:"name"`
	Email        *string `gorm:"not null;unique" json:"email"`
	PasswordHash *string `gorm:"not null" json:"-"`
	RoleID       *uint   `gorm:"not null" json:"role_id"`
	// ImagePath

	Role *Role `json:"role"`
	gorm.Model
}

type Role struct {
	Name *string `gorm:"not null;unique"`

	Users []User
	gorm.Model
}

type Item struct {
	Name        *string `gorm:"not null"`
	Price       *uint   `gorm:"not null"`
	Quantity    *uint   `gorm:"not null"`
	Description *string
	// image

	Tags []Tag `gorm:"many2many:item_tags;not null"`
	gorm.Model
}

type Tag struct {
	Name *string `gorm:"not null;unique"`

	Items []Item `gorm:"many2many:item_tags"`
	gorm.Model
}
