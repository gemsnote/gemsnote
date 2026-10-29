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
