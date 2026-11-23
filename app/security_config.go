package app

import (
	"strings"

	"github.com/gemsnote/gemsnote/app/lea"
	"github.com/revel/revel"
)

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
	})
}
