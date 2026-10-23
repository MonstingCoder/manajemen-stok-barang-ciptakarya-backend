package service

import (
	"context"
	"errors"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/helper"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

var (
	ErrLimitTooLarge = errors.New("ERROR LIMIT TO LARGE")
)

type userRepo interface {
	Create(ctx context.Context, user *model.User) (err error)
	Finds(ctx context.Context, offset, limit uint) (users []model.User, err error)
	FindByEmail(ctx context.Context, email string) (user *model.User, err error)
	Update(ctx context.Context, id uint, user *model.User) (err error)
	Delete(ctx context.Context, id uint) (err error)
}

// var _ userRepo = (*repository.UserRepo)(nil)

type UserService struct {
	r userRepo
}

func NewUserService(r userRepo) (s *UserService) {
	s = new(UserService{r})
	return
}

func (s *UserService) Create(ctx context.Context, in *model.UserInput) (id uint, err error) {
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

	user := new(model.User{
		Name:         in.Name,
		Email:        (*string)(in.Email),
		PasswordHash: pwHash,
		RoleID:       in.RoleID,
	})

	if err = s.r.Create(ctx, user); err != nil {
		return
	}

	id = user.ID
	return
}

func (s *UserService) Finds(ctx context.Context, limit, offset uint) (users []model.User, err error) {
	if limit == 0 {
		limit = 5
	}
	if limit > 50 {
		err = ErrLimitTooLarge
		return
	}

	users, err = s.r.Finds(ctx, offset, limit)
	return
}

func (s *UserService) FindByEmail(ctx context.Context, email helper.Email) (user *model.User, err error) {
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

	err = s.r.Update(ctx, id, new(model.User{
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
