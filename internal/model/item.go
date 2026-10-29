package model

type ItemInput struct {
	Name        *string
	Price       *uint
	Quantity    *uint
	Description *string
	Tags        []TagInput
}
