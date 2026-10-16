package service

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/info"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// Shared by API2 and browser mutations that replace a whole note snapshot.
var HistoryMutationMutex sync.Mutex

func (s *NoteContentHistoryService) SnapshotCurrent(noteID, userID string) bool {
	n := noteService.GetNoteByIdAndUserId(noteID, userID)
	if n.NoteId == "" || n.IsDeleted {
		return false
	}
	c := noteService.GetNoteContent(noteID, userID)
	if c.NoteId == "" {
		return false
	}
	return s.AddHistory(noteID, userID, info.EachHistory{
		UpdatedUserId: n.UserId, UpdatedTime: time.Now(), Content: c.Content,
		IsMarkdown: n.IsMarkdown, FilesComplete: true,
		Files:             noteService.getFiles([]bson.ObjectId{n.NoteId})[noteID],
		AttachmentRecords: attachService.ListAttachs(noteID, userID),
	})
}

func historyFilePath(path string) (string, bool) {
	path = strings.TrimLeft(path, "/")
	clean := filepath.Clean(path)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", false
	}
	if strings.HasPrefix(clean, "upload/") {
		clean = "public/" + clean
	}
	full := filepath.Join(revel.BasePath, clean)
	stat, err := os.Stat(full)
	return full, err == nil && stat.Mode().IsRegular()
}

// Validate the entire set before mutating anything. Legacy missing manifests
// are not equivalent to an empty attachment set.
func validateHistoryAttachments(h info.EachHistory, noteID string) ([]info.Attach, string) {
	if !h.FilesComplete {
		return nil, "historyFilesUnknown"
	}
	records := map[string]info.Attach{}
	for _, a := range h.AttachmentRecords {
		if a.AttachId.Valid() && a.NoteId.Hex() == noteID {
			records[a.AttachId.Hex()] = a
		}
	}
	result := []info.Attach{}
	seen := map[string]bool{}
	for _, f := range h.Files {
		if !bson.IsObjectIdHex(f.FileId) || seen[f.FileId] {
			return nil, "historyFilesInvalid"
		}
		seen[f.FileId] = true
		if !f.IsAttach {
			continue
		}
		a, ok := records[f.FileId]
		if !ok {
			return nil, "historyFileMissing"
		}
		if _, ok = historyFilePath(a.Path); !ok {
			return nil, "historyFileMissing"
		}
		a.Title, a.Type = f.Title, f.Type
		result = append(result, a)
	}
	return result, ""
}

func (s *NoteContentHistoryService) Restore(noteID, userID, historyID string, usn int) (bool, string) {
	HistoryMutationMutex.Lock()
	defer HistoryMutationMutex.Unlock()
	if !bson.IsObjectIdHex(noteID) || !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(historyID) {
		return false, "invalidId"
	}
	n := noteService.GetNoteByIdAndUserId(noteID, userID)
	if n.NoteId == "" || n.IsDeleted || n.IsTrash {
		return false, "noAuth"
	}
	if usn <= 0 || n.Usn != usn {
		return false, "conflict"
	}
	h, found := s.GetHistoryByID(noteID, userID, historyID)
	if !found {
		return false, "historyNotFound"
	}
	attachments, msg := validateHistoryAttachments(h, noteID)
	if msg != "" {
		return false, msg
	}
	for _, f := range h.Files {
		if f.IsAttach {
			continue
		}
		path := fileService.GetFile(userID, f.FileId)
		if _, ok := historyFilePath(path); path == "" || !ok {
			return false, "historyFileMissing"
		}
	}
	previousAttachments := attachService.ListAttachs(noteID, userID)
	previousContent := noteService.GetNoteContent(noteID, userID)
	if !s.SnapshotCurrent(noteID, userID) {
		return false, "historySaveFailed"
	}
	// CAS prevents a stale client from replacing a newer note. Files are retained
	// on disk, including on failure, so the saved pre-restore snapshot is usable.
	next := userService.IncrUsn(userID)
	if !db.UpdateByQMap(db.Notes, bson.M{"_id": n.NoteId, "UserId": n.UserId, "Usn": usn, "IsDeleted": false},
		bson.M{"Usn": next, "IsMarkdown": h.IsMarkdown, "UpdatedTime": time.Now(), "UpdatedUserId": n.UserId}) {
		return false, "conflict"
	}
	committed := false
	defer func() {
		if committed {
			return
		}
		// Best-effort rollback on storage failure. Never remove the saved history
		// or original bytes; even a failed rollback leaves a recoverable snapshot.
		current := noteService.GetNoteByIdAndUserId(noteID, userID)
		if current.Usn != next {
			return
		}
		for _, a := range previousAttachments {
			db.Upsert(db.Attachs, bson.M{"_id": a.AttachId, "NoteId": n.NoteId}, a)
		}
		previous := map[string]bool{}
		for _, a := range previousAttachments {
			previous[a.AttachId.Hex()] = true
		}
		for _, a := range attachService.ListAttachs(noteID, userID) {
			if !previous[a.AttachId.Hex()] {
				db.Delete(db.Attachs, bson.M{"_id": a.AttachId, "NoteId": n.NoteId})
			}
		}
		db.UpdateByIdAndUserIdMap(db.NoteContents, noteID, userID, bson.M{"Content": previousContent.Content, "Abstract": previousContent.Abstract})
		db.UpdateByQMap(db.Notes, bson.M{"_id": n.NoteId, "UserId": n.UserId, "Usn": next}, bson.M{"IsMarkdown": n.IsMarkdown, "UpdatedTime": time.Now()})
		noteImageService.UpdateNoteImages(userID, noteID, n.ImgSrc, previousContent.Content)
		attachService.updateNoteAttachNum(n.NoteId, 0)
	}()
	for _, a := range attachments {
		if !db.Upsert(db.Attachs, bson.M{"_id": a.AttachId, "NoteId": n.NoteId}, a) {
			return false, "historyRestoreFailed"
		}
	}
	keep := map[string]bool{}
	for _, a := range attachments {
		keep[a.AttachId.Hex()] = true
	}
	for _, a := range attachService.ListAttachs(noteID, userID) {
		if !keep[a.AttachId.Hex()] && !db.Delete(db.Attachs, bson.M{"_id": a.AttachId, "NoteId": n.NoteId}) {
			return false, "historyRestoreFailed"
		}
	}
	if !attachService.updateNoteAttachNum(n.NoteId, 0) {
		return false, "historyRestoreFailed"
	}
	ok, msg, _ := noteService.UpdateNoteContent(userID, noteID, h.Content, "", true, -1, time.Now())
	committed = ok
	if !ok && msg == "" {
		msg = "historyRestoreFailed"
	}
	return ok, msg
}
