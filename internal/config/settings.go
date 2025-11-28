package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/google/uuid"
)

const SYSTEM_ID_GOPENEHR = "gopenehr"
const NAMESPACE_LOCAL = "local"
const API_KEY_HEADER = "X-API-Key"
const TARGET_MIGRATION_VERSION uint64 = 20251113195000

var SystemUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

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

type Settings struct {
	Port                string
	Version             string
	DatabaseURL         string
	LogLevel            slog.Level
	APIKey              string
	OAuthTrustedIssuers []string
	OAuthAudience       string
}

func NewSettings() Settings {
	return Settings{}
}

func (s *Settings) Load() error {
	s.Version = Version

	port, err := getEnvString("APP_PORT", "3000", false)
	if err != nil {
		return err
	}
	s.Port = port

	dbURL, err := getEnvString("DATABASE_URL", "", true)
	if err != nil {
		return err
	}
	s.DatabaseURL = dbURL

	logLevel, err := getEnvLogLevel("LOG_LEVEL", slog.LevelInfo, false)
	if err != nil {
		return err
	}
	s.LogLevel = logLevel

	apiKey, err := getEnvString("API_KEY", "", false)
	if err != nil {
		return err
	}
	s.APIKey = apiKey

	trustedIssuersStr, err := getEnvString("OAUTH_TRUSTED_ISSUERS", "", false)
	if err != nil {
		return err
	}
	s.OAuthTrustedIssuers = []string{}
	for issuer := range strings.SplitSeq(trustedIssuersStr, ",") {
		issuer = strings.TrimSpace(issuer)
		if issuer == "" {
			continue
		}
		s.OAuthTrustedIssuers = append(s.OAuthTrustedIssuers, issuer)
	}

	audience, err := getEnvString("OAUTH_AUDIENCE", "", false)
	if err != nil {
		return err
	}
	s.OAuthAudience = audience

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
