package entity

import "gorm.io/gorm"

type Item struct {
	Name        *string `gorm:"not null"`
	Price       *uint   `gorm:"not null"`
	Quantity    *uint   `gorm:"not null"`
	Description *string

	Tags []Tag `gorm:"many2many:item_tags;not null"`
	gorm.Model
}
