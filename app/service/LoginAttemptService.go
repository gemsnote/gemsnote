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
	defaultLoginFailureLimit  = 5
	defaultLoginFailureWindow = 15 * time.Minute
	// A single source address may try many accounts, so its aggregate budget is
	// larger than one account/address pair while still stopping password sprays.
	ipFailureFactor = 10
)

type loginAttempt struct {
	failures int
	first    time.Time
	last     time.Time
}

// LoginAttemptService provides a bounded, process-local first line of defence
// against password guessing. Keys are hashed so usernames and client addresses
// are not retained in memory. Deployments with multiple application replicas
// should additionally enforce a shared rate limit at their reverse proxy.
type LoginAttemptService struct {
	mu       sync.Mutex
	attempts map[string]loginAttempt
	now      func() time.Time
}

var LoginAttempts = &LoginAttemptService{attempts: make(map[string]loginAttempt), now: time.Now}

func loginAttemptHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func loginAccountKey(identity, clientIP string) string {
	return "user:" + loginAttemptHash(strings.ToLower(strings.TrimSpace(identity))+"\x00"+strings.TrimSpace(clientIP))
}

func loginIPKey(clientIP string) string {
	return "ip:" + loginAttemptHash(strings.TrimSpace(clientIP))
}

func loginFailureLimit(key string) int {
	limit := defaultLoginFailureLimit
	if revel.Config != nil {
		if configured := revel.Config.IntDefault("login.failure.limit", defaultLoginFailureLimit); configured > 0 {
			limit = configured
		}
	}
	if strings.HasPrefix(key, "ip:") {
		return limit * ipFailureFactor
	}
	return limit
}

func loginFailureWindow() time.Duration {
	return configuredDuration("login.failure.window", defaultLoginFailureWindow)
}

// Allow reports whether the (identity, client IP) pair may attempt a login and,
// when it may not, how long the caller should wait before retrying.
func (s *LoginAttemptService) Allow(identity, clientIP string) (bool, time.Duration) {
	window := loginFailureWindow()
	now := s.now()
	s.mu.Lock()
	defer s.mu.Unlock()
	blocked := false
	var retry time.Duration
	for _, key := range []string{loginAccountKey(identity, clientIP), loginIPKey(clientIP)} {
		attempt, ok := s.attempts[key]
		if !ok {
			continue
		}
		if now.Sub(attempt.first) >= window {
			delete(s.attempts, key)
			continue
		}
		if attempt.failures >= loginFailureLimit(key) {
			blocked = true
			if wait := window - now.Sub(attempt.first); wait > retry {
				retry = wait
			}
		}
	}
	return !blocked, retry
}

func (s *LoginAttemptService) Failure(identity, clientIP string) {
	window := loginFailureWindow()
	now := s.now()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, key := range []string{loginAccountKey(identity, clientIP), loginIPKey(clientIP)} {
		attempt := s.attempts[key]
		if attempt.first.IsZero() || now.Sub(attempt.first) >= window {
			attempt = loginAttempt{first: now}
		}
		attempt.failures++
		attempt.last = now
		s.attempts[key] = attempt
	}
	// Keep attacker-controlled identities from growing memory without bound.
	if len(s.attempts) > 10000 {
		for candidate, value := range s.attempts {
			if now.Sub(value.last) >= window {
				delete(s.attempts, candidate)
			}
		}
	}
}

func (s *LoginAttemptService) Success(identity, clientIP string) {
	s.mu.Lock()
	delete(s.attempts, loginAccountKey(identity, clientIP))
	s.mu.Unlock()
}
