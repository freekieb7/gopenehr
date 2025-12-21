package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const SYSTEM_ID_GOPENEHR = "e6d14bbd-2a0c-474a-9964-11f0bfbe36bd"
const NAMESPACE_LOCAL = "local"
const API_KEY_HEADER = "X-API-Key"

var SystemUserID = uuid.MustParse(SYSTEM_ID_GOPENEHR)

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
	Name                string
	Port                string
	Version             string
	DatabaseURL         string
	LogLevel            slog.Level
	APIKey              string
	OAuthTrustedIssuers []string
	OAuthAudience       string
	OtelEndpoint        string
	OtelInsecure        bool
	KafkaBrokers        []string
	CapEHRs             int
}

func NewSettings() *Settings {
	return &Settings{
		Name:    SYSTEM_ID_GOPENEHR,
		Version: Version,
	}
}

func (s *Settings) Load() error {
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

	OtelEndpoint, err := getEnvString("OTEL_ENDPOINT", "", false)
	if err != nil {
		return err
	}
	s.OtelEndpoint = OtelEndpoint

	otelInsecure, err := getEnvBool("OTEL_INSECURE", false, false)
	if err != nil {
		return err
	}
	s.OtelInsecure = otelInsecure

	kafkaBrokers, err := getEnvString("KAFKA_BROKERS", "", false)
	if err != nil {
		return err
	}
	for broker := range strings.SplitSeq(kafkaBrokers, ",") {
		broker = strings.TrimSpace(broker)
		if broker == "" {
			continue
		}
		s.KafkaBrokers = append(s.KafkaBrokers, broker)
	}

	capEHRs, err := getEnvUint64("CAP_EHRS", 0, false)
	if err != nil {
		return err
	}
	s.CapEHRs = int(capEHRs)
	return nil
}

func (s *Settings) OAuthEnabled() bool {
	return len(s.OAuthTrustedIssuers) > 0 && s.OAuthAudience != ""
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

func getEnvBool(key string, defaultValue bool, required bool) (bool, error) {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		if required {
			return false, fmt.Errorf("environment variable %s is required", key)
		}
		return defaultValue, nil
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, fmt.Errorf("environment variable %s has invalid value: %s", key, valueStr)
	}
	return value, nil
}

func getEnvUint64(key string, defaultValue uint64, required bool) (uint64, error) {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		if required {
			return 0, fmt.Errorf("environment variable %s is required", key)
		}
		return defaultValue, nil
	}
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s has invalid value: %s", key, valueStr)
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
