package repository

import (
	"context"
	"errors"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
	"gorm.io/gorm"
)

var (
	ErrItemExists   = errors.New("ERROR ITEM EXISTS")
	ErrItemNotFound = errors.New("ERROR ITEM NOT FOUND")
)

type ItemRepo struct {
	db *gorm.DB
}

func NewItemRepo(db *gorm.DB) (r *ItemRepo) {
	r = new(ItemRepo{db})
	return
}

func (r *ItemRepo) Create(ctx context.Context, item *model.Item) (err error) {
	if err = gorm.G[model.Item](r.db).Create(ctx, item); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = ErrItemExists
	}

	return
}

func (r *ItemRepo) Finds(ctx context.Context, limit, offset uint) (items []model.Item, err error) {
	if items, err = gorm.G[model.Item](r.db).Limit(int(limit)).Offset(int(offset)).Find(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrItemNotFound
	}

	return
}

func (r *ItemRepo) FindByName(ctx context.Context, name string) (items []model.Item, err error) {
	if items, err = gorm.G[model.Item](r.db).Where("name = ?", name).Find(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrItemNotFound
	}

	return
}

func (r *ItemRepo) Update(ctx context.Context, id uint, item *model.Item) (err error) {
	if _, err = gorm.G[model.Item](r.db).Where("id = ?", id).Updates(ctx, *item); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrItemNotFound
	}

	return
}

func (r *ItemRepo) Delete(ctx context.Context, id uint) (err error) {
	if _, err = gorm.G[model.Item](r.db).Where("id = ?", id).Delete(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrItemNotFound
	}

	return
}
