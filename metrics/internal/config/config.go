package config

import (
	"os"
	"strings"
)

// Prefix for application specific environment variables
const envPrefix string = "RRM_METRICS"

func Key(p ...string) string {
	parts := append([]string{envPrefix}, p...)
	return strings.Join(parts, "__")
}

func FromEnv(p ...string) string {
	key := Key(p...)
	return os.Getenv(key)
}
