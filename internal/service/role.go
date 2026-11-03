package service

import (
	"context"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type RoleRepo interface {
	Creates(ctx context.Context, roles []model.Role) (err error)
	Delete(ctx context.Context, id uint) (err error)
}

// var _ RoleRepo = (*repository.RoleRepo)(nil)

type RoleService struct {
	r RoleRepo
}

func NewRoleService(r RoleRepo) (s *RoleService) {
	s = new(RoleService{r})
	return
}

func (s *RoleService) Creates(ctx context.Context, ins []model.RoleInput) (roles []model.Role, err error) {
	roles = make([]model.Role, len(ins))
	for i, in := range ins {
		roles[i] = model.Role{Name: in.Name, Users: make([]model.User, len(in.Users))}
		for j, user := range in.Users {
			if err = user.Email.Check(); err != nil {
				return
			}
			if err = user.Pw.Check(); err != nil {
				return
			}

			var pwHash *string
			if pwHash, err = user.Pw.Hash(); err != nil {
				return
			}
			roles[i].Users[j] = model.User{
				Name:         user.Name,
				Email:        (*string)(user.Email),
				PasswordHash: pwHash,
				RoleID:       user.RoleID,
			}
		}
	}

	err = s.r.Creates(ctx, roles)
	return
}

func (s *RoleService) Delete(ctx context.Context, id uint) (err error) {
	err = s.r.Delete(ctx, id)
	return
}
