package app

import (
	"context"
	"encoding/json"
	"time"

	"tiny-url/internal/domain"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, code string) (*domain.URL, error)
	Set(ctx context.Context, code string, url *domain.URL, ttl time.Duration) error
	Delete(ctx context.Context, code string) error
}

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{rdb: rdb}
}

func (c *RedisCache) Get(ctx context.Context, code string) (*domain.URL, error) {
	val, err := c.rdb.Get(ctx, code).Result()
	if err == redis.Nil {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var u domain.URL
	if err := json.Unmarshal([]byte(val), &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *RedisCache) Set(ctx context.Context, code string, url *domain.URL, ttl time.Duration) error {
	b, err := json.Marshal(url)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, code, b, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, code string) error {
	return c.rdb.Del(ctx, code).Err()
}
