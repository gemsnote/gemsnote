package app

import (
	"os"
	"strings"

	"github.com/gemsnote/gemsnote/app/lea"
	"github.com/revel/revel"
)

// databaseSecretFromEnvironment reports whether the deployment injects its
// database credential or URL instead of relying on the shipped configuration.
func databaseSecretFromEnvironment() bool {
	for _, name := range []string{"GEMSNOTE_DB_URL", "GEMSNOTE_DB_PASSWORD"} {
		if strings.TrimSpace(os.Getenv(name)) != "" {
			return true
		}
	}
	return false
}

// Production deployments must not sign sessions or PDF capabilities with the
// public default secret. Development only warns so a fresh checkout still runs.
func init() {
	revel.OnAppStart(func() {
		secret, _ := revel.Config.String("app.secret")
		production := strings.EqualFold(revel.RunMode, "prod")
		if err := lea.ValidateAppSecret(secret, production); err != nil {
			if production {
				revel.RevelLog.Error("insecure app.secret; refusing to start", "error", err)
				revel.RevelLog.Fatal("insecure app.secret; refusing to start", "error", err)
			}
			revel.RevelLog.Warn("insecure app.secret; set a unique value before deploying", "error", err, "mode", revel.RunMode)
		}
		if strings.TrimSpace(revel.Config.StringDefault("pdf.secret", "")) == "" {
			revel.RevelLog.Info("pdf.secret is unset; deriving the PDF capability key from app.secret")
		}
		// Revel hardcodes HttpOnly on its session cookie (v1.0.0), but Secure
		// still follows http.ssl unless explicitly configured.
		if production && !revel.CookieSecure {
			revel.RevelLog.Warn("session cookies are not marked Secure; terminate TLS and set http.ssl=true or cookie.secure=true")
		}
		if production && !databaseSecretFromEnvironment() {
			password, _ := revel.Config.String("db.password")
			if _, hasURL := revel.Config.String("db.url"); !hasURL && (strings.TrimSpace(password) == "" || password == "gemsnote") {
				revel.RevelLog.Warn("database uses an empty or default password; set GEMSNOTE_DB_PASSWORD or db.password")
			}
		}
	})
}
