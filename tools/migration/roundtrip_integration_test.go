package main

import (
	"os"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestRoundTripMigration(t *testing.T) {
	mongoURL := os.Getenv("GEMSNOTE_INTEGRATION_MONGO_URL")
	postgresURL := os.Getenv("GEMSNOTE_INTEGRATION_POSTGRES_URL")
	if mongoURL == "" {
		mongoURL = os.Getenv("LEANOTE_INTEGRATION_MONGO_URL")
	}
	if postgresURL == "" {
		postgresURL = os.Getenv("LEANOTE_INTEGRATION_POSTGRES_URL")
	}
	if mongoURL == "" || postgresURL == "" {
		t.Skip("set GEMSNOTE_INTEGRATION_MONGO_URL and GEMSNOTE_INTEGRATION_POSTGRES_URL")
	}
	m, err := connect(mongoURL, postgresURL, false)
	if err != nil {
		t.Fatal(err)
	}
	defer m.close()

	userID, notebookID, noteID := bson.NewObjectId(), bson.NewObjectId(), bson.NewObjectId()
	now := time.Now().UTC().Truncate(time.Microsecond)
	documents := map[string]bson.M{
		"users": {
			"_id": userID, "Email": userID.Hex() + "@migration.test", "Username": userID.Hex(),
			"UsernameRaw": userID.Hex(), "CreatedTime": now, "Verified": true,
		},
		"notebooks": {
			"_id": notebookID, "UserId": userID, "Title": "Round trip", "CreatedTime": now, "UpdatedTime": now,
		},
		"notes": {
			"_id": noteID, "UserId": userID, "NotebookId": notebookID, "Title": "Portable",
			"Desc": "Mongo to PostgreSQL and back", "Tags": []string{"mongo", "postgres"},
			"CreatedTime": now, "UpdatedTime": now,
		},
	}
	for collection, document := range documents {
		if err := m.mongo.C(collection).Insert(document); err != nil {
			t.Fatal(err)
		}
		defer m.mongo.C(collection).RemoveId(document["_id"])
	}
	defer func() {
		m.postgres.Exec("DELETE FROM notes WHERE id=$1", noteID.Hex())
		m.postgres.Exec("DELETE FROM notebooks WHERE id=$1", notebookID.Hex())
		m.postgres.Exec("DELETE FROM users WHERE id=$1", userID.Hex())
	}()

	for _, collection := range collections {
		if _, err := m.mongoToPostgres(collection); err != nil {
			t.Fatalf("mongo_to_pg %s: %v", collection, err)
		}
	}
	var storedID string
	if err := m.postgres.QueryRow("SELECT id FROM notes WHERE id=$1", noteID.Hex()).Scan(&storedID); err != nil {
		t.Fatal(err)
	}
	if storedID != noteID.Hex() {
		t.Fatalf("PostgreSQL id = %q, want %q", storedID, noteID.Hex())
	}

	for collection, document := range documents {
		if err := m.mongo.C(collection).RemoveId(document["_id"]); err != nil {
			t.Fatal(err)
		}
	}
	for _, collection := range collections {
		if _, err := m.postgresToMongo(collection); err != nil {
			t.Fatalf("pg_to_mongo %s: %v", collection, err)
		}
	}
	var note bson.M
	if err := m.mongo.C("notes").FindId(noteID).One(&note); err != nil {
		t.Fatal(err)
	}
	if note["UserId"] != userID || note["NotebookId"] != notebookID {
		t.Fatalf("relationships changed during round trip: %#v", note)
	}
	if err := m.validateCounts(); err != nil {
		t.Fatal(err)
	}
}
