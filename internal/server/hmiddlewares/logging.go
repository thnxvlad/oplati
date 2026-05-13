package hmiddlewares

import (
	"context"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.New(os.Stdout).
			With().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Logger()
		logger.Debug().Msg("got request")
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "logger", logger)))
	})
}
