package apperror

import "errors"

var (
	ErrRoleExists   = errors.New("ERROR ROLE EXISTS")
	ErrRoleNotFound = errors.New("ERROR ROLE NOT FOUND")
)
