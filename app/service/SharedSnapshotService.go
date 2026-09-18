package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
)

// SharedSnapshotItem is the server-side, fully materialized representation
// consumed by desktop clients. Exactly one of Notebook, Note and File is set.
type SharedSnapshotItem struct {
	Kind     string                  `json:"Kind"`
	Notebook *SharedSnapshotNotebook `json:"Notebook,omitempty"`
	Note     *SharedSnapshotNote     `json:"Note,omitempty"`
	File     *SharedSnapshotFile     `json:"File,omitempty"`
}

type SharedSnapshotNotebook struct {
	NotebookId       string `json:"NotebookId"`
	OwnerUserId      string `json:"OwnerUserId"`
	ParentNotebookId string `json:"ParentNotebookId,omitempty"`
	Title            string `json:"Title"`
	Seq              int    `json:"Seq"`
	Perm             int    `json:"Perm"`
	Version          string `json:"Version"`
}

type SharedSnapshotNote struct {
	NoteId      string    `json:"NoteId"`
	NotebookId  string    `json:"NotebookId,omitempty"`
	OwnerUserId string    `json:"OwnerUserId"`
	Title       string    `json:"Title"`
	Desc        string    `json:"Desc,omitempty"`
	Tags        []string  `json:"Tags,omitempty"`
	IsMarkdown  bool      `json:"IsMarkdown"`
	Perm        int       `json:"Perm"`
	Version     string    `json:"Version"`
	CreatedTime time.Time `json:"CreatedTime"`
	UpdatedTime time.Time `json:"UpdatedTime"`
}

type SharedSnapshotFile struct {
	NoteId  string `json:"NoteId"`
	FileId  string `json:"FileId"`
	Kind    string `json:"Kind"`
	Title   string `json:"Title,omitempty"`
	Size    int64  `json:"Size"`
	Version string `json:"Version"`
}

// BuildSharedSnapshot retries until two consecutive materializations match.
// This prevents publishing a deletion-capable snapshot assembled across a
// concurrent permission or note change on databases without snapshot reads.
func (this *ShareService) BuildSharedSnapshot(userId, basePath string) ([]SharedSnapshotItem, error) {
	var previous []SharedSnapshotItem
	var previousDigest string
	for attempt := 0; attempt < 4; attempt++ {
		items := this.buildSharedSnapshotOnce(userId, basePath)
		digest := sharedSnapshotDigest(items)
		if attempt > 0 && digest == previousDigest {
			return items, nil
		}
		previous, previousDigest = items, digest
	}
	_ = previous
	return nil, errors.New("shared snapshot changed while it was being created")
}

