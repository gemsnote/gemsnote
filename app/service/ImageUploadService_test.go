package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImageUploadIDOwnerIsolation(t *testing.T) {
	one, err := ImageUploadKey("owner-a", "persistent_upload_key")
	if err != nil {
		t.Fatal(err)
	}
	retry, _ := ImageUploadKey("owner-a", "persistent_upload_key")
	other, _ := ImageUploadKey("owner-b", "persistent_upload_key")
	if one != retry || one == other {
		t.Fatal("upload identity is not stable/isolated")
	}
	if _, err := ImageUploadKey("owner-a", "../bad-key"); err == nil {
		t.Fatal("invalid key accepted")
	}
}

func TestSafeUploadPathRejectsEscapesAndSymlinks(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "files", "images"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"../files/image.png", "/files/image.png", "files/../config", "public/upload/avatar.png"} {
		if _, err := safeUploadPath(root, name); err == nil {
			t.Fatalf("accepted unsafe path: %s", name)
		}
	}
	if _, err := safeUploadPath(root, "files/images/removed.png"); err != nil {
		t.Fatal("retry must allow already removed original")
	}
	if _, err := safeUploadPath(root, "files/images"); err == nil {
		t.Fatal("directory accepted as an uploaded file")
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(root, "files", "link")); err != nil {
		t.Skip("symlink unavailable")
	}
	if _, err := safeUploadPath(root, "files/link/image.png"); err == nil {
		t.Fatal("symlink escape accepted")
	}
}
