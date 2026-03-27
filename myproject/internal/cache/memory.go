package cache

import (
	"context"
	"sync"
	"time"

	apperrors "myproject/internal/errors"
)

type cacheItem struct {
	value     string
	expiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func NewMemory() Cache {
	mc := &MemoryCache{items: make(map[string]cacheItem)}
	go mc.cleanup()
	return mc
}

func (c *MemoryCache) Get(_ context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok || time.Now().After(item.expiresAt) {
		return "", apperrors.ErrNotFound
	}
	return item.value, nil
}

func (c *MemoryCache) Set(_ context.Context, key, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{value: value, expiresAt: time.Now().Add(ttl)}
	return nil
}

func (c *MemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *MemoryCache) Exists(_ context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok || time.Now().After(item.expiresAt) {
		return false, nil
	}
	return true, nil
}

func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.items {
			if now.After(v.expiresAt) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}