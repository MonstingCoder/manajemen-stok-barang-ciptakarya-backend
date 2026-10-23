package service

import (
	"context"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type roleRepo interface {
	Creates(ctx context.Context, roles []model.Role) (err error)
	Delete(ctx context.Context, id uint) (err error)
}

// var _ roleRepo = (*repository.RoleRepo)(nil)

type RoleService struct {
	r roleRepo
}

func NewRoleService(r roleRepo) (s *RoleService) {
	s = new(RoleService{r})
	return
}

func (s *RoleService) Creates(ctx context.Context, in model.RoleInput) (role model.Role, err error) {
	var pwHashes []*string
	for _, user := range in.Users {
		if err = user.Email.Check(); err != nil {
			return
		}
		if err = user.Password.Check(); err != nil {
			return
		}

		var pwHash *string
		pwHash, err = user.Password.Hash()
		if err != nil {
			return
		}

		pwHashes = append(pwHashes, pwHash)
		// users :=
		/*
			new(model.User{
			Name:         in.Name,
			Email:        (*string)(in.Email),
			PasswordHash: pwHash,
			RoleID:       in.RoleID,
			})
		*/
	}
	panic("unimplemented")
}
