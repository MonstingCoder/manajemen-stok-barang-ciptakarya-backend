package repository

import (
	"context"
	"errors"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrUserEmailTaken = errors.New("ERROR EMAIL TAKEN")
	ErrUserExists     = errors.New("ERROR USER ALREADY EXISTS")
	ErrUserNotFound   = errors.New("ERROR USER NOT FOUND")
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) (r *UserRepo) {
	r = new(UserRepo{db})
	return
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) (err error) {
	if err = gorm.G[model.User](r.db).Create(ctx, user); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = ErrUserExists
		if user.ID == 0 {
			err = ErrUserEmailTaken
		}
	}

	return
}

func (r *UserRepo) Finds(ctx context.Context, limit, offset uint) (users []model.User, err error) {
	if users, err = gorm.G[model.User](r.db).Preload(clause.Associations, nil).Limit(int(limit)).Offset(int(offset)).Find(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrUserNotFound
	}

	return
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (user *model.User, err error) {
	res, err := gorm.G[model.User](r.db).Preload(clause.Associations, nil).Where("email = ?", email).Take(ctx)
	if err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrUserNotFound
		return
	}

	user = &res
	return
}

func (r *UserRepo) Update(ctx context.Context, id uint, user *model.User) (err error) {
	if _, err = gorm.G[model.User](r.db).Where("id = ?", id).Updates(ctx, *user); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrUserNotFound
	}

	return
}

func (r *UserRepo) Delete(ctx context.Context, id uint) (err error) {
	if _, err = gorm.G[model.User](r.db).Where("id = ?", id).Delete(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrUserNotFound
	}

	return
}
