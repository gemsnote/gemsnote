package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSharedFileDigest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.bin")
	if err := os.WriteFile(path, []byte("gemsnote"), 0600); err != nil {
		t.Fatal(err)
	}
	const expected = "029d60ab82c57305820e5e60bc43f8913f14e3670f64f7a55ff06df134ba4163"
	if got := sharedFileDigest(dir, "file.bin"); got != expected {
		t.Fatalf("unexpected digest: %s", got)
	}
	if got := sharedFileDigest(dir, "../file.bin"); got != "" {
		t.Fatalf("path outside base must be rejected, got %s", got)
	}
}

func TestSharedSnapshotDigestIsDeterministic(t *testing.T) {
	items := []SharedSnapshotItem{{Kind: "note", Note: &SharedSnapshotNote{NoteId: "0123456789abcdef01234567", Version: "v1"}}}
	first := sharedSnapshotDigest(items)
	second := sharedSnapshotDigest(items)
	if first == "" || first != second {
		t.Fatalf("snapshot digest is not deterministic: %q != %q", first, second)
	}
	items[0].Note.Version = "v2"
	if sharedSnapshotDigest(items) == first {
		t.Fatal("snapshot digest did not change with content version")
	}
}
