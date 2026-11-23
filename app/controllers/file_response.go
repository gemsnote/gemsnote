package controllers

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/revel/revel"
)

// RenderStoredFile serves an already-authorized database file path. Callers
// must check access before calling it. Revel owns and closes a successful file
// response after writing it; closing the file here would break seeking/ranges.
func RenderStoredFile(c *revel.Controller, storedPath, title string, delivery revel.ContentDisposition) revel.Result {
	fullPath := filepath.Join(revel.BasePath, strings.TrimLeft(storedPath, "/"))
	file, err := openStoredFile(storedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			revel.AppLog.Warn("Download file missing", "path", storedPath, "error", err)
			c.Response.Status = http.StatusNotFound
			return c.RenderText("fileNotFound")
		}
		revel.AppLog.Error("Download file open failed", "path", storedPath, "error", err)
		c.Response.Status = http.StatusInternalServerError
		return c.RenderText("fileReadFailed")
	}
	stat, err := file.Stat()
	if err != nil || !stat.Mode().IsRegular() {
		file.Close()
		if err != nil {
			revel.AppLog.Error("Download file stat failed", "path", storedPath, "error", err)
			c.Response.Status = http.StatusInternalServerError
			return c.RenderText("fileReadFailed")
		}
		revel.AppLog.Warn("Download path is not a regular file", "path", storedPath)
		c.Response.Status = http.StatusNotFound
		return c.RenderText("fileNotFound")
	}
	if title == "" {
		title = filepath.Base(fullPath)
	}
	title = safeDownloadName(title)
	SetFileSecurityHeaders(c)
	return c.RenderBinary(file, title, delivery, stat.ModTime())
}

// SetFileSecurityHeaders applies the response hardening required for any file
// that leaves the server: no content sniffing, no script execution, no caching
// by shared proxies.
func SetFileSecurityHeaders(c *revel.Controller) {
	if c.Response.Out.Server == nil {
		return
	}
	header := c.Response.Out.Header()
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Content-Security-Policy", "sandbox; default-src 'none'; frame-ancestors 'none'")
	header.Set("Cache-Control", "private, no-store")
}

// OpenStoredFile validates a database-stored path before opening it.
func OpenStoredFile(stored string) (*os.File, error) { return openStoredFile(stored) }

// SafeDownloadName keeps user titles out of header and archive separators.
func SafeDownloadName(name string) string { return safeDownloadName(name) }

// Paths from old snapshots are untrusted too. Resolve links and reject any
// path outside the instance directory before opening it.
func openStoredFile(stored string) (*os.File, error) {
	if strings.Contains(stored, "\\") || strings.ContainsRune(stored, 0) {
		return nil, os.ErrNotExist
	}
	root, err := filepath.EvalSymlinks(revel.BasePath)
	if err != nil {
		return nil, err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	candidate := filepath.Join(root, strings.TrimLeft(stored, "/"))
	rel, err := filepath.Rel(root, candidate)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return nil, os.ErrNotExist
	}
	resolved, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return nil, err
	}
	rel, err = filepath.Rel(root, resolved)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return nil, os.ErrNotExist
	}
	file, err := os.Open(resolved)
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil || !stat.Mode().IsRegular() {
		file.Close()
		return nil, errors.Join(os.ErrNotExist, err)
	}
	return file, nil
}
