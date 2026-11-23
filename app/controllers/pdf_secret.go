package controllers

import (
	"strings"

	"github.com/gemsnote/gemsnote/app/lea"
	"github.com/revel/revel"
)

// pdfSecret returns the key used to sign single-note PDF capabilities.
// pdf.secret allows deployments to keep it fully separate from app.secret;
// otherwise a dedicated subkey is derived so the session signing key is never
// itself handed to the renderer or used as a bearer secret.
func pdfSecret() string {
	if configured := strings.TrimSpace(revel.Config.StringDefault("pdf.secret", "")); configured != "" {
		return configured
	}
	appSecret, _ := revel.Config.String("app.secret")
	return lea.DerivePDFSecret(appSecret)
}

// PDFSecret exposes the PDF capability key to the legacy API package.
func PDFSecret() string { return pdfSecret() }
