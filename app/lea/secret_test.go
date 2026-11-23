package lea

import (
	"testing"
	"time"
)

func TestValidateAppSecret(t *testing.T) {
	cases := []struct {
		name       string
		secret     string
		production bool
		ok         bool
	}{
		{"empty", "", false, false},
		{"empty prod", "", true, false},
		{"shipped default", ShippedAppSecret, false, false},
		{"shipped default prod", ShippedAppSecret, true, false},
		{"short prod", "dev-only-secret", true, false},
		{"short dev", "dev-only-secret", false, true},
		{"long prod", "0123456789abcdef0123456789abcdef", true, true},
	}
	for _, tc := range cases {
		err := ValidateAppSecret(tc.secret, tc.production)
		if (err == nil) != tc.ok {
			t.Fatalf("%s: got err=%v, want ok=%v", tc.name, err, tc.ok)
		}
	}
}

func TestPDFCapabilityDoesNotUseAppSecretDirectly(t *testing.T) {
	appSecret := "0123456789abcdef0123456789abcdef"
	pdfSecret := DerivePDFSecret(appSecret)
	if pdfSecret == "" || pdfSecret == appSecret {
		t.Fatal("PDF key must be derived, not the raw app secret")
	}
	now := time.Now()
	token := PDFToken(pdfSecret, "5b1f0b7c9b1e4b0a8a2f3c4d", now)
	if token == "" {
		t.Fatal("empty token")
	}
	if !ValidPDFToken(pdfSecret, "5b1f0b7c9b1e4b0a8a2f3c4d", token, now) {
		t.Fatal("derived PDF key rejected its own token")
	}
	// A token forged with the leaked/raw app secret must not be accepted.
	forged := PDFToken(appSecret, "5b1f0b7c9b1e4b0a8a2f3c4d", now)
	if ValidPDFToken(pdfSecret, "5b1f0b7c9b1e4b0a8a2f3c4d", forged, now) {
		t.Fatal("token signed with the raw app secret was accepted")
	}
	if ValidPDFToken(DerivePDFSecret("another-secret"), "5b1f0b7c9b1e4b0a8a2f3c4d", token, now) {
		t.Fatal("token accepted with another key")
	}
}
