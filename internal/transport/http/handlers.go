package http

import (
	"fmt"
	"net/http"
	"strconv"
	"user-service/internal/service"
)

type UserHandler struct {
	userService *service.UserService // Теперь зависим строго от слоя логики
}

// Конструктор теперь принимает указатель на сервис
func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

func (h *UserHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Запрашиваем профиль у сервиса, а не у репозитория напрямую
	name, err := h.userService.GetUserProfile(id)
	if err != nil {
		// Проверяем текст ошибки, которую вернул сервис
		if err.Error() == "пользователь не найден" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Любая другая системная ошибка (например, упала БД)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintf(w, "Привет, %s\n", name)
}
