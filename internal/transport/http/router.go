package http

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type Handlers struct {
	User *UserHandler
}

func InitRoutes(h Handlers, db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", h.User.GetUser)
	mux.HandleFunc("GET /health", healthHandler(db))
	return chain(mux, panicRecoveryMiddleware)
}

func healthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			http.Error(w, "db unhealthy", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
