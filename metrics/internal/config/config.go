package config

import (
	"os"
	"strings"
)

// Prefix for application specific environment variables
const envPrefix string = "RRM_METRICS"

func FromEnv(p ...string) string {
	p = append(p, envPrefix)
	k := strings.Join(p, "__")
	return os.Getenv(k)
}