func sharedSnapshotDigest(items []SharedSnapshotItem) string {
	data, _ := json.Marshal(items)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func (this *ShareService) buildSharedSnapshotOnce(userId, basePath string) []SharedSnapshotItem {
	q := this.getOrQ(userId)
	sharedNotes := []info.ShareNote{}
	sharedNotebooks := []info.ShareNotebook{}
	db.ListByQ(db.ShareNotes, q, &sharedNotes)
	db.ListByQ(db.ShareNotebooks, q, &sharedNotebooks)

	noteCandidates := map[bson.ObjectId]bool{}
	notebookCandidates := map[bson.ObjectId]bool{}
	for _, share := range sharedNotes {
		noteCandidates[share.NoteId] = true
	}
	for _, share := range sharedNotebooks {
		notebookCandidates[share.NotebookId] = true
	}
	if len(notebookCandidates) > 0 {
		ids := make([]bson.ObjectId, 0, len(notebookCandidates))
		for id := range notebookCandidates {
			ids = append(ids, id)
		}
		notes := []info.Note{}
		db.ListByQ(db.Notes, bson.M{"NotebookId": bson.M{"$in": ids}, "IsTrash": false, "IsDeleted": false}, &notes)
		for _, note := range notes {
			noteCandidates[note.NoteId] = true
		}
	}

	noteIds := make([]bson.ObjectId, 0, len(noteCandidates))
	for id := range noteCandidates {
		noteIds = append(noteIds, id)
	}
	notes := []info.Note{}
	if len(noteIds) > 0 {
		db.ListByQ(db.Notes, bson.M{"_id": bson.M{"$in": noteIds}, "IsTrash": false, "IsDeleted": false}, &notes)
	}
	sort.Slice(notes, func(i, j int) bool { return notes[i].NoteId.Hex() < notes[j].NoteId.Hex() })

	// Filter again through the canonical permission implementation. Candidate
	// relationships are only an efficient upper bound, never authorization.
	effectiveNotes := make([]info.Note, 0, len(notes))
	usedNotebooks := map[bson.ObjectId]bool{}
	for id := range notebookCandidates {
		usedNotebooks[id] = true
	}
	for _, note := range notes {
		ownerId, noteId := note.UserId.Hex(), note.NoteId.Hex()
		if ownerId == userId || !this.HasReadPerm(ownerId, userId, noteId) {
			continue
		}
		effectiveNotes = append(effectiveNotes, note)
	}

	items := []SharedSnapshotItem{}
	if len(usedNotebooks) > 0 {
		ids := make([]bson.ObjectId, 0, len(usedNotebooks))
		for id := range usedNotebooks {
			ids = append(ids, id)
		}
		notebooks := []info.Notebook{}
		db.ListByQ(db.Notebooks, bson.M{"_id": bson.M{"$in": ids}, "IsTrash": false, "IsDeleted": false}, &notebooks)
		sort.Slice(notebooks, func(i, j int) bool { return notebooks[i].NotebookId.Hex() < notebooks[j].NotebookId.Hex() })
		for _, notebook := range notebooks {
			if notebook.UserId.Hex() == userId || !this.HasReadNotebookPerm(notebook.UserId.Hex(), userId, notebook.NotebookId.Hex()) {
				continue
			}
			perm := 0
			if this.HasUpdateNotebookPerm(notebook.UserId.Hex(), userId, notebook.NotebookId.Hex()) {
				perm = 1
			}
			entry := &SharedSnapshotNotebook{NotebookId: notebook.NotebookId.Hex(), OwnerUserId: notebook.UserId.Hex(), ParentNotebookId: notebook.ParentNotebookId.Hex(), Title: notebook.Title, Seq: notebook.Seq, Perm: perm, Version: sharedValueVersion(notebook)}
			items = append(items, SharedSnapshotItem{Kind: "notebook", Notebook: entry})
		}
	}

	for _, note := range effectiveNotes {
		perm := 0
		if this.HasUpdatePerm(note.UserId.Hex(), userId, note.NoteId.Hex()) {
			perm = 1
		}
		content := noteService.GetNoteContent(note.NoteId.Hex(), note.UserId.Hex())
		contentSum := sha256.Sum256([]byte(content.Content))
		entry := &SharedSnapshotNote{NoteId: note.NoteId.Hex(), NotebookId: note.NotebookId.Hex(), OwnerUserId: note.UserId.Hex(), Title: note.Title, Desc: note.Desc, Tags: note.Tags, IsMarkdown: note.IsMarkdown, Perm: perm, Version: hex.EncodeToString(contentSum[:]), CreatedTime: note.CreatedTime, UpdatedTime: note.UpdatedTime}
		items = append(items, SharedSnapshotItem{Kind: "note", Note: entry})
	}

	if len(effectiveNotes) > 0 {
		emittedFiles := map[string]bool{}
		effectiveIds := make([]bson.ObjectId, len(effectiveNotes))
		for i := range effectiveNotes {
			effectiveIds[i] = effectiveNotes[i].NoteId
		}
		imagesByNote := map[string][]bson.ObjectId{}
		links := []info.NoteImage{}
		db.ListByQ(db.NoteImages, bson.M{"NoteId": bson.M{"$in": effectiveIds}}, &links)
		imageIds := []bson.ObjectId{}
		for _, link := range links {
			imagesByNote[link.NoteId.Hex()] = append(imagesByNote[link.NoteId.Hex()], link.ImageId)
			imageIds = append(imageIds, link.ImageId)
		}
		imageFiles := []info.File{}
		if len(imageIds) > 0 {
			db.ListByQ(db.Files, bson.M{"_id": bson.M{"$in": imageIds}}, &imageFiles)
		}
		imageMap := map[bson.ObjectId]info.File{}
		for _, file := range imageFiles {
			imageMap[file.FileId] = file
		}
		for noteId, ids := range imagesByNote {
			for _, id := range ids {
				if file, ok := imageMap[id]; ok {
					key := noteId + ":image:" + id.Hex()
					if emittedFiles[key] {
						continue
					}
					emittedFiles[key] = true
					f := &SharedSnapshotFile{NoteId: noteId, FileId: id.Hex(), Kind: "image", Title: file.Title, Size: file.Size, Version: sharedFileDigest(basePath, file.Path)}
					items = append(items, SharedSnapshotItem{Kind: "file", File: f})
				}
			}
		}
		attachments := []info.Attach{}
		db.ListByQ(db.Attachs, bson.M{"NoteId": bson.M{"$in": effectiveIds}}, &attachments)
		for _, attach := range attachments {
			key := attach.NoteId.Hex() + ":attachment:" + attach.AttachId.Hex()
			if emittedFiles[key] {
				continue
			}
			emittedFiles[key] = true
			f := &SharedSnapshotFile{NoteId: attach.NoteId.Hex(), FileId: attach.AttachId.Hex(), Kind: "attachment", Title: attach.Title, Size: attach.Size, Version: sharedFileDigest(basePath, attach.Path)}
			items = append(items, SharedSnapshotItem{Kind: "file", File: f})
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		left, right := items[i].Kind+sharedItemId(items[i]), items[j].Kind+sharedItemId(items[j])
		return left < right
	})
	return items
}

func sharedFileDigest(basePath, storedPath string) string {
	relative := filepath.Clean(strings.TrimLeft(storedPath, "/\\"))
	if relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return ""
	}
	file, err := os.Open(filepath.Join(basePath, relative))
	if err != nil {
		return ""
	}
	defer file.Close()
	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func sharedItemId(item SharedSnapshotItem) string {
	if item.Notebook != nil {
		return item.Notebook.NotebookId
	}
	if item.Note != nil {
		return item.Note.NoteId
	}
	if item.File != nil {
		return item.File.NoteId + item.File.Kind + item.File.FileId
	}
	return ""
}

func sharedValueVersion(value interface{}) string {
	data, _ := json.Marshal(value)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
