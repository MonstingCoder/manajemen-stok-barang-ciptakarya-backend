package model

import (
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/helper"
)

type UserInput struct {
	Name *string
	*helper.Email
	*helper.Password
	RoleID *uint
}
