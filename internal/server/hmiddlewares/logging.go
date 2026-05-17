package hmiddlewares

import (
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode 	int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {

	// вывод в json
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:		http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)								
	})
}
