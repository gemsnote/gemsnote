// Package version exposes the canonical versions embedded from the adjacent
// version files. Release and container tooling read the same files.
package version

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
)

//go:embed VERSION
var currentSource string

//go:embed API_VERSION
var apiSource string

var semanticVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

// Current is the server application and container release version.
var Current = mustVersion("VERSION", currentSource)

// API is the current API2 protocol version.
var API = mustVersion("API_VERSION", apiSource)

// MinimumClient remains empty while older clients are still accepted.
const MinimumClient = ""

func mustVersion(name, source string) string {
	value := strings.TrimSpace(source)
	if !semanticVersion.MatchString(value) {
		panic(fmt.Sprintf("invalid %s %q: expected MAJOR.MINOR.PATCH", name, value))
	}
	return value
}
