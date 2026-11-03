package model

import (
	"time"

	"gorm.io/gorm"
)

type UserResp struct {
	ID        uint
	Name      *string         `json:"name"`
	Email     *string         `json:"email"`
	Role      *RoleResponse   `json:"role"`
	CreatedAt *time.Time      `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at"`
}

type UserLoginReq struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type RoleResponse struct {
	ID        uint            `json:"id"`
	Name      *string         `json:"name"`
	Users     []UserResp      `json:"users"`
	CreatedAt *time.Time      `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at"`
}
