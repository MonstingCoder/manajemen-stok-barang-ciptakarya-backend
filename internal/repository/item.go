package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
	"gorm.io/gorm"
)

type ItemRepo struct {
	db *gorm.DB
}

func NewItemRepo(db *gorm.DB) (r *ItemRepo) {
	r = new(ItemRepo{db})
	return
}

var (
	ErrItemExists   = errors.New("ERROR ITEM EXISTS")
	ErrItemNotFound = errors.New("ERROR ITEM NOT FOUND")
)

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
	if items, err = gorm.G[model.Item](r.db).Preload("Tags", nil).Limit(int(limit)).Offset(int(offset)).Find(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrItemNotFound
	}

	return
}

func (r *ItemRepo) FindByName(ctx context.Context, name string) (items []model.Item, err error) {
	if items, err = gorm.G[model.Item](r.db).Preload("Tags", nil).Where("name ILIKE ?", fmt.Sprintf("%%%s%%", name)).Find(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrItemNotFound
	}

	return
}

func (r *ItemRepo) Update(ctx context.Context, id uint, item *model.Item) (err error) {
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.Item
		if err := tx.First(&existing, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrItemNotFound
			}
			return err
		}

		if err := tx.Model(&existing).
			Omit("Tags").
			Select("Name", "Price", "Quantity", "Description").
			Updates(item).Error; err != nil {
			return err
		}

		return tx.Model(&existing).Association("Tags").Replace(item.Tags)
	})
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
