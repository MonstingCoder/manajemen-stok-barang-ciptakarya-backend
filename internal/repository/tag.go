package repository

import (
	"context"
	"errors"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
	"gorm.io/gorm"
)

var (
	ErrTagExists    = errors.New("ERROR TAG EXISTS")
	ErrTagNameTaken = errors.New("ERROR TAG NAME TAKEN")
	ErrTagNotFound  = errors.New("ERROR TAG NOT FOUND")
)

type TagRepo struct {
	db *gorm.DB
}

func NewTagRepo(db *gorm.DB) (r *TagRepo) {
	r = new(TagRepo{db})
	return
}

func (r *TagRepo) Creates(ctx context.Context, tags []model.Tag) (err error) {
	if err = gorm.G[model.Tag](r.db).CreateInBatches(ctx, &tags, len(tags)); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = ErrTagExists
	}

	return
}

func (r *TagRepo) Create(ctx context.Context, tag *model.Tag) (err error) {
	if err = gorm.G[model.Tag](r.db).Create(ctx, tag); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = ErrTagExists
		if tag.ID == 0 {
			err = ErrTagNameTaken
		}
	}

	return
}

func (r *TagRepo) Finds(ctx context.Context, limit uint, offset uint) (tags []model.Tag, err error) {
	if tags, err = gorm.G[model.Tag](r.db).Limit(int(limit)).Offset(int(offset)).Find(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrTagNotFound
	}

	return
}

func (r *TagRepo) FindByName(ctx context.Context, name string) (tags []model.Tag, err error) {
	if tags, err = gorm.G[model.Tag](r.db).Where("name = ?", name).Find(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrTagNotFound
	}

	return
}

func (r *TagRepo) Delete(ctx context.Context, id uint) (err error) {
	if _, err = gorm.G[model.Tag](r.db).Where("id = ?", id).Delete(ctx); err != nil {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = gorm.ErrRecordNotFound
	}

	return
}
