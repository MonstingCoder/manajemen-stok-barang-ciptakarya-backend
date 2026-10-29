package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/apperror"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/entity"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/helper"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/service"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/utility"
)

type UserService interface {
	Create(ctx context.Context, in *model.UserInput) (*entity.User, error)
	Finds(ctx context.Context, limit, offset uint) ([]entity.User, error)
	FindByEmail(ctx context.Context, email helper.Email) (*entity.User, error)
	Update(ctx context.Context, id uint, in *model.UserInput) error
	Delete(ctx context.Context, id uint) error
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(s UserService) (h *UserHandler) {
	h = new(UserHandler{service: s})
	return
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	//
}

// code dibawah hasil generate copilot :

type userRequest struct {
	Name     *string          `json:"name"`
	Email    *helper.Email    `json:"email"`
	Password *helper.Password `json:"password"`
	RoleID   *uint            `json:"role_id"`
}

type userResponse struct {
	ID        uint      `json:"id"`
	Name      *string   `json:"name"`
	Email     *string   `json:"email"`
	RoleID    *uint     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	in, ok := decodeUserInput(w, r)
	if !ok {
		return
	}

	user, err := h.service.Create(r.Context(), in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if user == nil {
		utility.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	utility.WriteResponse(w, http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) Finds(w http.ResponseWriter, r *http.Request) {
	limit, err := queryUint(r, "limit")
	if err != nil {
		utility.WriteError(w, http.StatusBadRequest, "invalid limit")
		return
	}
	offset, err := queryUint(r, "offset")
	if err != nil {
		utility.WriteError(w, http.StatusBadRequest, "invalid offset")
		return
	}

	users, err := h.service.Finds(r.Context(), limit, offset)
	if err != nil {
		utility.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]userResponse, len(users))
	for i := range users {
		responses[i] = toUserResponse(&users[i])
	}
	utility.WriteError(w, http.StatusOK, responses)
}

func (h *UserHandler) FindByEmail(w http.ResponseWriter, r *http.Request) {
	email := helper.Email(r.URL.Query().Get("email"))
	user, err := h.service.FindByEmail(r.Context(), email)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if user == nil {
		utility.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	utility.WriteError(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	in, ok := decodeUserInput(w, r)
	if !ok {
		return
	}

	if err := h.service.Update(r.Context(), id, in); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeUserInput(w http.ResponseWriter, r *http.Request) (*model.UserInput, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var body userRequest
	if err := decoder.Decode(&body); err != nil {
		utility.WriteError(w, http.StatusBadRequest, "invalid request body")
		return nil, false
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		utility.WriteError(w, http.StatusBadRequest, "request body must contain one JSON object")
		return nil, false
	}
	if body.Name == nil || body.Email == nil || body.Password == nil || body.RoleID == nil {
		utility.WriteError(w, http.StatusBadRequest, "name, email, password, and role_id are required")
		return nil, false
	}

	return &model.UserInput{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
		RoleID:   body.RoleID,
	}, true
}

func queryUint(r *http.Request, key string) (uint, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseUint(value, 10, strconv.IntSize)
	return uint(parsed), err
}

func pathID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		utility.WriteError(w, http.StatusBadRequest, "invalid user id")
		return 0, false
	}
	return uint(id), true
}

func toUserResponse(user *entity.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrLimitTooLarge):
		utility.WriteError(w, http.StatusBadRequest, "limit must not exceed 50")
	case errors.Is(err, helper.ErrEmptyEmail),
		errors.Is(err, helper.ErrEmailLen),
		errors.Is(err, helper.ErrInvalidEmail),
		errors.Is(err, helper.ErrPwLen),
		errors.Is(err, helper.ErrInvalidPw):
		utility.WriteError(w, http.StatusBadRequest, "invalid user data")
	case errors.Is(err, apperror.ErrUserNotFound):
		utility.WriteError(w, http.StatusNotFound, "user not found")
	case errors.Is(err, apperror.ErrUserEmailTaken),
		errors.Is(err, apperror.ErrUserExists):
		utility.WriteError(w, http.StatusConflict, "user already exists")
	default:
		utility.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
