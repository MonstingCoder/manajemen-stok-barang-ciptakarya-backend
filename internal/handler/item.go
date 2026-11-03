package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type ItemService interface {
	Create(ctx context.Context, in *model.ItemIn) (item *model.Item, err error)
	Finds(ctx context.Context, limit, offset uint) (items []model.Item, err error)
	FindByName(ctx context.Context, name string) (item []model.Item, err error)
	Update(ctx context.Context, id uint, in *model.ItemIn) (err error)
	Delete(ctx context.Context, id uint) (err error)
}

// var _ ItemService = (*service.ItemService)(nil)

type ItemHandler struct {
	s ItemService
}

func NewItemHandler(s ItemService) (h *ItemHandler) {
	h = new(ItemHandler{s})
	return
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	body := new(model.ItemIn)
	if err := dec.Decode(body); err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	item, err := h.s.Create(r.Context(), body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusCreated, item)
}

func (h *ItemHandler) Finds(w http.ResponseWriter, r *http.Request) {
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

	items, err := h.s.Finds(r.Context(), uint(limit), uint(offset))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, items)
}

func (h *ItemHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "name path parameter is required"})
		return
	}

	items, err := h.s.FindByName(r.Context(), name)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, items)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "id path parameter must be a positive integer"})
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	body := new(model.ItemIn)
	if err = dec.Decode(body); err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	if err = h.s.Update(r.Context(), uint(id), body); err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, map[string]string{"message": "item updated successfully"})
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "id path parameter must be a positive integer"})
		return
	}

	if err = h.s.Delete(r.Context(), uint(id)); err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, map[string]string{"message": "item deleted successfully"})
}
