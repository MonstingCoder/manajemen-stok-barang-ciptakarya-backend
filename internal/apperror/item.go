package apperror

import "errors"

var (
	ErrItemExists   = errors.New("ERROR ITEM EXISTS")
	ErrItemNotFound = errors.New("ERROR ITEM NOT FOUND")
)
