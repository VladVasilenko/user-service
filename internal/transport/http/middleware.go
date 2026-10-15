package http

import (
	"log"
	"net/http"
	"runtime/debug"
)

type responseWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.wroteHeader = true
	}
	return rw.ResponseWriter.Write(b)
}

func panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w}

		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC] %v\n%s", rec, debug.Stack())
				if !rw.wroteHeader {
					http.Error(rw, "internal server error", http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(rw, r)
	})
}
