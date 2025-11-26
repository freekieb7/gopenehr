package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

const SYSTEM_ID_GOPENEHR = "gopenehr"
const NAMESPACE_LOCAL = "local"
const API_KEY_HEADER = "X-API-Key"
const TARGET_MIGRATION_VERSION uint64 = 20251113195000

var SystemUserID = uuid.MustParse("9f5c8bd8-4c88-43f0-90e6-baadf00dc0de")

var (
	Version = "dev"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

func (e Environment) IsValid() bool {
	switch e {
	case Development, Production:
		return true
	default:
		return false
	}
}

type Config struct {
	Port        string
	Version     string
	LogLevel    slog.Level
	DatabaseURL string
	APIKey      string
}

func (c *Config) Load() error {
	c.Version = Version

	logLevel, err := getEnvLogLevel("LOG_LEVEL", slog.LevelInfo, false)
	if err != nil {
		return fmt.Errorf("invalid LOG_LEVEL value")
	}
	c.LogLevel = logLevel

	port, err := getEnvString("APP_PORT", "3000", false)
	if err != nil {
		return err
	}
	c.Port = port

	dbURL, err := getEnvString("DATABASE_URL", "", true)
	if err != nil {
		return err
	}
	c.DatabaseURL = dbURL

	apiKey, err := getEnvString("API_KEY", "", false)
	if err != nil {
		return err
	}
	c.APIKey = apiKey

	return nil
}

func getEnvString(key string, defaultValue string, required bool) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		if required {
			return "", fmt.Errorf("environment variable %s is required", key)
		}
		return defaultValue, nil
	}
	return value, nil
}

func getEnvLogLevel(key string, defaultValue slog.Level, required bool) (slog.Level, error) {
	levelStr, exists := os.LookupEnv(key)
	if !exists {
		if required {
			return slog.LevelInfo, fmt.Errorf("environment variable %s is required", key)
		}
		return defaultValue, nil
	}
	var level slog.Level
	err := level.UnmarshalText([]byte(levelStr))
	if err != nil {
		return slog.LevelInfo, fmt.Errorf("environment variable %s has invalid value: %s", key, levelStr)
	}
	return level, nil
}
