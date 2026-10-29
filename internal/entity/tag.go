package entity

import "gorm.io/gorm"

type Tag struct {
	Name *string `gorm:"not null;unique"`

	Items []Item `gorm:"many2many:item_tags"`
	gorm.Model
}
