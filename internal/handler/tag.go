package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/utility"
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
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	tag, err := h.s.Create(r.Context(), body)
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	utility.WriteResponse(w, http.StatusCreated, tag)
}

func (h *TagHandler) Finds(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, err := strconv.ParseUint(q.Get("limit"), 10, 0)
	if err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "limit query parameter is required"})
		return
	}
	offset, err := strconv.ParseUint(q.Get("offset"), 10, 0)
	if err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "offset query parameter is required"})
		return
	}

	tags, err := h.s.Finds(r.Context(), uint(limit), uint(offset))
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	utility.WriteResponse(w, http.StatusOK, tags)
}

func (h *TagHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "name path parameter is required"})
		return
	}

	tags, err := h.s.FindByName(r.Context(), name)
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	utility.WriteResponse(w, http.StatusOK, tags)
}

func (h *TagHandler) FindItemsByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "tag name path parameter is required"})
		return
	}

	q := r.URL.Query()
	limit, err := strconv.ParseUint(q.Get("limit"), 10, 0)
	if err != nil || limit == 0 {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "limit query parameter must be a positive integer"})
		return
	}
	offset, err := strconv.ParseUint(q.Get("offset"), 10, 0)
	if err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "offset query parameter must be a non-negative integer"})
		return
	}

	tags, err := h.s.FindByName(r.Context(), name)
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	items := []model.Item{}
	if len(tags) > 0 && offset < uint64(len(tags[0].Items)) {
		start := int(offset)
		end := len(tags[0].Items)
		if limit < uint64(end-start) {
			end = start + int(limit)
		}
		items = tags[0].Items[start:end]
		if items == nil {
			items = []model.Item{}
		}
	}
	utility.WriteResponse(w, http.StatusOK, items)
}

func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "id path parameter must be a positive integer"})
		return
	}

	err = h.s.Delete(r.Context(), uint(id))
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	utility.WriteResponse(w, http.StatusNoContent, map[string]string{"message": "Tag deleted successfully"})
}
