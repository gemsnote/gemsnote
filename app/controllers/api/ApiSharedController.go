package api

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/info"
	. "github.com/gemsnote/gemsnote/app/lea"
	"github.com/gemsnote/gemsnote/app/service"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

const sharedSnapshotPageSize = 200
const sharedSnapshotLifetime = 10 * time.Minute
const maxSharedSnapshotItems = 10000
const maxSharedSnapshots = 256

type ApiShared struct{ ApiBaseContrller }

type materializedSharedSnapshot struct {
	UserId    string
	ExpiresAt time.Time
	Items     []service.SharedSnapshotItem
}

var sharedSnapshots = struct {
	sync.Mutex
	items map[string]materializedSharedSnapshot
}{items: map[string]materializedSharedSnapshot{}}

func (c ApiShared) Capabilities() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"Ok": true, "ProtocolVersion": 1, "Snapshot": true,
	})
}

func (c ApiShared) CreateSnapshot() revel.Result {
	items, err := shareService.BuildSharedSnapshot(c.getUserId(), revel.BasePath)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "snapshotUnstable"})
	}
	if len(items) > maxSharedSnapshotItems {
		return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "snapshotTooLarge"})
	}
	now := time.Now().UTC()
	expiresAt := now.Add(sharedSnapshotLifetime)
	snapshotId := NewGuid()
	sharedSnapshots.Lock()
	for id, snapshot := range sharedSnapshots.items {
		if !snapshot.ExpiresAt.After(now) || snapshot.UserId == c.getUserId() {
			delete(sharedSnapshots.items, id)
		}
	}
	if len(sharedSnapshots.items) >= maxSharedSnapshots {
		oldestId := ""
		var oldestExpiry time.Time
		for id, snapshot := range sharedSnapshots.items {
			if oldestId == "" || snapshot.ExpiresAt.Before(oldestExpiry) {
				oldestId, oldestExpiry = id, snapshot.ExpiresAt
			}
		}
		delete(sharedSnapshots.items, oldestId)
	}
	sharedSnapshots.items[snapshotId] = materializedSharedSnapshot{UserId: c.getUserId(), ExpiresAt: expiresAt, Items: items}
	sharedSnapshots.Unlock()
	return c.RenderJSON(map[string]interface{}{"Ok": true, "SnapshotId": snapshotId, "ExpiresAt": expiresAt, "Total": len(items)})
}

func (c ApiShared) SnapshotItems(snapshotId, pageToken string) revel.Result {
	sharedSnapshots.Lock()
	snapshot, ok := sharedSnapshots.items[snapshotId]
	sharedSnapshots.Unlock()
	if !ok || snapshot.UserId != c.getUserId() || !snapshot.ExpiresAt.After(time.Now().UTC()) {
		return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "snapshotExpired"})
	}
	offset := 0
	if pageToken != "" {
		parsed, err := strconv.Atoi(pageToken)
		if err != nil || parsed < 0 || parsed > len(snapshot.Items) {
			return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "pageTokenInvalid"})
		}
		offset = parsed
	}
	end := offset + sharedSnapshotPageSize
	if end > len(snapshot.Items) {
		end = len(snapshot.Items)
	}
	next := ""
	if end < len(snapshot.Items) {
		next = strconv.Itoa(end)
	}
	page := snapshot.Items[offset:end]
	if page == nil {
		page = []service.SharedSnapshotItem{}
	}
	return c.RenderJSON(map[string]interface{}{"Ok": true, "Items": page, "NextPageToken": next, "Complete": end == len(snapshot.Items), "Total": len(snapshot.Items)})
}

func (c ApiShared) NoteContent(noteId string) revel.Result {
	if !bson.IsObjectIdHex(noteId) {
		return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "noteIdInvalid"})
	}
	for attempt := 0; attempt < 3; attempt++ {
		note := noteService.GetNoteById(noteId)
		if note.NoteId == "" || note.IsTrash || note.IsDeleted || note.UserId.Hex() == c.getUserId() || !shareService.HasReadPerm(note.UserId.Hex(), c.getUserId(), noteId) {
			return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "noPermission"})
		}
		content := noteService.GetNoteContent(noteId, note.UserId.Hex())
		noteAfter := noteService.GetNoteById(noteId)
		if noteAfter.NoteId != "" && note.Usn == noteAfter.Usn && note.UpdatedTime.Equal(noteAfter.UpdatedTime) && content.NoteId != "" {
			sum := sha256.Sum256([]byte(content.Content))
			digest := hex.EncodeToString(sum[:])
			return c.RenderJSON(map[string]interface{}{"Ok": true, "NoteId": noteId, "Version": digest, "Digest": digest, "Content": content.Content})
		}
	}
	return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "contentUnstable"})
}

func (c ApiShared) NoteFile(noteId, fileId string) revel.Result {
	if !bson.IsObjectIdHex(noteId) || !bson.IsObjectIdHex(fileId) {
		return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "fileIdInvalid"})
	}
	note := noteService.GetNoteById(noteId)
	if note.NoteId == "" || note.IsTrash || note.IsDeleted || note.UserId.Hex() == c.getUserId() || !shareService.HasReadPerm(note.UserId.Hex(), c.getUserId(), noteId) {
		return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "noPermission"})
	}
	path, title, inline := "", "", false
	attach := info.Attach{}
	db.GetByQ(db.Attachs, bson.M{"_id": bson.ObjectIdHex(fileId), "NoteId": bson.ObjectIdHex(noteId)}, &attach)
	if attach.AttachId != "" {
		path, title = attach.Path, attach.Title
	} else {
		link := info.NoteImage{}
		db.GetByQ(db.NoteImages, bson.M{"NoteId": bson.ObjectIdHex(noteId), "ImageId": bson.ObjectIdHex(fileId)}, &link)
		if link.NoteImageId == "" {
			return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "fileNotFound"})
		}
		fileInfo := info.File{}
		db.Get(db.Files, fileId, &fileInfo)
		if fileInfo.FileId == "" {
			return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "fileNotFound"})
		}
		path, title, inline = fileInfo.Path, fileInfo.Title, true
	}
	fullPath := revel.BasePath + "/" + strings.TrimLeft(path, "/")
	file, err := os.Open(fullPath)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "fileNotFound"})
	}
	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		file.Close()
		return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "fileReadFailed"})
	}
	if _, err = file.Seek(0, 0); err != nil {
		file.Close()
		return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "fileReadFailed"})
	}
	digest := hex.EncodeToString(hash.Sum(nil))
	c.Response.Out.Header().Set("X-Gemsnote-SHA256", digest)
	c.Response.Out.Header().Set("ETag", `"`+digest+`"`)
	if inline {
		return c.RenderFile(file, revel.Inline)
	}
	return c.RenderBinary(file, title, revel.Attachment, time.Now())
}
