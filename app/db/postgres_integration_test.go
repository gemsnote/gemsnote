package db

import (
	"os"
	"testing"
	"time"

	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
)

func TestPostgresCRUDContract(t *testing.T) {
	dsn := os.Getenv("GEMSNOTE_INTEGRATION_POSTGRES_URL")
	if dsn == "" {
		t.Skip("set GEMSNOTE_INTEGRATION_POSTGRES_URL to run PostgreSQL integration tests")
	}
	previousURL, previousDB := initURL, initDBName
	initURL, initDBName = dsn, "gemsnote"
	defer func() { initURL, initDBName = previousURL, previousDB }()

	p := &PostgresDatabase{}
	if err := p.Initialize(); err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	userID, notebookID, noteID := bson.NewObjectId(), bson.NewObjectId(), bson.NewObjectId()
	now := time.Now().UTC().Truncate(time.Microsecond)
	user := info.User{UserId: userID, Email: userID.Hex() + "@example.test", Username: userID.Hex(), UsernameRaw: userID.Hex(), CreatedTime: now}
	if !p.Insert("users", user) {
		t.Fatal("insert user failed")
	}
	defer p.DeleteAll("users", bson.M{"_id": userID})

	notebook := info.Notebook{NotebookId: notebookID, UserId: userID, Title: "Contract", CreatedTime: now, UpdatedTime: now}
	if !p.Insert("notebooks", notebook) {
		t.Fatal("insert notebook failed")
	}
	defer p.DeleteAll("notebooks", bson.M{"_id": notebookID})

	note := info.Note{NoteId: noteID, UserId: userID, NotebookId: notebookID, Title: "Hello PostgreSQL", Tags: []string{"go", "db"}, CreatedTime: now, UpdatedTime: now}
	if !p.Insert("notes", note) {
		t.Fatal("insert note failed")
	}
	defer p.DeleteAll("notes", bson.M{"_id": noteID})

	var got info.Note
	p.Get("notes", noteID.Hex(), &got)
	if got.NoteId != noteID || got.UserId != userID || len(got.Tags) != 2 {
		t.Fatalf("ObjectId/array round trip failed: %#v", got)
	}

	query := bson.M{"UserId": userID, "$or": []bson.M{
		{"Title": bson.M{"$regex": bson.RegEx{Pattern: "postgres", Options: "i"}}},
		{"Tags": bson.M{"$all": []string{"missing"}}},
	}}
	var notes []info.Note
	p.ListByQOptions("notes", query, &notes, QueryOptions{Sort: []string{"-UpdatedTime"}, Limit: 5})
	if len(notes) != 1 || notes[0].NoteId != noteID {
		t.Fatalf("query contract failed: %#v", notes)
	}

	if !p.Update("notes", bson.M{"_id": noteID}, bson.M{"$inc": bson.M{"ReadNum": 2}}) {
		t.Fatal("atomic increment failed")
	}
	p.Get("notes", noteID.Hex(), &got)
	if got.ReadNum != 2 {
		t.Fatalf("read_num = %d, want 2", got.ReadNum)
	}

	if !p.Upsert("tags", bson.M{"_id": userID}, bson.M{"$addToSet": bson.M{"Tags": "one"}}) ||
		!p.Upsert("tags", bson.M{"_id": userID}, bson.M{"$addToSet": bson.M{"Tags": "two"}}) ||
		!p.Upsert("tags", bson.M{"_id": userID}, bson.M{"$addToSet": bson.M{"Tags": "two"}}) {
		t.Fatal("array upsert failed")
	}
	defer p.DeleteAll("tags", bson.M{"_id": userID})
	var tags info.Tag
	p.Get("tags", userID.Hex(), &tags)
	if len(tags.Tags) != 2 {
		t.Fatalf("addToSet produced %#v", tags.Tags)
	}

	share1, share2 := bson.NewObjectId(), bson.NewObjectId()
	if !p.Insert("share_notes", info.ShareNote{ShareNoteId: share1, UserId: userID, NoteId: noteID, CreatedTime: now}) ||
		!p.Insert("share_notes", info.ShareNote{ShareNoteId: share2, UserId: userID, NoteId: noteID, CreatedTime: now}) {
		t.Fatal("insert shares failed")
	}
	defer p.DeleteAll("share_notes", bson.M{"_id": bson.M{"$in": []bson.ObjectId{share1, share2}}})
	var distinctUsers []bson.ObjectId
	p.Distinct("share_notes", bson.M{"NoteId": noteID}, "UserId", &distinctUsers)
	if len(distinctUsers) != 1 || distinctUsers[0] != userID {
		t.Fatalf("distinct ObjectId scan failed: %#v", distinctUsers)
	}
}
