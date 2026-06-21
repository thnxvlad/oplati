package hmiddlewares

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type RateLimiterStorage interface {
	IncrementWithTTL(ctx context.Context, key string, limit int64, interval time.Duration) (bool, error)
}

const maxRequestPerInterval = 3
const timeIntervalRequests time.Duration = 60 * time.Second

func RateLimiterRedisMiddleware(storage RateLimiterStorage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIP, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			userAgent := r.UserAgent()
			key := fmt.Sprintf("%s:%s", userIP, userAgent)
			
			ok, err := storage.IncrementWithTTL(r.Context(), key, maxRequestPerInterval, timeIntervalRequests)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !ok {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}