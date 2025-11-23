package config

import (
	"fmt"
	"os"
)

const SYSTEM_ID_GOPENEHR = "gopenehr"
const NAMESPACE_LOCAL = "local"

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

func (e Environment) IsProduction() bool {
	return e == Production
}

type Config struct {
	Host          string
	Version       string
	Environment   Environment
	DatabaseURL   string
	MigrationsDir string
}

func (c *Config) Load() error {
	c.Version = Version

	host, err := getEnvString("APP_HOST", "http://localhost:3000", false)
	if err != nil {
		return err
	}
	c.Host = host

	env, err := getEnvEnvironment("APP_ENV", Production, false)
	if err != nil {
		return err
	}
	c.Environment = env

	dbURL, err := getEnvString("DATABASE_URL", "", true)
	if err != nil {
		return err
	}
	c.DatabaseURL = dbURL

	migrationsDir, err := getEnvString("MIGRATIONS_DIR", "./migrations", false)
	if err != nil {
		return err
	}
	c.MigrationsDir = migrationsDir

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

func getEnvEnvironment(key string, defaultValue Environment, required bool) (Environment, error) {
	env, exists := os.LookupEnv(key)
	if !exists {
		if required {
			return "", fmt.Errorf("environment variable %s is required", key)
		}
		return defaultValue, nil
	}
	envValue := Environment(env)
	if !envValue.IsValid() {
		return "", fmt.Errorf("environment variable %s has invalid value: %s", key, env)
	}
	return envValue, nil
}
