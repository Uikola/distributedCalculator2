package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c Cache) Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, val, ttl).Err()
}
