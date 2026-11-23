package app

import (
	"net/http"
	"testing"
)

func TestApplySecurityHeaders(t *testing.T) {
	header := http.Header{}
	applySecurityHeaders(header, "/api2/bootstrap", false)
	for name, want := range map[string]string{
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "SAMEORIGIN",
		"Referrer-Policy":         "no-referrer",
		"Content-Security-Policy": "frame-ancestors 'self'",
	} {
		if got := header.Get(name); got != want {
			t.Fatalf("%s = %q, want %q", name, got, want)
		}
	}
	if header.Get("Strict-Transport-Security") != "" {
		t.Fatal("HSTS must not be advertised over plain HTTP")
	}
}

func TestApplySecurityHeadersKeepsActionCSP(t *testing.T) {
	header := http.Header{}
	header.Set("Content-Security-Policy", "sandbox; default-src 'none'")
	applySecurityHeaders(header, "/api2/attach/download", true)
	if got := header.Get("Content-Security-Policy"); got != "sandbox; default-src 'none'" {
		t.Fatalf("action CSP overwritten: %q", got)
	}
	if header.Get("Strict-Transport-Security") == "" {
		t.Fatal("HSTS missing on TLS responses")
	}
}
