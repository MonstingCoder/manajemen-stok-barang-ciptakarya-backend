package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/config"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type UserService interface {
	Reg(ctx context.Context, cfg *config.Config, in *model.UserIn) (token string, err error)
	Login(ctx context.Context, cfg *config.Config, email model.Email, pw model.Pw) (token string, err error)
	Finds(ctx context.Context, limit, offset uint) (users []model.User, err error)
	Update(ctx context.Context, id uint, in *model.UserIn) (err error)
	Delete(ctx context.Context, id uint) (err error)
}

// var _ UserService = (*service.UserService)(nil)

type UserHandler struct {
	s   UserService
	cfg *config.Config
}

func NewUserHandler(s UserService, cfg *config.Config) (h *UserHandler) {
	h = new(UserHandler{s, cfg})
	return
}

func (h *UserHandler) Reg(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	body := new(model.UserIn)
	if err := dec.Decode(body); err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	token, err := h.s.Reg(r.Context(), h.cfg, body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusCreated, map[string]string{"token": token})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	body := new(model.UserLoginReq)
	if err := dec.Decode(body); err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	token, err := h.s.Login(r.Context(), h.cfg, model.Email(*body.Email), model.Pw(*body.Password))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) Finds(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit, err := strconv.ParseUint(q.Get("limit"), 10, 0)
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "limit query parameter must be a positive integer"})
		return
	}

	offset, err := strconv.ParseUint(q.Get("offset"), 10, 0)
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "offset query parameter must be a positive integer"})
		return
	}

	users, err := h.s.Finds(r.Context(), uint(limit), uint(offset))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, users)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "id path parameter must be a positive integer"})
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	body := new(model.UserIn)
	if err := dec.Decode(body); err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	if err = h.s.Update(r.Context(), uint(id), body); err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, map[string]string{"message": "user updated successfully"})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": "id path parameter must be a positive integer"})
		return
	}

	if err = h.s.Delete(r.Context(), uint(id)); err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}
