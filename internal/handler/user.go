package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/config"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/middleware"
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
	if body.Name == nil || body.Email == nil || body.Pw == nil || body.RoleID == nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "name, email, password, and role_id are required"})
		return
	}

	user, err := h.s.Reg(r.Context(), body)
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	token, err := createUserClaims(user, h.cfg)
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
	if body.Email == nil || body.Password == nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "email and password are required"})
		return
	}

	user, err := h.s.Login(r.Context(), model.Email(*body.Email), model.Pw(*body.Password))
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	token, err := createUserClaims(user, h.cfg)
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

	ctx := r.Context()
	userIDKey, ok := ctx.Value(middleware.UserIDKey{}).(uint)
	if !ok {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "bad request"})
		return
	}
	if userIDKey != uint(id) {
		userRoleKey, ok := ctx.Value(middleware.UserRoleKey{}).(string)
		if !ok {
			utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "bad request"})
			return
		}

		if userRoleKey != "admin" {
			utility.WriteError(w, http.StatusForbidden, map[string]string{"message": "not allowed"})
			return
		}
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	body := new(model.UserIn)
	if err := dec.Decode(body); err != nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	if body.Name == nil && body.Email == nil && body.Pw == nil && body.RoleID == nil {
		utility.WriteError(w, http.StatusBadRequest, map[string]string{"message": "at least one user field is required"})
		return
	}

	if err = h.s.Update(r.Context(), uint(id), body); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, model.ErrEmptyEmail) ||
			errors.Is(err, model.ErrEmailLen) ||
			errors.Is(err, model.ErrInvalidEmail) ||
			errors.Is(err, model.ErrPwLen) ||
			errors.Is(err, model.ErrInvalidPw) {
			status = http.StatusBadRequest
		}
		utility.WriteError(w, status, map[string]string{"message": err.Error()})
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

	utility.WriteResponse(w, http.StatusNoContent, map[string]string{"message": "user deleted successfully"})
}

func createUserClaims(user *model.User, cfg *config.Config) (token string, err error) {
	if user == nil || user.Role == nil || user.Role.Name == nil || cfg == nil || cfg.JWTSecret == nil || cfg.JWTExpireTime == nil {
		return "", errors.New("cannot create token: user role and JWT configuration are required")
	}

	now := time.Now()
	claims := new(model.UserClaims{
		UserID: user.ID,
		Role:   user.Role.Name,
		RegisteredClaims: new(jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(*cfg.JWTExpireTime)),
		}),
	})

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(*cfg.JWTSecret))
	return
}
