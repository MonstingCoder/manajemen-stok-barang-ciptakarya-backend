package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/config"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) (err error)
	Finds(ctx context.Context, limit, offset uint) (users []model.User, err error)
	FindByEmail(ctx context.Context, email string) (user *model.User, err error)
	Update(ctx context.Context, id uint, user *model.User) (err error)
	Delete(ctx context.Context, id uint) (err error)
}

var (
	// _ UserRepo = (*repository.UserRepo)(nil)

	ErrLimitTooLarge = errors.New("ERROR LIMIT TO LARGE")
	ErrUserNotFound  = errors.New("ERROR USER NOT FOUND")
	ErrInvalidUserPw = errors.New("ERROR INVALID USER PASSWORD")
)

type UserService struct {
	r UserRepo
}

func NewUserService(r UserRepo) (s *UserService) {
	s = new(UserService{r})
	return
}

func (s *UserService) Reg(ctx context.Context, cfg *config.Config, in *model.UserIn) (token string, err error) {
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

	user := new(model.User{
		Name:         in.Name,
		Email:        (*string)(in.Email),
		PasswordHash: pwHash,
		RoleID:       in.RoleID,
	})
	err = s.r.Create(ctx, user)

	now := time.Now()
	token, err = CreateUserClaims(user, now, cfg)
	return
}

func (s *UserService) Login(ctx context.Context, cfg *config.Config, email model.Email, pass model.Pw) (token string, err error) {
	if err = email.Check(); err != nil {
		return
	}
	if err = pass.Check(); err != nil {
		return
	}

	user, err := s.r.FindByEmail(ctx, string(email))
	if err != nil {
		return
	}
	if user == nil {
		err = ErrUserNotFound
		return
	}
	if !pass.Verify(*user.PasswordHash) {
		err = ErrInvalidUserPw
		return
	}

	now := time.Now()
	token, err = CreateUserClaims(user, now, cfg)
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

func (s *UserService) Update(ctx context.Context, id uint, in *model.UserIn) (err error) {
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

func CreateUserClaims(user *model.User, now time.Time, cfg *config.Config) (token string, err error) {
	claims := new(model.UserClaims{
		User: new(model.UserOut{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role: new(model.RoleOut{
				ID:        *user.RoleID,
				Name:      user.Name,
				CreatedAt: &user.Role.CreatedAt,
				UpdatedAt: &user.Role.UpdatedAt,
				DeletedAt: &user.Role.DeletedAt,
			}),
			CreatedAt: &user.CreatedAt,
			UpdatedAt: &user.UpdatedAt,
			DeletedAt: &user.DeletedAt,
		}),
		RegisteredClaims: new(jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
		}),
	})

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(cfg.JWTSecret)
	return
}
