package lea

import (
	"testing"
	"time"
)

const (
	testImageUser = "507f1f77bcf86cd799439011"
	testImageFile = "0123456789abcdef01234567"
)

func TestImageTokenRoundTrip(t *testing.T) {
	now := time.Now()
	secret := DeriveImageSecret("0123456789abcdef0123456789abcdef")
	if secret == "" {
		t.Fatal("empty derived image secret")
	}
	token := ImageToken(secret, testImageUser, testImageFile, now, 10*time.Minute)
	if token == "" {
		t.Fatal("empty image token")
	}
	if !ValidImageToken(secret, testImageUser, testImageFile, token, now, 10*time.Minute) {
		t.Fatal("valid capability rejected")
	}
	if ValidImageToken(secret, testImageUser, "0123456789abcdef01234568", token, now, 10*time.Minute) {
		t.Fatal("capability accepted for another file")
	}
	if ValidImageToken(secret, "507f1f77bcf86cd799439012", testImageFile, token, now, 10*time.Minute) {
		t.Fatal("capability accepted for another user")
	}
	if ValidImageToken(DeriveImageSecret("another-secret"), testImageUser, testImageFile, token, now, 10*time.Minute) {
		t.Fatal("capability accepted with another key")
	}
	if ValidImageToken(secret, testImageUser, testImageFile, token, now.Add(11*time.Minute), 10*time.Minute) {
		t.Fatal("expired capability accepted")
	}
	// A token cannot outlive the configured lifetime even if its timestamp says so.
	long := ImageToken(secret, testImageUser, testImageFile, now, time.Hour)
	if ValidImageToken(secret, testImageUser, testImageFile, long, now, 10*time.Minute) {
		t.Fatal("capability with excessive lifetime accepted")
	}
	for _, broken := range []string{"", "x", "1.2.3", token + "0", "." + token} {
		if ValidImageToken(secret, testImageUser, testImageFile, broken, now, 10*time.Minute) {
			t.Fatalf("malformed capability accepted: %q", broken)
		}
	}
}

func TestImageSecretDiffersFromOtherKeys(t *testing.T) {
	appSecret := "0123456789abcdef0123456789abcdef"
	if DeriveImageSecret(appSecret) == DerivePDFSecret(appSecret) {
		t.Fatal("image and PDF keys must be independent")
	}
	if DeriveImageSecret(appSecret) == appSecret {
		t.Fatal("image key must be derived, not the raw secret")
	}
	if DeriveImageSecret("") != "" {
		t.Fatal("empty app secret must not derive a key")
	}
}
