package service

import (
	"testing"
	"time"
)

func newTestLoginAttempts(now time.Time) *LoginAttemptService {
	return &LoginAttemptService{attempts: make(map[string]loginAttempt), now: func() time.Time { return now }}
}

func TestLoginAttemptsBlocksAfterLimit(t *testing.T) {
	now := time.Now()
	throttle := newTestLoginAttempts(now)
	for i := 0; i < defaultLoginFailureLimit; i++ {
		if allowed, _ := throttle.Allow("victim@example.test", "203.0.113.9"); !allowed {
			t.Fatalf("attempt %d blocked too early", i+1)
		}
		throttle.Failure("victim@example.test", "203.0.113.9")
	}
	allowed, retry := throttle.Allow("victim@example.test", "203.0.113.9")
	if allowed {
		t.Fatal("expected the account/address pair to be blocked")
	}
	if retry <= 0 || retry > defaultLoginFailureWindow {
		t.Fatalf("unexpected retry duration %v", retry)
	}
	// Another source address keeps its own budget.
	if allowed, _ := throttle.Allow("victim@example.test", "198.51.100.7"); !allowed {
		t.Fatal("unrelated address must not inherit the block")
	}
}

func TestLoginAttemptsExpireAndReset(t *testing.T) {
	now := time.Now()
	throttle := newTestLoginAttempts(now)
	for i := 0; i < defaultLoginFailureLimit; i++ {
		throttle.Failure("user@example.test", "203.0.113.9")
	}
	if allowed, _ := throttle.Allow("user@example.test", "203.0.113.9"); allowed {
		t.Fatal("expected block")
	}
	now = now.Add(defaultLoginFailureWindow + time.Second)
	throttle.now = func() time.Time { return now }
	if allowed, _ := throttle.Allow("user@example.test", "203.0.113.9"); !allowed {
		t.Fatal("block should expire after the window")
	}

	throttle.Failure("user@example.test", "203.0.113.9")
	throttle.Success("user@example.test", "203.0.113.9")
	if allowed, _ := throttle.Allow("user@example.test", "203.0.113.9"); !allowed {
		t.Fatal("success should clear the account counter")
	}
}

func TestLoginAttemptsLimitSprayingFromOneAddress(t *testing.T) {
	now := time.Now()
	throttle := newTestLoginAttempts(now)
	targets := defaultLoginFailureLimit * ipFailureFactor
	for i := 0; i < targets; i++ {
		throttle.Failure("user"+string(rune('a'+i%26))+"@example.test", "203.0.113.9")
	}
	allowed, _ := throttle.Allow("fresh@example.test", "203.0.113.9")
	if allowed {
		t.Fatal("one address must not be able to spray many accounts")
	}
	if allowed, _ := throttle.Allow("fresh@example.test", "198.51.100.7"); !allowed {
		t.Fatal("another address must stay usable")
	}
}
