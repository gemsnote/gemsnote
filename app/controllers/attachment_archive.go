package controllers

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
	"unicode"

	"github.com/gemsnote/gemsnote/app/info"
	"github.com/revel/revel"
)

func safeDownloadName(name string) string {
	name = path.Base(strings.ReplaceAll(name, "\\", "/"))
	name = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || strings.ContainsRune(":*?\"<>|", r) {
			return '_'
		}
		return r
	}, name)
	if name == "" || name == "." || name == ".." || name == "/" {
		return "file"
	}
	return name
}

type temporaryArchive struct{ *os.File }

func (f *temporaryArchive) Close() error {
	err := f.File.Close()
	os.Remove(f.Name())
	return err
}

// Callers must authorize the attachment list first. Neither note titles nor
// attachment names are used as disk paths. The response owns the temporary file.
func RenderAttachmentArchive(c *revel.Controller, title string, attachments []info.Attach) revel.Result {
	file, err := buildAttachmentArchive(attachments)
	if err != nil {
		c.Response.Status = 500
		if os.IsNotExist(err) {
			c.Response.Status = 404
		}
		return c.RenderText("attachmentArchiveFailed")
	}
	return c.RenderBinary(file, safeDownloadName(title)+".tar.gz", revel.Attachment, time.Now())
}

func buildAttachmentArchive(attachments []info.Attach) (_ *temporaryArchive, err error) {
	f, err := os.CreateTemp("", "gemsnote-attachments-*.tar.gz")
	if err != nil {
		return nil, err
	}
	result := &temporaryArchive{f}
	defer func() {
		if err != nil {
			result.Close()
		}
	}()
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	names := map[string]bool{}
	for _, attach := range attachments {
		err = func() error {
			src, err := openStoredFile(attach.Path)
			if err != nil {
				return err
			}
			defer src.Close()
			stat, err := src.Stat()
			if err != nil {
				return err
			}
			base := safeDownloadName(attach.Title)
			name := base
			for i := 2; names[strings.ToLower(name)]; i++ {
				name = fmt.Sprintf("%d-%s", i, base)
			}
			names[strings.ToLower(name)] = true
			if err := tw.WriteHeader(&tar.Header{Name: name, Size: stat.Size(), Mode: 0600, ModTime: stat.ModTime()}); err != nil {
				return err
			}
			_, err = io.CopyN(tw, src, stat.Size())
			return err
		}()
		if err != nil {
			tw.Close()
			gz.Close()
			return nil, err
		}
	}
	if err = tw.Close(); err != nil {
		gz.Close()
		return nil, err
	}
	if err = gz.Close(); err != nil {
		return nil, err
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	return result, nil
}
