package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"user-service/internal/model"
	"user-service/internal/service"
)

type UserService interface {
	GetUserProfile(ctx context.Context, id int64) (*model.User, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{userService: s}
}

type userResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	user, err := h.userService.GetUserProfile(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		case errors.Is(err, service.ErrUserNotFound):
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "user not found"})
		default:
			log.Printf("get user profile: %v", err)
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, userResponse{
		ID:   user.ID,
		Name: user.Name,
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}
