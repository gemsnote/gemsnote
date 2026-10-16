package db

import (
	"encoding/json"
	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestHistoryArchiveSurvivesPostgresAndBSONEncoding(t *testing.T) {
	id := bson.NewObjectId()
	original := []info.EachHistory{{HistoryId: bson.NewObjectId(), UpdatedUserId: id, FilesComplete: true, Content: "", AttachmentRecords: []info.Attach{{AttachId: id, UploadUserId: id, NoteId: bson.NewObjectId(), Path: "retained-original"}}}}
	raw, ok := postgresValue(original).([]byte)
	if !ok {
		t.Fatal("not JSON")
	}
	var decoded []info.EachHistory
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if !decoded[0].FilesComplete || decoded[0].AttachmentRecords[0].AttachId != id {
		t.Fatal(decoded)
	}
	mongo, err := bson.Marshal(info.NoteContentHistory{NoteId: id, UserId: id, Histories: original})
	if err != nil {
		t.Fatal(err)
	}
	var value info.NoteContentHistory
	if err = bson.Unmarshal(mongo, &value); err != nil {
		t.Fatal(err)
	}
	if !value.Histories[0].FilesComplete || value.Histories[0].AttachmentRecords[0].Path != "retained-original" {
		t.Fatal(value)
	}
}
