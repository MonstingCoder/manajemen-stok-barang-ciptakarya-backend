package service

import (
	"context"
	"errors"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) (err error)
	Finds(ctx context.Context, limit, offset uint) (users []model.User, err error)
	FindByEmail(ctx context.Context, email string) (user *model.User, err error)
	FindByID(ctx context.Context, id uint) (user *model.User, err error)
	Update(ctx context.Context, id uint, user *model.User) (err error)
	Delete(ctx context.Context, id uint) (err error)
}

var (
	// _ UserRepo = (*repository.UserRepo)(nil)

	ErrLimitTooLarge = errors.New("ERROR LIMIT TO LARGE")
)

type UserService struct {
	r UserRepo
}

func NewUserService(r UserRepo) (s *UserService) {
	s = new(UserService{r})
	return
}

func (s *UserService) Reg(ctx context.Context, in *model.UserIn) (user *model.User, err error) {
	if err = in.Email.Check(); err != nil {
		return
	}
	if err = in.Pw.Check(); err != nil {
		return
	}

	pwHash, err := in.Pw.Hash()
	if err != nil {
		return
	}

	user = new(model.User{
		Name:         in.Name,
		Email:        (*string)(in.Email),
		PasswordHash: pwHash,
		RoleID:       in.RoleID,
	})
	err = s.r.Create(ctx, user)

	return
}

func (s *UserService) Login(ctx context.Context, email model.Email, pw model.Pw) (user *model.User, err error) {
	if err = email.Check(); err != nil {
		return
	}
	if err = pw.Check(); err != nil {
		return
	}

	user, err = s.r.FindByEmail(ctx, string(email))
	if err != nil {
		return
	}

	err = pw.Verify(*user.PasswordHash)
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

	users, err = s.r.Finds(ctx, limit, offset)
	return
}

func (s *UserService) FindByID(ctx context.Context, id uint) (user *model.User, err error) {
	user, err = s.r.FindByID(ctx, id)
	return
}

func (s *UserService) Update(ctx context.Context, id uint, in *model.UserIn) (err error) {
	if in == nil {
		return errors.New("user update is empty")
	}

	user := new(model.User)
	user.Name = in.Name
	user.RoleID = in.RoleID

	if in.Email != nil {
		if err = in.Email.Check(); err != nil {
			return
		}
		email := string(*in.Email)
		user.Email = &email
	}

	if in.Pw != nil {
		if err = in.Pw.Check(); err != nil {
			return
		}
		if user.PasswordHash, err = in.Pw.Hash(); err != nil {
			return
		}
	}

	if user.Name == nil && user.Email == nil && user.PasswordHash == nil && user.RoleID == nil {
		return errors.New("user update is empty")
	}

	return s.r.Update(ctx, id, user)
}

func (s *UserService) Delete(ctx context.Context, id uint) (err error) {
	err = s.r.Delete(ctx, id)
	return
}
