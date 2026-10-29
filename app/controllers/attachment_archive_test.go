package controllers

import (
	"archive/tar"
	"compress/gzip"
	"github.com/gemsnote/gemsnote/app/info"
	"github.com/revel/revel"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAttachmentArchiveNamesAndLifetime(t *testing.T) {
	previous := revel.BasePath
	revel.BasePath = t.TempDir()
	t.Cleanup(func() { revel.BasePath = previous })
	if err := os.WriteFile(filepath.Join(revel.BasePath, "original"), []byte("exact bytes"), 0600); err != nil {
		t.Fatal(err)
	}
	file, err := buildAttachmentArchive([]info.Attach{{Title: "../../secret", Path: "original"}, {Title: `C:\outside\secret`, Path: "original"}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { file.Close() })
	name := file.Name()
	gz, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(gz)
	for _, want := range []string{"secret", "2-secret"} {
		h, err := tr.Next()
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(tr)
		if err != nil || h.Name != want || string(body) != "exact bytes" {
			t.Fatalf("archive entry: %+v %q %v", h, body, err)
		}
	}
	if _, err := tr.Next(); err != io.EOF {
		t.Fatal(err)
	}
	gz.Close()
	file.Close()
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		t.Fatal("temporary archive retained")
	}
	if safeDownloadName("../../private\r\nX: value") != "private__X_ value" {
		t.Fatal("unsafe download name")
	}
}

func TestStoredPathsCannotEscapeInstance(t *testing.T) {
	previous := revel.BasePath
	root := t.TempDir()
	revel.BasePath = filepath.Join(root, "instance")
	os.Mkdir(revel.BasePath, 0700)
	t.Cleanup(func() { revel.BasePath = previous })
	os.WriteFile(filepath.Join(root, "secret"), []byte("private"), 0600)
	os.Symlink(filepath.Join(root, "secret"), filepath.Join(revel.BasePath, "link"))
	for _, name := range []string{"../secret", "/../secret", `..\secret`, "link"} {
		if file, err := openStoredFile(name); err == nil {
			file.Close()
			t.Fatalf("opened escaped path %q", name)
		}
		if file, err := buildAttachmentArchive([]info.Attach{{Title: "safe", Path: name}}); err == nil {
			file.Close()
			t.Fatalf("archived escaped path %q", name)
		}
	}
	for _, name := range []string{"../report", `C:\temp\report`, "a\r\nb"} {
		if strings.ContainsAny(safeDownloadName(name), "/\\\r\n:") {
			t.Fatalf("unsafe name %q", name)
		}
	}
}
