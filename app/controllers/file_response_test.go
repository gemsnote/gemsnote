package controllers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/revel/revel"
)

func TestStoredFileDownloadRejectsMissingFilesAndDirectories(t *testing.T) {
	previous := revel.BasePath
	revel.BasePath = t.TempDir()
	t.Cleanup(func() { revel.BasePath = previous })
	if err := os.Mkdir(filepath.Join(revel.BasePath, "directory"), 0700); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"missing.bin", "directory"} {
		t.Run(path, func(t *testing.T) {
			c := &revel.Controller{Response: &revel.Response{}}
			result := RenderStoredFile(c, path, "original.bin", revel.Attachment)
			if _, binary := result.(*revel.BinaryResult); binary || c.Response.Status != http.StatusNotFound {
				t.Fatalf("invalid file was passed to the seek/stream handler: result=%T status=%d", result, c.Response.Status)
			}
			if _, ok := result.(*revel.RenderTextResult); !ok {
				t.Fatalf("unexpected missing-file response: %+v", result)
			}
		})
	}
}

func TestStoredFileDownloadRemainsSeekableUntilResponse(t *testing.T) {
	previous := revel.BasePath
	revel.BasePath = t.TempDir()
	t.Cleanup(func() { revel.BasePath = previous })
	for _, body := range []string{"", "0123456789"} {
		path := "unknown-extension.download-test"
		if err := os.WriteFile(filepath.Join(revel.BasePath, path), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		c := &revel.Controller{Response: &revel.Response{}}
		result := RenderStoredFile(c, "/"+path, "原件.bin", revel.Attachment)
		binary, ok := result.(*revel.BinaryResult)
		if !ok || binary.Name != "原件.bin" || binary.Delivery != revel.Attachment {
			t.Fatalf("unexpected download result: %+v", result)
		}
		file := binary.Reader.(*os.File)
		req := httptest.NewRequest(http.MethodGet, "/api2/file/getAttach", nil)
		if body != "" {
			req.Header.Set("Range", "bytes=2-5")
		}
		rec := httptest.NewRecorder()
		// Exercise the same seeking used by Revel's native stream writer.
		http.ServeContent(rec, req, binary.Name, binary.ModTime, file)
		wantStatus, wantBody := http.StatusOK, ""
		if body != "" {
			wantStatus, wantBody = http.StatusPartialContent, "2345"
		}
		file.Close()
		if rec.Code != wantStatus || rec.Body.String() != wantBody {
			t.Fatalf("file response = %d %q; want %d %q", rec.Code, rec.Body, wantStatus, wantBody)
		}
		if _, err := file.Seek(0, io.SeekStart); err == nil {
			t.Fatal("test did not close the response reader")
		}
	}
}
