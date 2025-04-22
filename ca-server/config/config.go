package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	ServerPort     int // Regular HTTP port
	TLSServerPort  int // TLS server port
	MTLSServerPort int // mTLS server port
	Mode           string
	LogLevel       string
	// TLS configuration
	TLSEnabled  bool
	TLSCertPath string
	TLSKeyPath  string
	// mTLS configuration
	MTLSEnabled      bool
	ClientCACertPath string
}

// New creates a new Config with values from environment
func New() *Config {
	return &Config{
		ServerPort:     getEnvAsInt("SERVER_PORT", 8080),
		TLSServerPort:  getEnvAsInt("TLS_SERVER_PORT", 8080),
		MTLSServerPort: getEnvAsInt("MTLS_SERVER_PORT", 8443),
		Mode:           getEnv("GIN_MODE", "debug"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		// TLS configuration with defaults
		TLSEnabled:  getEnvAsBool("TLS_ENABLED", false),
		TLSCertPath: getEnv("TLS_CERT_PATH", "server/certs/cert.pem"),
		TLSKeyPath:  getEnv("TLS_KEY_PATH", "server/certs/key.pem"),
		// mTLS configuration
		MTLSEnabled:      getEnvAsBool("MTLS_ENABLED", false),
		ClientCACertPath: getEnv("CLIENT_CA_CERT_PATH", "cert.pem"),
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

// Helper to get env as boolean with fallback
func getEnvAsBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}
