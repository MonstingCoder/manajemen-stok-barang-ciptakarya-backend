package repository

import (
	"context"
	"errors"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
	"gorm.io/gorm"
)

var (
	ErrRoleExists   = errors.New("ERROR ROLE EXISTS")
	ErrRoleNotFound = errors.New("ERROR ROLE NOT FOUND")
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) (r *RoleRepo) {
	r = new(RoleRepo{db})
	return
}

func (r *RoleRepo) Creates(ctx context.Context, roles []model.Role) (err error) {
	if err = gorm.G[model.Role](r.db).CreateInBatches(ctx, &roles, len(roles)); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = ErrRoleExists
	}

	return
}

func (r *RoleRepo) Delete(ctx context.Context, id uint) (err error) {
	if _, err = gorm.G[model.Role](r.db).Where("id = ?", id).Delete(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrRoleNotFound
	}

	return
}
