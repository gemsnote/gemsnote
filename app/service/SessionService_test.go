package service

import (
	"encoding/base64"
	"testing"
	"time"
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

func TestStoredSessionKeyHidesCredential(t *testing.T) {
	token, err := NewAPIToken()
	if err != nil {
		t.Fatal(err)
	}
	key := storedSessionKey(token)
	if key == token || len(key) == 0 {
		t.Fatalf("session key must not be the raw credential: %q", key)
	}
	if storedSessionKey(token) != key {
		t.Fatal("session key must be deterministic")
	}
	if storedSessionKey(key) == key {
		t.Fatal("a stored key must not be usable as a raw credential")
	}
	if storedSessionKey("") != "" {
		t.Fatal("empty credential must not map to a stored key")
	}
}

func TestSessionExpiredHonoursIdleAndAbsoluteLimits(t *testing.T) {
	now := time.Now()
	if !sessionExpired(time.Time{}, now, time.Hour, time.Hour) {
		t.Fatal("sessions without timestamps must be rejected")
	}
	if !sessionExpired(now.Add(-26*time.Hour), now.Add(-2*time.Hour), time.Hour, 48*time.Hour) {
		t.Fatal("idle session should expire")
	}
	if !sessionExpired(now.Add(-25*time.Hour), now, 24*time.Hour, 24*time.Hour) {
		t.Fatal("absolute lifetime should expire an active session")
	}
	if sessionExpired(now.Add(-time.Hour), now, 3*time.Hour, 24*time.Hour) {
		t.Fatal("fresh session should stay valid")
	}
}
