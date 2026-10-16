package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gemsnote/gemsnote/app/info"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

func TestHistoryManifestValidation(t *testing.T) {
	before := revel.BasePath
	revel.BasePath = t.TempDir()
	defer func() { revel.BasePath = before }()
	if err := os.WriteFile(filepath.Join(revel.BasePath, "original"), []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}
	noteID := bson.NewObjectId()
	fileID := bson.NewObjectId()
	h := info.EachHistory{FilesComplete: true, Files: []info.NoteFile{{FileId: fileID.Hex(), IsAttach: true}}, AttachmentRecords: []info.Attach{{AttachId: fileID, NoteId: noteID, Path: "original"}}}
	if files, msg := validateHistoryAttachments(h, noteID.Hex()); msg != "" || len(files) != 1 {
		t.Fatalf("%v %s", files, msg)
	}
	for _, tc := range []struct {
		name   string
		change func(*info.EachHistory)
		want   string
	}{
		{"legacy", func(h *info.EachHistory) { h.FilesComplete = false }, "historyFilesUnknown"},
		{"missing metadata", func(h *info.EachHistory) { h.AttachmentRecords = nil }, "historyFileMissing"},
		{"other note", func(h *info.EachHistory) {
			h.AttachmentRecords = []info.Attach{{AttachId: fileID, NoteId: bson.NewObjectId(), Path: "original"}}
		}, "historyFileMissing"},
		{"missing bytes", func(h *info.EachHistory) {
			h.AttachmentRecords = []info.Attach{{AttachId: fileID, NoteId: noteID, Path: "missing"}}
		}, "historyFileMissing"},
		{"traversal", func(h *info.EachHistory) {
			h.AttachmentRecords = []info.Attach{{AttachId: fileID, NoteId: noteID, Path: "../original"}}
		}, "historyFileMissing"},
		{"duplicate", func(h *info.EachHistory) { h.Files = append(h.Files, h.Files[0]) }, "historyFilesInvalid"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			copy := h
			tc.change(&copy)
			if _, msg := validateHistoryAttachments(copy, noteID.Hex()); msg != tc.want {
				t.Fatalf("got %q", msg)
			}
		})
	}
	if _, msg := validateHistoryAttachments(info.EachHistory{FilesComplete: true}, noteID.Hex()); msg != "" {
		t.Fatal(msg)
	}
}
