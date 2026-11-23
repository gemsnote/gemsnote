package app

import (
	"strings"

	"github.com/revel/revel"
)

// headerWriter is satisfied by both http.Header and Revel's header wrapper.
type headerWriter interface {
	Get(string) string
	Set(string, string)
}

// applySecurityHeaders centralises defaults for every response. File responses
// refine the CSP (sandboxing untrusted documents) and upload paths get a
// stricter sandbox; those values win because they are written later.
func applySecurityHeaders(header headerWriter, path string, tls bool) {
	if header == nil {
		return
	}
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("X-Frame-Options", "SAMEORIGIN")
	header.Set("Referrer-Policy", "no-referrer")
	if tls {
		header.Set("Strict-Transport-Security", "max-age=31536000")
	}
	// Uploaded themes/avatars live under a static route. Sandbox them so an
	// uploaded document can never run script with the application's origin.
	if strings.HasPrefix(path, "/upload/") || strings.HasPrefix(path, "/public/upload/") {
		header.Set("Content-Security-Policy", "sandbox; default-src 'none'; frame-ancestors 'none'")
		return
	}
	if header.Get("Content-Security-Policy") == "" {
		header.Set("Content-Security-Policy", "frame-ancestors 'self'")
	}
}

// SecurityHeadersFilter adds the headers before the action runs so that actions
// which need a different CSP (downloads, uploads) can override them.
func SecurityHeadersFilter(c *revel.Controller, fc []revel.Filter) {
	if c.Response != nil && c.Response.Out.Server != nil {
		path := ""
		if c.Request != nil && c.Request.URL != nil {
			path = strings.TrimSpace(c.Request.URL.Path)
		}
		applySecurityHeaders(c.Response.Out.Header(), path, revel.HTTPSsl)
	}
	fc[0](c, fc[1:])
}
