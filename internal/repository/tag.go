package repository

import (
	"context"
	"errors"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/apperror"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/entity"
	"gorm.io/gorm"
)

type TagRepo struct {
	db *gorm.DB
}

func NewTagRepo(db *gorm.DB) (r *TagRepo) {
	r = new(TagRepo{db})
	return
}

func (r *TagRepo) Creates(ctx context.Context, tags []entity.Tag) (err error) {
	if err = gorm.G[entity.Tag](r.db).CreateInBatches(ctx, &tags, len(tags)); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = apperror.ErrTagExists
	}

	return
}

func (r *TagRepo) Create(ctx context.Context, tag *entity.Tag) (err error) {
	if err = gorm.G[entity.Tag](r.db).Create(ctx, tag); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = apperror.ErrTagExists
		if tag.ID == 0 {
			err = apperror.ErrTagNameTaken
		}
	}

	return
}

func (r *TagRepo) Finds(ctx context.Context, limit uint, offset uint) (tags []entity.Tag, err error) {
	if tags, err = gorm.G[entity.Tag](r.db).Limit(int(limit)).Offset(int(offset)).Find(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = apperror.ErrTagNotFound
	}

	return
}

func (r *TagRepo) FindByName(ctx context.Context, name string) (tags []entity.Tag, err error) {
	if tags, err = gorm.G[entity.Tag](r.db).Where("name = ?", name).Find(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = apperror.ErrTagNotFound
	}

	return
}

func (r *TagRepo) Delete(ctx context.Context, id uint) (err error) {
	if _, err = gorm.G[entity.Tag](r.db).Where("id = ?", id).Delete(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = gorm.ErrRecordNotFound
	}

	return
}
