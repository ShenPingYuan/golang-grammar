package retry

import (
	"context"
	"math"
	"time"
)

type Config struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

func DefaultConfig() Config {
	return Config{
		MaxRetries: 3,
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   5 * time.Second,
	}
}

// Do 带指数退避的重试
func Do(ctx context.Context, cfg Config, fn func() error) error {
	var lastErr error
	for i := 0; i <= cfg.MaxRetries; i++ {
		if err := fn(); err != nil {
			lastErr = err
			if i == cfg.MaxRetries {
				break
			}
			delay := time.Duration(math.Pow(2, float64(i))) * cfg.BaseDelay
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			continue
		}
		return nil
	}
	return lastErr
}