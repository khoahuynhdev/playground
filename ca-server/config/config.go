package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	ServerPort int
	Mode       string
	LogLevel   string
}

// New creates a new Config with values from environment
func New() *Config {
	return &Config{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		Mode:       getEnv("GIN_MODE", "debug"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}
}

// Helper to get env with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper to get env as integer with fallback
func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
