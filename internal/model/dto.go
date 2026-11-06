package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserLoginReq struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserClaims struct {
	UserID uint
	Role   *string
	*jwt.RegisteredClaims
}
