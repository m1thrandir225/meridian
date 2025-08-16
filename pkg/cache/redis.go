package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client  *redis.Client
	metrics *CacheMetrics
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client:  client,
		metrics: &CacheMetrics{},
	}
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) GetWithMetrics(ctx context.Context, key string, dest interface{}) (bool, error) {
	err := c.Get(ctx, key, dest)
	if err != nil {
		c.metrics.IncrementMisses()
		return false, err
	}
	c.metrics.IncrementHits()
	return true, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, json, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	return nil
}

func (c *RedisCache) GetMetrics() *CacheMetrics {
	return c.metrics
}
