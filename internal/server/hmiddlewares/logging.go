package hmiddlewares

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	// вывод в json
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		var event *zerolog.Event
		if wrapped.statusCode >= 200 && wrapped.statusCode <= 399 {
			event = log.Info()
		} else if wrapped.statusCode <= 499 {
			event = log.Warn()
		} else {
			event = log.Error()
		}

		event.
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", wrapped.statusCode).
			Dur("latency", time.Since(start)).
			Msg("HTTP Request processed")
	})
}
