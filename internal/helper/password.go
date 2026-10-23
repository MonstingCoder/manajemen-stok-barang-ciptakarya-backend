package helper

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPwLen     = errors.New("ERROR PASSWORD LENGTH")
	ErrInvalidPw = errors.New("ERROR INVALID PASSWORD")
)

type Password string

func (p *Password) Check() (err error) {
	if lenP := len(*p); lenP < 8 || lenP >= 72 {
		err = ErrPwLen
		return
	}
	if !strings.ContainsFunc(string(*p), func(r rune) (res bool) {
		res = unicode.IsUpper(r) &&
			unicode.IsLower(r) &&
			unicode.IsDigit(r) &&
			unicode.IsSpace(r) &&
			unicode.IsPunct(r) &&
			unicode.IsSymbol(r)
		return
	}) {
		err = ErrInvalidPw
	}

	return
}

func (p *Password) Hash() (pwHash *string, err error) {
	res, err := bcrypt.GenerateFromPassword([]byte(*p), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	pwHash = new(string(res))
	return
}
