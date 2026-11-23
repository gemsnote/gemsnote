package lea

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

// The renderer receives a short-lived capability for one note, never the
// application's session-signing secret in its URL, logs or process arguments.
func PDFToken(secret, noteID string, now time.Time) string {
	if secret == "" || noteID == "" {
		return ""
	}
	expires := strconv.FormatInt(now.Add(5*time.Minute).Unix(), 10)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte("pdf:" + noteID + ":" + expires))
	return expires + "." + hex.EncodeToString(mac.Sum(nil))
}

func ValidPDFToken(secret, noteID, token string, now time.Time) bool {
	parts := strings.Split(token, ".")
	if secret == "" || noteID == "" || len(parts) != 2 {
		return false
	}
	expires, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || expires <= now.Unix() || expires > now.Add(5*time.Minute).Unix() {
		return false
	}
	actual, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte("pdf:" + noteID + ":" + parts[0]))
	return hmac.Equal(actual, mac.Sum(nil))
}

// DerivePDFSecret keeps PDF capability signing independent from the session
// cookie signing key, even when only app.secret is configured.
func DerivePDFSecret(appSecret string) string {
	if appSecret == "" {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write([]byte("gemsnote/pdf-credential/v1"))
	return hex.EncodeToString(mac.Sum(nil))
}

// PDFRenderArgs keeps the external renderer away from local files and plugins.
// Markdown notes need JavaScript (the export template renders them client-side
// and signals --window-status), so JS is only disabled for HTML notes; the
// renderer still has to be isolated at the OS level.
func PDFRenderArgs(markdown bool) []string {
	args := []string{"--lowquality", "--disable-local-file-access", "--no-plugins"}
	if markdown {
		args = append(args, "--window-status", "done")
	} else {
		args = append(args, "--disable-javascript")
	}
	return args
}
