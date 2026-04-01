package cache

import (
	"context"
	"log/slog"
	"time"
)

// RedisCache 基于 Redis 的缓存实现
// 生产环境请使用 github.com/redis/go-redis/v9
type RedisCache struct {
	addr string
}

func NewRedis(addr string) Cache {
	slog.Info("RedisCache configured", "addr", addr)
	return &RedisCache{addr: addr}
}

func (c *RedisCache) Get(_ context.Context, key string) (string, error) {
	// TODO: 使用 go-redis 客户端实现
	slog.Debug("RedisCache.Get", "key", key)
	return "", nil
}

func (c *RedisCache) Set(_ context.Context, key, value string, ttl time.Duration) error {
	slog.Debug("RedisCache.Set", "key", key, "ttl", ttl)
	return nil
}

func (c *RedisCache) Delete(_ context.Context, key string) error {
	slog.Debug("RedisCache.Delete", "key", key)
	return nil
}

func (c *RedisCache) Exists(_ context.Context, key string) (bool, error) {
	slog.Debug("RedisCache.Exists", "key", key)
	return false, nil
}