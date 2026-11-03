package handler

import (
	"context"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type TagService interface {
	Create(ctx context.Context, in *model.TagIn) (tag *model.Tag, err error)
	Finds(ctx context.Context, limit, offset uint) (tags []model.Tag, err error)
	FindByName(ctx context.Context, name string) (tag []model.Tag, err error)
	Delete(ctx context.Context, id uint) (err error)
}

// var _ TagService = (*service.TagService)(nil)

type TagHandler struct {
	s TagService
}

func NewTagHandler(s TagService) (h *TagHandler) {
	h = new(TagHandler{s})
	return
}
