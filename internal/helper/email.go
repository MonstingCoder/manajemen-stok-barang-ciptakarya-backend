package helper

import (
	"errors"
	"regexp"
)

var (
	ErrEmptyEmail   = errors.New("ERROR EMPTY EMAIL")
	ErrEmailLen     = errors.New("ERROR EMAIL LENGTH")
	ErrInvalidEmail = errors.New("ERROR INVALID EMAIL")

	reEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9](?:[a-zA-Z0-9\-]*[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9\-]*[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)
)

type Email string

func (e *Email) Check() (err error) {
	if *e == "" {
		err = ErrEmptyEmail
		return
	}
	if len([]rune(*e)) >= 256 {
		err = ErrEmailLen
		return
	}
	if !reEmail.MatchString(string(*e)) {
		err = ErrInvalidEmail
	}

	return
}
