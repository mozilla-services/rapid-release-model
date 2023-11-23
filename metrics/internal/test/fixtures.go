package test

import (
	"os"
	"path/filepath"
)

// Load the given file from the fixtures directory
func LoadFixture(p ...string) ([]byte, error) {
	parts := append([]string{"./fixtures"}, p...)
	return os.ReadFile(filepath.Join(parts...))
}

// Load the fixture at the given location
func NewFixture(p ...string) func() ([]byte, error) {
	return func() ([]byte, error) {
		return LoadFixture(p...)
	}
}
