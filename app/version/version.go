// Package version is the single source of truth for Gemsnote versions.
package version

// Current is the server application release version. Keep this declaration on
// one line because release tooling reads it without compiling the package.
const Current = "1.1.0"

// API is the current API2 protocol version.
const API = "1.1.0"

// MinimumClient remains empty while older clients are still accepted.
const MinimumClient = ""
