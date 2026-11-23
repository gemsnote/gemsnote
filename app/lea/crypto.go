// contains two cryptographic functions for both storing and comparing passwords.
package lea

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptCost is the work factor for new password hashes. Existing hashes keep
// their recorded cost and keep verifying.
const BcryptCost = 12

// GenerateHash generates bcrypt hash from plaintext password
func GenerateHash(password string) ([]byte, error) {
	hex := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(hex, BcryptCost)
	if err != nil {
		return hashedPassword, err
	}
	return hashedPassword, nil
}

// CompareHash compares bcrypt password with a plaintext one. Returns true if passwords match
// and false if they do not.
func CompareHash(digest []byte, password string) bool {
	hex := []byte(password)
	if err := bcrypt.CompareHashAndPassword(digest, hex); err == nil {
		return true
	}
	return false
}
