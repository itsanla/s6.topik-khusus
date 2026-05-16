package config

import (
	"os"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Port          string
	AppEnv        string
}

func Load() *Config {
	return &Config{
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
		Port:          getEnv("PORT", "8080"),
		AppEnv:        getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
