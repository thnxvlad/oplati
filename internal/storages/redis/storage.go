package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	db *redis.Client
}

func New(client *redis.Client) *Storage {
	return &Storage{db: client}
}

const script = `
	local current = redis.call("INCR", KEYS[1])
	if current == 1 then
		redis.call("EXPIRE", KEYS[1], ARGV[1])
	end
	return current
`

func (s *Storage) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	val, err := s.db.Eval(ctx, script, []string{key}, int(window.Seconds())).Int64()
	if err != nil {
		return false, err
	}

	return val <= limit, nil
}
