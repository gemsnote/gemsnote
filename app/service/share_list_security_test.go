package service

import (
	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestSharedListRejectsForeignAndDeletedNotes(t *testing.T) {
	owner := bson.NewObjectId()
	note := info.Note{NoteId: bson.NewObjectId(), UserId: owner}
	if !shareListNoteAllowed(note, owner.Hex()) {
		t.Fatal("valid shared note rejected")
	}
	if shareListNoteAllowed(note, bson.NewObjectId().Hex()) {
		t.Fatal("forged grant exposes foreign note")
	}
	note.IsDeleted = true
	if shareListNoteAllowed(note, owner.Hex()) {
		t.Fatal("deleted note exposed")
	}
	note.IsDeleted = false
	note.IsTrash = true
	if shareListNoteAllowed(note, owner.Hex()) {
		t.Fatal("trashed note exposed")
	}
	if shareListNoteAllowed(info.Note{}, owner.Hex()) {
		t.Fatal("missing note exposed")
	}
}
