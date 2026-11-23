package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	"github.com/revel/revel"
)

const (
	defaultPasswordResetLimit  = 3
	defaultPasswordResetWindow = 15 * time.Minute
)

type passwordResetWindow struct {
	count int
	first time.Time
}

// PasswordResetLimiter bounds how many recovery emails one address may trigger
// in a window, so the public reset endpoint cannot be used to flood a mailbox.
// Keys are hashed; addresses are never retained in memory.
type PasswordResetLimiter struct {
	mu       sync.Mutex
	requests map[string]passwordResetWindow
	now      func() time.Time
}

var PasswordResetRequests = &PasswordResetLimiter{requests: make(map[string]passwordResetWindow), now: time.Now}

func passwordResetKey(email string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(email))))
	return hex.EncodeToString(sum[:])
}

func passwordResetLimit() int {
	if revel.Config != nil {
		if configured := revel.Config.IntDefault("password.reset.limit", defaultPasswordResetLimit); configured > 0 {
			return configured
		}
	}
	return defaultPasswordResetLimit
}

func passwordResetWindowDuration() time.Duration {
	return configuredDuration("password.reset.window", defaultPasswordResetWindow)
}

// Allow records a reset request for the address and reports whether an email
// may be sent. Requests over the limit are counted as denied but do not extend
// the window, so a blocked address recovers on its own.
func (s *PasswordResetLimiter) Allow(email string) bool {
	if strings.TrimSpace(email) == "" {
		return false
	}
	key := passwordResetKey(email)
	window := passwordResetWindowDuration()
	now := s.now()
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := s.requests[key]
	if entry.first.IsZero() || now.Sub(entry.first) >= window {
		entry = passwordResetWindow{first: now}
	}
	if entry.count >= passwordResetLimit() {
		return false
	}
	entry.count++
	s.requests[key] = entry
	// Keep memory bounded if an attacker cycles through addresses.
	if len(s.requests) > 10000 {
		for candidate, value := range s.requests {
			if now.Sub(value.first) >= window {
				delete(s.requests, candidate)
			}
		}
	}
	return true
}
