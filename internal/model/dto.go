package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserLoginReq struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserClaims struct {
	// User *UserOut `json:"user"`
	Role *string
	*jwt.RegisteredClaims
}
