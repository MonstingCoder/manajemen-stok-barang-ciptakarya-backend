package model

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type Email string

var (
	ErrEmptyEmail   = errors.New("ERROR EMPTY EMAIL")
	ErrEmailLen     = errors.New("ERROR EMAIL LENGTH")
	ErrInvalidEmail = errors.New("ERROR INVALID EMAIL")
	reEmail         = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9](?:[a-zA-Z0-9\-]*[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9\-]*[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)

	ErrPwLen     = errors.New("ERROR PASSWORD LENGTH")
	ErrInvalidPw = errors.New("ERROR INVALID PASSWORD")
	pwUpper      = regexp.MustCompile(`\p{Lu}`)
	pwLower      = regexp.MustCompile(`\p{Ll}`)
	pwDigit      = regexp.MustCompile(`\p{N}`)
	pwSpecial    = regexp.MustCompile(`[\p{P}\p{S}]`)
)

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

type Pw string

func (p *Pw) Check() (err error) {
	if lenP := len(*p); lenP < 8 || lenP >= 72 {
		err = ErrPwLen
		return
	}

	pw := string(*p)
	if !pwUpper.MatchString(pw) || !pwLower.MatchString(pw) || !pwDigit.MatchString(pw) || !pwSpecial.MatchString(pw) {
		err = ErrInvalidPw
	}
	return
}

func (p *Pw) Hash() (pwHash *string, err error) {
	res, err := bcrypt.GenerateFromPassword([]byte(*p), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	pwHash = new(string(res))
	return
}

func (p *Pw) Verify(hash string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(*p))
	return
}

type UserIn struct {
	Name   *string `json:"name"`
	*Email `json:"email"`
	*Pw    `json:"password"`
	RoleID *uint `json:"role_id"`
}

type RoleIn struct {
	Name  *string
	Users []UserIn
}

type ItemIn struct {
	Name        *string
	Price       *uint
	Quantity    *uint
	Description *string
	Tags        []TagIn
}

type TagIn struct {
	ID    uint
	Name  *string
	Items []ItemIn
}
