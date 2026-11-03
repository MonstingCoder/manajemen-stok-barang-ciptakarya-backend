package service

import (
	"context"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type TagRepo interface {
	Creates(ctx context.Context, tags []model.Tag) (err error)
	Create(ctx context.Context, tag *model.Tag) (err error)
	Finds(ctx context.Context, limit, offset uint) (tags []model.Tag, err error)
	FindByName(ctx context.Context, name string) (tag []model.Tag, err error)
	Delete(ctx context.Context, id uint) (err error)
}

// var _ TagRepo = (*repository.TagRepo)(nil)

type TagService struct {
	r TagRepo
}

func NewTagService(r TagRepo) (s *TagService) {
	s = new(TagService{r})
	return
}

func (s *TagService) Creates(ctx context.Context, ins []model.TagIn) (tags []model.Tag, err error) {
	tags = make([]model.Tag, len(ins))
	for i, in := range ins {
		tags[i] = model.Tag{
			Name:  in.Name,
			Items: make([]model.Item, len(in.Items)),
		}

		for j, item := range in.Items {
			tags[i].Items[j] = model.Item{
				Name:        item.Name,
				Price:       item.Price,
				Quantity:    item.Quantity,
				Description: item.Description,
			}
		}
	}

	err = s.r.Creates(ctx, tags)
	return
}

func (s *TagService) Create(ctx context.Context, in *model.TagIn) (tag *model.Tag, err error) {
	tag = new(model.Tag{
		Name:  in.Name,
		Items: make([]model.Item, len(in.Items)),
	})

	for i, item := range in.Items {
		tag.Items[i] = model.Item{
			Name:        item.Name,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Description: item.Description,
		}
	}

	err = s.r.Create(ctx, tag)
	return
}

func (s *TagService) Finds(ctx context.Context, limit, offset uint) (tags []model.Tag, err error) {
	tags, err = s.r.Finds(ctx, limit, offset)
	return
}

func (s *TagService) FindByName(ctx context.Context, name string) (tag []model.Tag, err error) {
	tag, err = s.r.FindByName(ctx, name)
	return
}

func (s *TagService) Delete(ctx context.Context, id uint) (err error) {
	err = s.r.Delete(ctx, id)
	return
}
