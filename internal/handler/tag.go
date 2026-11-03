package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	body := new(model.TagIn)
	if err := dec.Decode(body); err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	tag, err := h.s.Create(r.Context(), body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusCreated, tag)
}

func (h *TagHandler) Finds(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, err := strconv.ParseUint(q.Get("limit"), 10, 0)
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "limit query parameter is required"})
		return
	}
	offset, err := strconv.ParseUint(q.Get("offset"), 10, 0)
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "offset query parameter is required"})
		return
	}

	tags, err := h.s.Finds(r.Context(), uint(limit), uint(offset))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, tags)
}

func (h *TagHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "name query parameter is required"})
		return
	}

	tags, err := h.s.FindByName(r.Context(), name)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, tags)
}

func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 0)
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "id query parameter is required"})
		return
	}

	err = h.s.Delete(r.Context(), uint(id))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, map[string]string{"message": "Tag deleted successfully"})
}
