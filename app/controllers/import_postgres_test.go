package controllers

import (
	"encoding/json"
	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/info"
	"github.com/gemsnote/gemsnote/app/service"
	"github.com/revel/config"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPostgresImportedMetadataAndNotebookDelete(t *testing.T) {
	dsn := os.Getenv("GEMSNOTE_INTEGRATION_POSTGRES_URL")
	if dsn == "" {
		t.Skip("set GEMSNOTE_INTEGRATION_POSTGRES_URL to an isolated PostgreSQL database")
	}
	base, _ := filepath.Abs("../..")
	t.Setenv("GEMSNOTE_DATABASE_DIR", filepath.Join(base, "database"))
	revel.BasePath = base
	var err error
	revel.Config, err = config.LoadContext("app.conf", []string{filepath.Join(base, "conf")})
	if err != nil {
		t.Fatal(err)
	}
	revel.Config.SetOption("db.type", "postgresql")
	db.Init(dsn, "")
	defer db.Close()
	service.InitService()
	noteService, notebookService, tagService = service.NoteS, service.NotebookS, service.TagS
	user, book := bson.NewObjectId(), bson.NewObjectId()
	if !db.Insert(db.Users, info.User{UserId: user, Username: user.Hex(), Flag: info.FlagAvailable, CreatedTime: time.Now()}) || !db.Insert(db.Notebooks, info.Notebook{NotebookId: book, UserId: user, Title: "Import", CreatedTime: time.Now(), UpdatedTime: time.Now()}) {
		t.Fatal("fixture insert")
	}
	makeController := func(path string, payload any) BaseController {
		body, _ := json.Marshal(payload)
		requestURL, _ := url.Parse(path)
		return BaseController{&revel.Controller{Request: &revel.Request{URL: requestURL}, Response: &revel.Response{Status: 200}, Params: &revel.Params{JSON: body}, Session: map[string]interface{}{"UserId": user.Hex()}, ViewArgs: map[string]interface{}{}}}
	}
	created, _ := time.Parse(time.RFC3339Nano, "2020-01-02T03:04:05.123+08:00")
	updated, _ := time.Parse(time.RFC3339Nano, "2021-02-03T04:05:06.456Z")
	noteID := bson.NewObjectId().Hex()
	(Api2{makeController("/api2/save", map[string]any{"noteId": noteID, "notebookId": book.Hex(), "title": "Hugo Title", "content": "# Body\n", "tags": "Go,中文", "isNew": true, "isMarkdown": true, "createdTime": created, "updatedTime": updated})}).Save()
	note := service.NoteS.GetNoteById(noteID)
	body := service.NoteS.GetNoteContent(noteID, user.Hex())
	if note.Title != "Hugo Title" || !note.CreatedTime.Equal(created) || !note.UpdatedTime.Equal(updated) || len(note.Tags) != 2 || body.Content != "# Body\n" || !body.UpdatedTime.Equal(updated) {
		t.Fatalf("metadata lost: %+v body=%+v", note, body)
	}
	if ok, _ := service.NotebookS.DeleteNotebook(user.Hex(), book.Hex()); ok {
		t.Fatal("deleted without explicit note deletion")
	}
	if ok, _ := service.NotebookS.DeleteNotebookAndNotes(bson.NewObjectId().Hex(), book.Hex()); ok {
		t.Fatal("foreign notebook deleted")
	}
	if ok, msg := service.NotebookS.DeleteNotebookAndNotes(user.Hex(), book.Hex()); !ok {
		t.Fatal(msg)
	}
	note = service.NoteS.GetNoteById(noteID)
	if !note.IsTrash || note.IsDeleted || service.NoteS.GetNoteContent(noteID, user.Hex()).Content != "# Body\n" || !service.NotebookS.GetNotebookById(book.Hex()).IsDeleted {
		t.Fatal("notes must remain recoverable in Trash")
	}
	empty := bson.NewObjectId()
	if !db.Insert(db.Notebooks, info.Notebook{NotebookId: empty, UserId: user, Title: "Empty", CreatedTime: time.Now(), UpdatedTime: time.Now()}) {
		t.Fatal("empty fixture")
	}
	(Notebook{makeController("/api2/notebook/deleteNotebook", nil)}).DeleteNotebook(empty.Hex(), true)
	if !service.NotebookS.GetNotebookById(empty.Hex()).IsDeleted {
		t.Fatal("empty notebook was not deleted")
	}
}
