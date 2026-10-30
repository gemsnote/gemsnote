package db

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type migration struct {
	id string
	up func(Database) error
}

// Migrations run in declaration order and have stable IDs independent of the
// application release. Keep the existing numeric IDs (including 1.2.0/1.3.0)
// so databases already containing them remain usable with release 1.1.0.
// New migrations use descriptive IDs; no release bump is needed for a schema change.
var migrations = []migration{
	{id: "1.0.0", up: func(Database) error { return nil }},
	{id: "1.1.0", up: func(database Database) error {
		if postgres, ok := database.(*PostgresDatabase); ok {
			_, err := postgres.db.Exec(`ALTER TABLE files
				ADD COLUMN IF NOT EXISTS upload_state TEXT NOT NULL DEFAULT '',
				ADD COLUMN IF NOT EXISTS upload_note_id CHAR(24),
				ADD COLUMN IF NOT EXISTS upload_digest TEXT NOT NULL DEFAULT '',
				ADD COLUMN IF NOT EXISTS upload_key TEXT NOT NULL DEFAULT ''`)
			if err != nil {
				return err
			}
			_, err = postgres.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_files_upload_key ON files(upload_key) WHERE upload_key <> ''`)
			return err
		}
		if mongo, ok := database.(*MongoDatabase); ok {
			return mongo.db.C("files").EnsureIndex(mgo.Index{Key: []string{"UploadKey"}, Unique: true, Sparse: true})
		}
		return nil // MongoDB documents need no backfill; absent state is protected.
	}},
	{id: "1.2.0", up: func(database Database) error {
		if postgres, ok := database.(*PostgresDatabase); ok {
			_, err := postgres.db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT ''`)
			return err
		}
		return nil
	}},
	{id: "1.3.0", up: func(database Database) error {
		if postgres, ok := database.(*PostgresDatabase); ok {
			_, err := postgres.db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS flag TEXT NOT NULL DEFAULT 'A' CHECK (flag IN ('A','D'))`)
			return err
		}
		if mongo, ok := database.(*MongoDatabase); ok {
			_, err := mongo.db.C("users").UpdateAll(bson.M{"$or": []bson.M{{"Flag": bson.M{"$exists": false}}, {"Flag": ""}}}, bson.M{"$set": bson.M{"Flag": "A"}})
			return err
		}
		return nil
	}},
	{id: "shared-note-placements", up: func(database Database) error {
		if postgres, ok := database.(*PostgresDatabase); ok {
			_, err := postgres.db.Exec(`CREATE TABLE IF NOT EXISTS shared_note_placements (
                user_id CHAR(24) NOT NULL, note_id CHAR(24) NOT NULL,
                notebook_id CHAR(24) NOT NULL, PRIMARY KEY(user_id,note_id))`)
			return err
		}
		return nil
	}},
}

func migrate(database Database) error {
	appliedIDs, err := database.AppliedMigrations()
	if err != nil {
		return fmt.Errorf("read database migrations: %w", err)
	}
	supported := make(map[string]bool, len(migrations))
	for _, item := range migrations {
		if item.id == "" || supported[item.id] {
			return fmt.Errorf("invalid or duplicate migration ID %q", item.id)
		}
		supported[item.id] = true
	}
	applied := make(map[string]bool, len(appliedIDs))
	for _, id := range appliedIDs {
		if !supported[id] {
			return fmt.Errorf("database migration %s is not supported by this application", id)
		}
		applied[id] = true
	}
	for _, item := range migrations {
		if applied[item.id] {
			continue
		}
		if err := item.up(database); err != nil {
			return fmt.Errorf("apply database migration %s: %w", item.id, err)
		}
		if err := database.RecordMigration(item.id); err != nil {
			return fmt.Errorf("record database migration %s: %w", item.id, err)
		}
	}
	return nil
}
