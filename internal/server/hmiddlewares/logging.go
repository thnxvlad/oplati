package hmiddlewares

import (
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		// TODO  реализовать логирование информации о запросах через zerolog
	})
}
