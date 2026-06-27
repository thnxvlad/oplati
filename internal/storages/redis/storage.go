package redislimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	db *redis.Client
}

func New(pool *redis.Client) *Storage {
	return &Storage{db: pool}
}

func (s *Storage) IncrementWithTTL(ctx context.Context, key string, limit int64, ttl time.Duration) (bool, error) {
	value, err := s.db.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if value == 1 {
		err = s.db.Expire(ctx, key, ttl).Err()
		if err != nil {
			return false, err
		}
	}

	if value > limit {
		return false, nil
	}

	return true, nil
}
