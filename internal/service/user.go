package service

import (
	"context"
	"errors"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/entity"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/helper"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

var ErrLimitTooLarge = errors.New("ERROR LIMIT TO LARGE")

type UserRepo interface {
	Create(ctx context.Context, user *entity.User) (err error)
	Finds(ctx context.Context, limit, offset uint) (users []entity.User, err error)
	FindByEmail(ctx context.Context, email string) (user *entity.User, err error)
	Update(ctx context.Context, id uint, user *entity.User) (err error)
	Delete(ctx context.Context, id uint) (err error)
}

// var _ UserRepo = (*repository.UserRepo)(nil)

type UserService struct {
	r UserRepo
}

func NewUserService(r UserRepo) (s *UserService) {
	s = new(UserService{r})
	return
}

func (s *UserService) Create(ctx context.Context, in *model.UserInput) (user *entity.User, err error) {
	if err = in.Email.Check(); err != nil {
		return
	}
	if err = in.Password.Check(); err != nil {
		return
	}

	pwHash, err := in.Password.Hash()
	if err != nil {
		return
	}

	user = new(entity.User{
		Name:         in.Name,
		Email:        (*string)(in.Email),
		PasswordHash: pwHash,
		RoleID:       in.RoleID,
	})
	err = s.r.Create(ctx, user)
	return
}

func (s *UserService) Finds(ctx context.Context, limit, offset uint) (users []entity.User, err error) {
	if limit == 0 {
		limit = 5
	}
	if limit > 50 {
		err = ErrLimitTooLarge
		return
	}

	users, err = s.r.Finds(ctx, limit, offset)
	return
}

func (s *UserService) FindByEmail(ctx context.Context, email helper.Email) (user *entity.User, err error) {
	if err = email.Check(); err != nil {
		return
	}

	user, err = s.r.FindByEmail(ctx, string(email))
	return
}

func (s *UserService) Update(ctx context.Context, id uint, in *model.UserInput) (err error) {
	if err = in.Email.Check(); err != nil {
		return
	}
	if err = in.Password.Check(); err != nil {
		return
	}

	pwHash, err := in.Password.Hash()
	if err != nil {
		return
	}

	err = s.r.Update(ctx, id, new(entity.User{
		Name:         in.Name,
		Email:        (*string)(in.Email),
		PasswordHash: pwHash,
		RoleID:       in.RoleID,
	}))
	return
}

func (s *UserService) Delete(ctx context.Context, id uint) (err error) {
	err = s.r.Delete(ctx, id)
	return
}
