package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrMiss = errors.New("cache miss")

type Store struct{ client *redis.Client }

var rateLimitScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
return current
`)

func NewRedis(addr string) *Store {
	return &Store{client: redis.NewClient(&redis.Options{Addr: addr, DialTimeout: 2 * time.Second, ReadTimeout: time.Second, WriteTimeout: time.Second})}
}

func (s *Store) Ping(ctx context.Context) error { return s.client.Ping(ctx).Err() }
func (s *Store) Close() error                   { return s.client.Close() }

func (s *Store) GetJSON(ctx context.Context, key string, target any) error {
	value, err := s.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return ErrMiss
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(value, target)
}

func (s *Store) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, encoded, ttl).Err()
}

func (s *Store) DeletePrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	for {
		keys, next, err := s.client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err = s.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}

func (s *Store) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	count, err := rateLimitScript.Run(ctx, s.client, []string{key}, window.Milliseconds()).Int64()
	if err != nil {
		return true, err
	}
	return count <= limit, nil
}
