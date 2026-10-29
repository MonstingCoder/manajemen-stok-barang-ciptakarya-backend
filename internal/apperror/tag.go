package apperror

import "errors"

var (
	ErrTagExists    = errors.New("ERROR TAG EXISTS")
	ErrTagNameTaken = errors.New("ERROR TAG NAME TAKEN")
	ErrTagNotFound  = errors.New("ERROR TAG NOT FOUND")
)
