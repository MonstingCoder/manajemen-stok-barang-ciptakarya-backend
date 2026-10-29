package helper

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPwLen     = errors.New("ERROR PASSWORD LENGTH")
	ErrInvalidPw = errors.New("ERROR INVALID PASSWORD")

	pwUpper   = regexp.MustCompile(`\p{Lu}`)
	pwLower   = regexp.MustCompile(`\p{Ll}`)
	pwDigit   = regexp.MustCompile(`\p{N}`)
	pwSpecial = regexp.MustCompile(`[\p{P}\p{S}]`)
)

type Password string

func (p *Password) Check() (err error) {
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

func (p *Password) Hash() (pwHash *string, err error) {
	res, err := bcrypt.GenerateFromPassword([]byte(*p), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	pwHash = new(string(res))
	return
}
