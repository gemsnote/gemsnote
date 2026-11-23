package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/revel/revel"
)

func TestThemeTemplateAllowed(t *testing.T) {
	allowed := []string{"index.html", "note.html", "style.css", "theme.json"}
	for _, name := range allowed {
		if !themeTemplateAllowed(name) {
			t.Fatalf("allowed template rejected: %s", name)
		}
	}
	rejected := []string{"", "x.js", "x.htm", "x.svg", "../index.html", "a/b.html", "index.html\x00.js", "theme.json.bak"}
	for _, name := range rejected {
		if themeTemplateAllowed(name) {
			t.Fatalf("unsafe filename accepted: %q", name)
		}
	}
}

func TestThemeDirectoryWritable(t *testing.T) {
	base := t.TempDir()
	previous := revel.BasePath
	revel.BasePath = base
	defer func() { revel.BasePath = previous }()

	upload := filepath.Join(base, "public", "upload", "user", "themes", "t")
	if err := os.MkdirAll(upload, 0o755); err != nil {
		t.Fatal(err)
	}
	if !themeDirectoryWritable(upload) {
		t.Fatal("theme inside public/upload rejected")
	}
	if themeDirectoryWritable(filepath.Join(base, "public", "blog", "themes", "default")) {
		t.Fatal("theme outside public/upload accepted")
	}
	if themeDirectoryWritable(filepath.Join(base, "public", "upload", "..", "..")) {
		t.Fatal("escaping path accepted")
	}
	if themeDirectoryWritable("") {
		t.Fatal("empty path accepted")
	}
}

func TestThemeDirectoryWritableRejectsSymlinkEscape(t *testing.T) {
	base := t.TempDir()
	previous := revel.BasePath
	revel.BasePath = base
	defer func() { revel.BasePath = previous }()

	outside := filepath.Join(base, "outside")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	upload := filepath.Join(base, "public", "upload")
	if err := os.MkdirAll(upload, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(upload, "escape")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if themeDirectoryWritable(link) {
		t.Fatal("symlinked theme escaped public/upload")
	}
}
