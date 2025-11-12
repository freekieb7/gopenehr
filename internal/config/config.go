package internal

import (
	"fmt"
	"os"
)

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
	Version     string
	Environment Environment
	DatabaseURL string
}

func (c *Config) Load() error {
	c.Version = Version

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
