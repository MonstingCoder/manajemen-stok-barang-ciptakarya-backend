package handler

import (
	"context"
	"net/http"
	"strconv"
)

type RoleService interface {
	Delete(ctx context.Context, id uint) (err error)
}

// var _ RoleService = (*service.RoleService)(nil)

type RoleHandler struct {
	s RoleService
}

func NewRoleHandler(s RoleService) (h *RoleHandler) {
	h = new(RoleHandler{s})
	return
}

func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	if err = h.s.Delete(r.Context(), uint(id)); err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	WriteResponse(w, http.StatusOK, map[string]string{"message": "role deleted successfully"})
}
