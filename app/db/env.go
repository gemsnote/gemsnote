package db

import (
	"os"
	"strings"
)

// Environment override lets deployments keep database credentials out of the
// configuration file that is mounted into the container.
func envValue(name string) (string, bool) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return "", false
	}
	return value, true
}
