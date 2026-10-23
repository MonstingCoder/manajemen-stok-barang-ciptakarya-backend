package model

import "gorm.io/gorm"

type Item struct {
	Name     *string `gorm:"not null"`
	Price    *uint   `gorm:"not null"`
	Quantity *uint   `gorm:"not null"`

	Tags []Tag `gorm:"many2many:item_tags"`
	gorm.Model
}
