package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gemsnote/gemsnote/app/db"
)

var imageUploadKeyPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{16,128}$`)

func ImageUploadKey(userID, uploadID string) (string, error) {
	if !imageUploadKeyPattern.MatchString(uploadID) {
		return "", errors.New("invalidUploadId")
	}
	return userID + ":" + uploadID, nil
}

// Unlike the idempotency key, this ID must not be knowable before the file row
// is committed. Otherwise a pre-existing reference could bypass retention.
func NewImageUploadID() (string, error) {
	var bytes [12]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes[:]), nil
}

// Only files carrying the new never-referenced receipt may be cleaned. A
// tombstone stays in the database to prevent later adoption of the same ID.
func CleanupNoteImage(root, userID, fileID string) error {
	file, err := db.ReserveImageCleanup(userID, fileID)
	if err != nil {
		return err
	}
	if file.UploadState == "deleted" {
		return nil
	}
	name, err := safeUploadPath(root, file.Path)
	if err != nil {
		return err
	}
	if err := os.Remove(name); err != nil && !os.IsNotExist(err) {
		return errors.New("cleanupFileFailed")
	}
	if !db.CompleteImageCleanup(file) {
		return errors.New("databaseUnavailable")
	}
	return nil
}

func safeUploadPath(root, name string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(name))
	if filepath.IsAbs(clean) || !strings.HasPrefix(clean, "files"+string(filepath.Separator)) || clean != filepath.FromSlash(name) {
		return "", errors.New("unsafeUploadPath")
	}
	base, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", errors.New("unsafeUploadPath")
	}
	current := base
	for _, part := range strings.Split(clean, string(filepath.Separator)) {
		current = filepath.Join(current, part)
		stat, err := os.Lstat(current)
		if os.IsNotExist(err) {
			continue
		} // interrupted cleanup can be retried
		if err != nil || stat.Mode()&os.ModeSymlink != 0 {
			return "", errors.New("unsafeUploadPath")
		}
		if current == filepath.Join(base, clean) && !stat.Mode().IsRegular() {
			return "", errors.New("unsafeUploadPath")
		}
	}
	return current, nil
}
