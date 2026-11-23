package lea

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

// DeriveImageSecret keeps image capabilities independent from the session
// signing key.
func DeriveImageSecret(appSecret string) string {
	if appSecret == "" {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write([]byte("gemsnote/image-credential/v1"))
	return hex.EncodeToString(mac.Sum(nil))
}

func imageCredentialMAC(secret, userId, fileId, expires string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte("image:" + userId + ":" + fileId + ":" + expires))
	return mac.Sum(nil)
}

// ImageToken returns a short-lived capability for one file on behalf of one
// user. It lets a plain <img> URL authorize a read without exposing a session
// cookie or long-lived API token.
func ImageToken(secret, userId, fileId string, now time.Time, ttl time.Duration) string {
	if secret == "" || userId == "" || fileId == "" || ttl <= 0 {
		return ""
	}
	expires := strconv.FormatInt(now.Add(ttl).Unix(), 10)
	return expires + "." + hex.EncodeToString(imageCredentialMAC(secret, userId, fileId, expires))
}

// ValidImageToken verifies the capability and rejects tokens that outlive the
// configured lifetime, so a captured URL stays useless quickly.
func ValidImageToken(secret, userId, fileId, token string, now time.Time, ttl time.Duration) bool {
	parts := strings.Split(token, ".")
	if secret == "" || userId == "" || fileId == "" || ttl <= 0 || len(parts) != 2 {
		return false
	}
	expires, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || expires <= now.Unix() || expires > now.Add(ttl).Unix() {
		return false
	}
	actual, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}
	return hmac.Equal(actual, imageCredentialMAC(secret, userId, fileId, parts[0]))
}
