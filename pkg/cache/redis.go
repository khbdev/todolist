package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
	ttl    time.Duration
}


func NewCache(client *redis.Client) *Cache {
	return &Cache{
		Client: client,
		ttl:    1 * time.Minute, 
	}
}


func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}


func (c *Cache) Set(ctx context.Context, key string, value string, ttl ...time.Duration) error {
	expire := c.ttl
	if len(ttl) > 0 {
		expire = ttl[0]
	}
	return c.Client.Set(ctx, key, value, expire).Err()
}


func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}
