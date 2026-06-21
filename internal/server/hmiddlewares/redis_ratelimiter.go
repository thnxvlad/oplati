package hmiddlewares

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type RateLimiterStorage interface {
	Allow(ctx context.Context,key string, limit int64, window time.Duration) (bool, error)
}

func RedisRateLimiterMiddleware(storage RateLimiterStorage, limit int64, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIP, _, err := net.SplitHostPort(r.RemoteAddr)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			userAgent := r.UserAgent()
			key := fmt.Sprintf("%s:%s", userIP, userAgent)

			ok, err := storage.Allow(r.Context(), key, limit, window)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !ok {
				w.WriteHeader(http.StatusTooManyRequests)
				return 
			}
			next.ServeHTTP(w,r)
		})
	}
}
