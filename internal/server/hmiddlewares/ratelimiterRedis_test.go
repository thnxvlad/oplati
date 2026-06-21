package hmiddlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockStorage struct {
	mockIncrement func(ctx context.Context, key string, limit int64, interval time.Duration) (bool, error)
}

func (m *MockStorage) IncrementWithTTL(ctx context.Context, key string, limit int64, interval time.Duration) (bool, error) {
	return m.mockIncrement(ctx, key, limit, interval)
}

func TestRateLimiterRedisMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	tests := []struct {
		name           string
		remoteAddr     string
		userAgent      string
		mockOk         bool
		mockErr        error
		expectedStatus int
		expectedBody   string
		expectedKey    string 
	}{
		{
			name:           "Успешный запрос (лимит не превышен)",
			remoteAddr:     "192.168.1.1:12345",
			userAgent:      "Mozilla/5.0",
			mockOk:         true,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
			expectedKey:    "192.168.1.1:Mozilla/5.0",
		},
		{
			name:           "Превышен лимит запросов",
			remoteAddr:     "10.0.0.1:8080",
			userAgent:      "curl/7.68.0",
			mockOk:         false,
			mockErr:        nil,
			expectedStatus: http.StatusTooManyRequests,
			expectedBody:   "Too Many Requests\n",
			expectedKey:    "10.0.0.1:curl/7.68.0",
		},
		{
			name:           "Ошибка хранилища (Redis упал)",
			remoteAddr:     "127.0.0.1:54321",
			userAgent:      "Postman",
			mockOk:         false,
			mockErr:        errors.New("redis connection error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "redis connection error\n",
			expectedKey:    "127.0.0.1:Postman",
		},
		{
			name:           "Ошибка парсинга IP (нет порта)",
			remoteAddr:     "invalid-ip-without-port",
			userAgent:      "Test",
			mockOk:         false,
			mockErr:        nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "missing port in address",
			expectedKey:    "", 
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{
				mockIncrement: func(ctx context.Context, key string, limit int64, interval time.Duration) (bool, error) {
					if key != tt.expectedKey {
						t.Errorf("expected key %q, got %q", tt.expectedKey, key)
					}
					if limit != maxRequestPerInterval {
						t.Errorf("expected limit %d, got %d", maxRequestPerInterval, limit)
					}
					if interval != timeIntervalRequests {
						t.Errorf("expected interval %v, got %v", timeIntervalRequests, interval)
					}
					return tt.mockOk, tt.mockErr
				},
			}

			middleware := RateLimiterRedisMiddleware(mockStorage)
			handlerToTest := middleware(nextHandler)

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			req.RemoteAddr = tt.remoteAddr
			req.Header.Set("User-Agent", tt.userAgent)

			rr := httptest.NewRecorder()

			handlerToTest.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want it to contain %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}