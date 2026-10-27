package db

import (
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
)

func TestImageReferenceNormalization(t *testing.T) {
	id := bson.NewObjectId()
	for _, content := range []string{
		`![img](/api2/file/getImage?fileId=` + id.Hex() + `)`,
		`<img src="/api2/file/getImage?x=1&amp;fileId=` + id.Hex() + `">`,
		`<img src="/api2/file/getImage%3FfileId%3D` + strings.ToUpper(id.Hex()) + `">`,
	} {
		query, err := imageReferenceQuery(bson.M{"Content": content})
		if err != nil || query == nil {
			t.Fatalf("missing reference: %v", err)
		}
		found := false
		for _, value := range query["$or"].([]bson.M)[0]["_id"].(bson.M)["$in"].([]bson.ObjectId) {
			if value == id {
				found = true
			}
		}
		if !found {
			t.Fatal("reference normalization lost ID")
		}
	}
	query, err := imageReferenceQuery(bson.M{"Content": `<img src="https://notes.test/files/a/../b/image.png">`})
	if err != nil || query["$or"].([]bson.M)[0]["Path"].(bson.M)["$in"].([]string)[0] != "files/b/image.png" {
		t.Fatal("path normalization failed")
	}
}

func TestPostgresImageUploadLifecycle(t *testing.T) {
	url := os.Getenv("GEMSNOTE_INTEGRATION_POSTGRES_URL")
	if url == "" {
		t.Skip("set GEMSNOTE_INTEGRATION_POSTGRES_URL for isolated database tests")
	}
	previousURL, previousName, previous := initURL, initDBName, globalDatabase
	initURL, initDBName = url, ""
	database := &PostgresDatabase{}
	if err := database.Initialize(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })
	defer func() { initURL, initDBName, globalDatabase = previousURL, previousName, previous }()
	if err := migrate(database); err != nil {
		t.Fatal(err)
	}
	globalDatabase = database
	userID := bson.NewObjectId()
	newFile := func(state string) info.File {
		file := info.File{FileId: bson.NewObjectId(), UserId: userID, AlbumId: bson.NewObjectId(), CreatedTime: time.Now(), UploadNoteId: bson.NewObjectId(), UploadState: state, Path: "files/test/image.png", UploadDigest: strings.Repeat("a", 64)}
		if !database.Insert(imageCollection(database), file) {
			t.Fatal("insert receipt")
		}
		t.Cleanup(func() { database.Delete(imageCollection(database), bson.M{"_id": file.FileId}) })
		return file
	}
	// Both current content and history manifests claim the original first.
	for _, table := range []string{"notes", "note_contents", "note_content_histories", "note_images", "users"} {
		file := newFile("pending")
		if !protectImageReferences(database, table, bson.M{"Content": "/api2/file/getImage?fileId=" + file.FileId.Hex()}) {
			t.Fatal("claim failed")
		}
		if _, err := ReserveImageCleanup(userID.Hex(), file.FileId.Hex()); err == nil {
			t.Fatalf("deleted referenced image in %s", table)
		}
	}
	legacy := newFile("")
	if _, err := ReserveImageCleanup(userID.Hex(), legacy.FileId.Hex()); err == nil {
		t.Fatal("legacy image cleanup allowed")
	}
	keyed := newFile("pending")
	keyed.UploadKey = userID.Hex() + ":unique_upload_key"
	if !database.UpdateByQMap(imageCollection(database), bson.M{"_id": keyed.FileId}, bson.M{"UploadKey": keyed.UploadKey}) {
		t.Fatal("set upload key")
	}
	keyed.FileId = bson.NewObjectId()
	if database.Insert(imageCollection(database), keyed) {
		t.Fatal("duplicate upload key accepted")
	}
	file := newFile("pending")
	if _, err := ImageUpload(bson.NewObjectId().Hex(), file.FileId.Hex()); err == nil {
		t.Fatal("cross-owner receipt visible")
	}
	reserved, err := ReserveImageCleanup(userID.Hex(), file.FileId.Hex())
	if err != nil || !CompleteImageCleanup(reserved) {
		t.Fatalf("cleanup failed: %v", err)
	}
	if protectImageReferences(database, "note_contents", bson.M{"Content": file.FileId.Hex()}) {
		t.Fatal("accepted a deleted image")
	}
	if _, err := ReserveImageCleanup(userID.Hex(), file.FileId.Hex()); err != nil {
		t.Fatal("cleanup retry failed")
	}

	// Separate competing CAS operations; exactly one may win. This is not
	// serialized by a Go mutex and also works with multiple server processes.
	for i := 0; i < 24; i++ {
		file := newFile("pending")
		var claimed, cleaned bool
		start := make(chan struct{})
		var wait sync.WaitGroup
		wait.Add(2)
		go func() {
			defer wait.Done()
			<-start
			claimed = protectImageReferences(database, "note_contents", bson.M{"Content": file.FileId.Hex()})
		}()
		go func() {
			defer wait.Done()
			<-start
			_, err := ReserveImageCleanup(userID.Hex(), file.FileId.Hex())
			cleaned = err == nil
		}()
		close(start)
		wait.Wait()
		if claimed == cleaned {
			t.Fatalf("race invariant failed: claim=%v cleanup=%v", claimed, cleaned)
		}
	}
}
