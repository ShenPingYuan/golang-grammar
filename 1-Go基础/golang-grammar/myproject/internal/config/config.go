package config

import "time"

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	JWT      JWTConfig      `yaml:"jwt"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Log      LogConfig      `yaml:"log"`
	MQ       MQConfig       `yaml:"mq"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type GRPCConfig struct {
	Port string `yaml:"port"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expiry string `yaml:"expiry"`
}

func (c JWTConfig) ExpiryDuration() time.Duration {
	d, err := time.ParseDuration(c.Expiry)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type MQConfig struct {
	Driver  string   `yaml:"driver"`
	Brokers []string `yaml:"brokers"`
}