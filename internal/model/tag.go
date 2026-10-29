package model

type TagInput struct {
	ID    uint
	Name  *string
	Items []ItemInput
}
