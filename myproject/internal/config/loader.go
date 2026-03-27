package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Load 从 YAML 文件加载配置，并用环境变量覆盖
func Load(path string) *Config {
	cfg := &Config{
		Server:   ServerConfig{Port: "8080", Mode: "development"},
		GRPC:     GRPCConfig{Port: "9090"},
		JWT:      JWTConfig{Secret: "change-me-in-production", Expiry: "24h"},
		Database: DatabaseConfig{Driver: "memory"},
		Redis:    RedisConfig{Addr: "localhost:6379"},
		Log:      LogConfig{Level: "debug", Format: "text"},
		MQ:       MQConfig{Driver: "memory"},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("config file not found (%s), using defaults\n", path)
		ApplyEnvOverrides(cfg)
		return cfg
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	ApplyEnvOverrides(cfg)
	return cfg
}