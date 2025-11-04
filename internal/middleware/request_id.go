package middleware

import (
	"net/http"

	"github.com/sariya23/game_service/internal/lib/generate"
)

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generate.GenerateRequestID()
		w.Header().Set("X-Request-Id", requestID)
		next.ServeHTTP(w, r)
	})
}
