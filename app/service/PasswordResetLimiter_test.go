package service

import (
	"testing"
	"time"
)

func newTestResetLimiter(now *time.Time) *PasswordResetLimiter {
	return &PasswordResetLimiter{requests: make(map[string]passwordResetWindow), now: func() time.Time { return *now }}
}

func TestPasswordResetLimiterPerAddress(t *testing.T) {
	now := time.Now()
	limiter := newTestResetLimiter(&now)
	for i := 0; i < defaultPasswordResetLimit; i++ {
		if !limiter.Allow("User@Example.test") {
			t.Fatalf("request %d blocked too early", i+1)
		}
	}
	if limiter.Allow("  user@example.test ") {
		t.Fatal("address casing/whitespace bypassed the limit")
	}
	if !limiter.Allow("other@example.test") {
		t.Fatal("unrelated address inherited the limit")
	}
	if limiter.Allow("") {
		t.Fatal("empty address must not be accepted")
	}
}

func TestPasswordResetLimiterRecoversAfterWindow(t *testing.T) {
	now := time.Now()
	limiter := newTestResetLimiter(&now)
	for i := 0; i < defaultPasswordResetLimit; i++ {
		limiter.Allow("user@example.test")
	}
	if limiter.Allow("user@example.test") {
		t.Fatal("expected the address to be limited")
	}
	// Denied attempts must not extend the window.
	now = now.Add(defaultPasswordResetWindow - time.Minute)
	if limiter.Allow("user@example.test") {
		t.Fatal("limit expired before the window elapsed")
	}
	now = now.Add(2 * time.Minute)
	if !limiter.Allow("user@example.test") {
		t.Fatal("address did not recover after the window")
	}
}

func TestFindPwdThrottlesBeforeLookingUpTheAccount(t *testing.T) {
	// An empty address is refused by the limiter, so the request must look
	// successful without touching the database.
	ok, msg := (&PwdService{}).FindPwd("")
	if !ok || msg != "" {
		t.Fatalf("throttled request must stay indistinguishable: ok=%v msg=%q", ok, msg)
	}
}
