package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment  string
	Port         int
	DatabasePath string
}

func Load() Config {
	return Config{
		Environment:  getEnv("APP_ENV", "development"),
		Port:         getEnvInt("APP_PORT", 8080),
		DatabasePath: getEnv("DATABASE_PATH", "./data/app.db"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
