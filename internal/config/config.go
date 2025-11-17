package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// Config contains runtime configuration for the API.
type Config struct {
	ServerPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
	JWTTTL    time.Duration
}

// Load reads configuration from environment variables, applying sensible defaults.
func Load() (Config, error) {
	cfg := Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "todo"),
		DBPassword: getEnv("DB_PASSWORD", "todo"),
		DBName:     getEnv("DB_NAME", "todo"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET must be provided")
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, errors.New("JWT_SECRET must be at least 32 characters")
	}

	cfg.JWTTTL = 15 * time.Minute
	if ttlStr := os.Getenv("JWT_TTL_MINUTES"); ttlStr != "" {
		minutes, err := strconv.Atoi(ttlStr)
		if err != nil || minutes <= 0 {
			return Config{}, errors.New("JWT_TTL_MINUTES must be a positive integer")
		}
		cfg.JWTTTL = time.Duration(minutes) * time.Minute
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
