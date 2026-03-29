package database

import (
	"log/slog"
)

type Redis struct {
	Addr     string
	Password string
	DB       int
}

// NewRedis 创建 Redis 连接
// 生产环境请使用 github.com/redis/go-redis/v9
func NewRedis(addr, password string, db int) *Redis {
	slog.Info("Redis configured", "addr", addr)
	return &Redis{Addr: addr, Password: password, DB: db}
}