package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/config"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/utility"
)

type UserService interface {
	Reg(ctx context.Context, in *model.UserIn) (user *model.User, err error)
	Login(ctx context.Context, email model.Email, pw model.Pw) (user *model.User, err error)
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
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	user, err := h.s.Reg(r.Context(), body)
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	token, err := CreateUserClaims(user, time.Now(), h.cfg)
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	utility.WriteResponse(w, http.StatusCreated, map[string]string{"token": token})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	body := new(model.UserLoginReq)
	if err := dec.Decode(body); err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	user, err := h.s.Login(r.Context(), model.Email(*body.Email), model.Pw(*body.Password))
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	token, err := CreateUserClaims(user, time.Now(), h.cfg)
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	utility.WriteResponse(w, http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) Finds(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit, err := strconv.ParseUint(q.Get("limit"), 10, 0)
	if err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "limit query parameter must be a positive integer"})
		return
	}

	offset, err := strconv.ParseUint(q.Get("offset"), 10, 0)
	if err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "offset query parameter must be a positive integer"})
		return
	}

	users, err := h.s.Finds(r.Context(), uint(limit), uint(offset))
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	utility.WriteResponse(w, http.StatusOK, users)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "id path parameter must be a positive integer"})
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	body := new(model.UserIn)
	if err := dec.Decode(body); err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	if err = h.s.Update(r.Context(), uint(id), body); err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	utility.WriteResponse(w, http.StatusOK, map[string]string{"message": "user updated successfully"})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "id path parameter must be a positive integer"})
		return
	}

	if err = h.s.Delete(r.Context(), uint(id)); err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	utility.WriteResponse(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}

func CreateUserClaims(user *model.User, now time.Time, cfg *config.Config) (token string, err error) {
	claims := new(model.UserClaims{
		RegisteredClaims: new(jwt.RegisteredClaims{
			ID: strconv.FormatUint(uint64(user.ID), 10),
		}),
		Role: user.Role.Name,
	})

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(cfg.JWTSecret)
	return
}
