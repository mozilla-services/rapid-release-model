package config

import (
	"os"
	"strings"
)

// Prefix for application specific environment variables
const envPrefix string = "RRM_METRICS"

func FromEnv(p ...string) string {
	parts := append([]string{envPrefix}, p...)
	key := strings.Join(parts, "__")
	return os.Getenv(key)
}
