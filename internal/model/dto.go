package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserLoginReq struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserProfile struct {
	ID     uint             `json:"id"`
	Name   *string          `json:"name"`
	Email  *string          `json:"email"`
	RoleID *uint            `json:"role_id"`
	Role   *UserProfileRole `json:"role"`
}

type UserProfileRole struct {
	ID   uint    `json:"id"`
	Name *string `json:"name"`
}

type UserClaims struct {
	UserID uint
	Role   *string
	*jwt.RegisteredClaims
}
