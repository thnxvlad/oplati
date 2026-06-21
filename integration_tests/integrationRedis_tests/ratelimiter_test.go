package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"

	hm "github.com/thnxvlad/oplati/internal/server/hmiddlewares"
)

var redisCli *redis.Client

type RedisRateLimiterStorage struct {
	rdb *redis.Client
}

func NewRedisRateLimiterStorage(rdb *redis.Client) *RedisRateLimiterStorage {
	return &RedisRateLimiterStorage{rdb: rdb}
}

func (s *RedisRateLimiterStorage) IncrementWithTTL(ctx context.Context, key string, limit int64, interval time.Duration) (bool, error) {
	pipe := s.rdb.TxPipeline()

	incr := pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, interval)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	if incr.Val() > limit {
		return false, nil
	}

	return true, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	redisURL := os.Getenv("REDIS_URL")

	var cli *redis.Client
	var rContainer *rediscontainer.RedisContainer
	var err error

	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}

		cli = redis.NewClient(opt)
		if err := cli.Ping(ctx).Err(); err != nil {
			panic(err)
		}
	} else {
		rContainer, err = rediscontainer.Run(ctx, "redis:7-alpine")
		if err != nil {
			panic(err)
		}

		u, err := rContainer.ConnectionString(ctx)
		if err != nil {
			panic(err)
		}

		opt, err := redis.ParseURL(u)
		if err != nil {
			panic(err)
		}

		cli = redis.NewClient(opt)
		if err := cli.Ping(ctx).Err(); err != nil {
			panic(err)
		}
	}

	redisCli = cli

	code := m.Run()

	_ = redisCli.Close()
	if rContainer != nil {
		_ = rContainer.Terminate(ctx)
	}

	os.Exit(code)
}

func TestRateLimiterRedisMiddleware_AllowsUpToLimit(t *testing.T) {
	ctx := context.Background()

	_ = redisCli.Del(ctx, "192.168.100.1:TestAgent")

	st := NewRedisRateLimiterStorage(redisCli)
	mw := hm.RateLimiterRedisMiddleware(st)

	called := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(http.StatusOK)
	})

	handler := mw(next)

	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://test", nil)
		req.RemoteAddr = "192.168.100.1:5555"
		req.Header.Set("User-Agent", "TestAgent")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("запрос %d: ожидали 200, получили %d", i, rr.Code)
		}
	}

	if called != 3 {
		t.Fatalf("next вызван %d раз, ожидали 3", called)
	}

	req := httptest.NewRequest(http.MethodGet, "http://test", nil)
	req.RemoteAddr = "192.168.100.1:5555"
	req.Header.Set("User-Agent", "TestAgent")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("4-й запрос: ожидали 429, получили %d", rr.Code)
	}

	if called != 3 {
		t.Fatalf("на 429 next не должен вызываться, вызван %d раз", called)
	}
}

func TestRateLimiterRedisMiddleware_DifferentKeysSeparate(t *testing.T) {
	ctx := context.Background()

	_ = redisCli.Del(ctx, "10.0.0.2:curl")
	_ = redisCli.Del(ctx, "10.0.0.2:browser")

	st := NewRedisRateLimiterStorage(redisCli)
	mw := hm.RateLimiterRedisMiddleware(st)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := mw(next)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://test", nil)
		req.RemoteAddr = "10.0.0.2:1111"
		req.Header.Set("User-Agent", "curl")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("curl %d: ожидали 200, получили %d", i, rr.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "http://test", nil)
	req.RemoteAddr = "10.0.0.2:1111"
	req.Header.Set("User-Agent", "curl")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("curl 4-й: ожидали 429, получили %d", rr.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "http://test", nil)
	req2.RemoteAddr = "10.0.0.2:2222"
	req2.Header.Set("User-Agent", "browser")

	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Fatalf("browser: ожидали 200, получили %d", rr2.Code)
	}
}
