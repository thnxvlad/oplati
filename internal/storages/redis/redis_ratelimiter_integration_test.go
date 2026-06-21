package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	
	myredis "github.com/thnxvlad/oplati/internal/storages/redis"
)

func TestStorage_Allow(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := testredis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("* Ready to accept connections"),
		),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	endpoint, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	options, err := redis.ParseURL(endpoint)
	if err != nil {
		t.Fatalf("failed to parse redis url: %s", err)
	}
	client := redis.NewClient(options)
	defer client.Close()

	storage := myredis.New(client)


	key := "test_ip_127.0.0.1"
	limit := int64(2)
	window := 1 * time.Second

	ok, err := storage.Allow(ctx, key, limit, window)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected first request to be allowed")
	}

	ok, _ = storage.Allow(ctx, key, limit, window)
	if !ok {
		t.Error("expected second request to be allowed")
	}

	ok, _ = storage.Allow(ctx, key, limit, window)
	if ok {
		t.Error("expected third request to be blocked (limit exceeded)")
	}

	time.Sleep(window + 2*time.Second)

	ok, _ = storage.Allow(ctx, key, limit, window)
	if !ok {
		t.Error("expected request to be allowed after window expired")
	}
}