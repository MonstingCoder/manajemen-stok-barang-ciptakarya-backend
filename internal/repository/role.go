package repository

import (
	"context"
	"errors"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/apperror"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/entity"
	"gorm.io/gorm"
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) (r *RoleRepo) {
	r = new(RoleRepo{db})
	return
}

func (r *RoleRepo) Creates(ctx context.Context, roles []entity.Role) (err error) {
	if err = gorm.G[entity.Role](r.db).CreateInBatches(ctx, &roles, len(roles)); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = apperror.ErrRoleExists
	}

	return
}

func (r *RoleRepo) Delete(ctx context.Context, id uint) (err error) {
	if _, err = gorm.G[entity.Role](r.db).Where("id = ?", id).Delete(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = apperror.ErrRoleNotFound
	}

	return
}
