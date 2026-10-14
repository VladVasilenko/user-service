package http

import (
	"log"
	"net/http"
)

// InitRoutes настраивает роутер, регистрирует хендлеры и навешивает мидлвари
func InitRoutes(handler *UserHandler) http.Handler {
	mux := http.NewServeMux()

	// Регистрируем эндпоинт, привязывая его к методу структуры UserHandler
	mux.HandleFunc("/user", handler.GetPlayer)

	// Оборачиваем роутер в наше защитное Middleware от паник
	return panicRecoveryMiddleware(mux)
}

// Внутренний Middleware для перехвата паник
func panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] Перехвачена критическая ошибка: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Внутренняя ошибка сервера"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
