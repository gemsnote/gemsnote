package service

import (
	"encoding/base64"
	"testing"
)

func TestNewAPITokenUsesFullRandomPayload(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 128; i++ {
		token, err := NewAPIToken()
		if err != nil {
			t.Fatal(err)
		}
		raw, err := base64.RawURLEncoding.DecodeString(token)
		if err != nil || len(raw) != apiTokenBytes {
			t.Fatalf("invalid API token %q: bytes=%d err=%v", token, len(raw), err)
		}
		if seen[token] {
			t.Fatalf("duplicate API token generated: %q", token)
		}
		seen[token] = true
	}
}
