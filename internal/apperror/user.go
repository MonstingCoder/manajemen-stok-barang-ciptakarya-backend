package apperror

import "errors"

var (
	ErrUserEmailTaken = errors.New("ERROR EMAIL TAKEN")
	ErrUserExists     = errors.New("ERROR USER ALREADY EXISTS")
	ErrUserNotFound   = errors.New("ERROR USER NOT FOUND")
)
