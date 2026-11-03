package service

import (
	"context"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type ItemRepo interface {
	Create(ctx context.Context, item *model.Item) (err error)
	Finds(ctx context.Context, limit, offset uint) (items []model.Item, err error)
	FindByName(ctx context.Context, name string) (item []model.Item, err error)
	Update(ctx context.Context, id uint, item *model.Item) (err error)
	Delete(ctx context.Context, id uint) (err error)
}

// var _ ItemRepo = (*repository.ItemRepo)(nil)

type ItemService struct {
	r ItemRepo
}

func NewItemService(r ItemRepo) (s *ItemService) {
	s = new(ItemService{r})
	return
}

func (s *ItemService) Create(ctx context.Context, in *model.ItemIn) (item *model.Item, err error) {
	item = new(model.Item{
		Name:        in.Name,
		Price:       in.Price,
		Quantity:    in.Quantity,
		Description: in.Description,
		Tags:        make([]model.Tag, len(in.Tags)),
	})

	for i, tag := range in.Tags {
		item.Tags[i] = model.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	err = s.r.Create(ctx, item)
	return
}

func (s *ItemService) Finds(ctx context.Context, limit, offset uint) (items []model.Item, err error) {
	items, err = s.r.Finds(ctx, limit, offset)
	return
}

func (s *ItemService) FindByName(ctx context.Context, name string) (item []model.Item, err error) {
	item, err = s.r.FindByName(ctx, name)
	return
}

func (s *ItemService) Update(ctx context.Context, id uint, in *model.ItemIn) (err error) {
	item := new(model.Item{
		Name:        in.Name,
		Price:       in.Price,
		Quantity:    in.Quantity,
		Description: in.Description,
		Tags:        make([]model.Tag, len(in.Tags)),
	})

	for i, tag := range in.Tags {
		item.Tags[i] = model.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	err = s.r.Update(ctx, id, item)
	return
}

func (s *ItemService) Delete(ctx context.Context, id uint) (err error) {
	err = s.r.Delete(ctx, id)
	return
}
