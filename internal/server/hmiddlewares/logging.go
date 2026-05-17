package hmiddlewares

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
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
	output := zerolog.ConsoleWriter{
		Out: 		os.Stdout,
		TimeFormat: "15:04:05",
		NoColor: 	false,		
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:		http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", wrapped.statusCode).
				Msg("HTTP Request processed")
								
	})
}
