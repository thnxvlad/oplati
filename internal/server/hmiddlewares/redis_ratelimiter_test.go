package hmiddlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockRateLimitStorage struct {
	allowFunc func(key string) (bool, error)
}

func (m *mockRateLimitStorage) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	return m.allowFunc(key)
}

func TestRedisRateLimiterMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		mockAllow      bool
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "Success - Request Allowed",
			mockAllow:      true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Rate Limit Exceed",
			mockAllow:      false,
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "Storage Error",
			mockAllow:      false,
			mockErr:        errors.New("redis is down"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRateLimitStorage{
				allowFunc: func(key string) (bool, error) {
					return tt.mockAllow, tt.mockErr
				},
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := RedisRateLimiterMiddleware(mock, 3, time.Second*3)
			testHandler := middleware(nextHandler)

			req := httptest.NewRequest("POST", "/signup", nil)
			req.RemoteAddr = "127.0.0.1:8083"
			rr := httptest.NewRecorder()

			testHandler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
