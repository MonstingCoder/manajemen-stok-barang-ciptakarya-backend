package service

import "github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/repository"

type tagRepo interface {
	//
}

var _ tagRepo = (*repository.TagRepo)(nil)

type TagService struct {
	r tagRepo
}

func NewTagService(r tagRepo) (s *TagService) {
	s = new(TagService{r})
	return
}
