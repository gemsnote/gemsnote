package controllers

import (
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
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
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
	return c.RenderBinary(file, title, delivery, stat.ModTime())
}
