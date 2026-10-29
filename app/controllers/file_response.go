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
	if c.Response.Out.Server != nil {
		c.Response.Out.Header().Set("X-Content-Type-Options", "nosniff")
		c.Response.Out.Header().Set("Content-Security-Policy", "sandbox; default-src 'none'")
		c.Response.Out.Header().Set("Cache-Control", "private, no-store")
	}
	return c.RenderBinary(file, title, delivery, stat.ModTime())
}

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
