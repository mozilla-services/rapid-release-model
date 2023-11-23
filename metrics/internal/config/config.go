package config

import (
	"fmt"
	"os"
	"strings"
)

// Prefix for application specific environment variables
const envPrefix string = "RRM_METRICS"

// EnvKey returns the full key for the given key parts.
func EnvKey(p ...string) string {
	parts := append([]string{envPrefix}, p...)
	return strings.Join(parts, "__")
}

// ReadFromEnv reads the environment variable for the given parts.
func ReadFromEnv(p ...string) string {
	key := EnvKey(p...)
	return os.Getenv(key)
}

// ReadFromEnvE reads the environment variable for the given parts and returns
// an error if the environment variable is not set or if the value is empty.
func ReadFromEnvE(p ...string) (string, error) {
	key := EnvKey(p...)
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("Required environment variable %v not set.", key)
	}
	return val, nil
}
