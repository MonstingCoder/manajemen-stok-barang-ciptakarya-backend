package service

import (
	"context"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/entity"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type ItemRepo interface {
	Create(ctx context.Context, item *entity.Item) (err error)
	Finds(ctx context.Context, limit, offset uint) (items []entity.Item, err error)
	FindByName(ctx context.Context, name string) (item []entity.Item, err error)
	Update(ctx context.Context, id uint, item *entity.Item) (err error)
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

func (s *ItemService) Create(ctx context.Context, in *model.ItemInput) (item *entity.Item, err error) {
	item = new(entity.Item{
		Name:        in.Name,
		Price:       in.Price,
		Quantity:    in.Quantity,
		Description: in.Description,
		Tags:        make([]entity.Tag, len(in.Tags)),
	})

	for i, tag := range in.Tags {
		item.Tags[i] = entity.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	err = s.r.Create(ctx, item)
	return
}

func (s *ItemService) Finds(ctx context.Context, limit, offset uint) (items []entity.Item, err error) {
	items, err = s.r.Finds(ctx, limit, offset)
	return
}

func (s *ItemService) FindByName(ctx context.Context, name string) (item []entity.Item, err error) {
	item, err = s.r.FindByName(ctx, name)
	return
}

func (s *ItemService) Update(ctx context.Context, id uint, in *model.ItemInput) (err error) {
	item := new(entity.Item{
		Name:        in.Name,
		Price:       in.Price,
		Quantity:    in.Quantity,
		Description: in.Description,
		Tags:        make([]entity.Tag, len(in.Tags)),
	})

	for i, tag := range in.Tags {
		item.Tags[i] = entity.Tag{
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
