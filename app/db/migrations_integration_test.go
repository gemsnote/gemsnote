package db

import (
	"os"
	"testing"
)

func TestMongoMigrationVersionPersistence(t *testing.T) {
	url := os.Getenv("GEMSNOTE_INTEGRATION_MONGO_URL")
	if url == "" {
		t.Skip("set GEMSNOTE_INTEGRATION_MONGO_URL to run MongoDB migration tests")
	}
	previousURL, previousDB := initURL, initDBName
	initURL, initDBName = url, ""
	defer func() { initURL, initDBName = previousURL, previousDB }()

	mongo := &MongoDatabase{}
	if err := mongo.Initialize(); err != nil {
		t.Fatal(err)
	}
	defer mongo.Close()
	assertMigrationVersion(t, mongo)
}

func TestPostgresMigrationVersionPersistence(t *testing.T) {
	url := os.Getenv("GEMSNOTE_INTEGRATION_POSTGRES_URL")
	if url == "" {
		t.Skip("set GEMSNOTE_INTEGRATION_POSTGRES_URL to run PostgreSQL migration tests")
	}
	previousURL, previousDB := initURL, initDBName
	initURL, initDBName = url, ""
	defer func() { initURL, initDBName = previousURL, previousDB }()

	postgres := &PostgresDatabase{}
	if err := postgres.Initialize(); err != nil {
		t.Fatal(err)
	}
	defer postgres.Close()
	assertMigrationVersion(t, postgres)
}

func assertMigrationVersion(t *testing.T, database Database) {
	t.Helper()
	if err := migrate(database); err != nil {
		t.Fatal(err)
	}
	// A second run verifies that migrations and version recording are idempotent.
	if err := migrate(database); err != nil {
		t.Fatal(err)
	}
	versions, err := database.AppliedMigrations()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range migrations {
		count := 0
		for _, applied := range versions {
			if applied == item.id {
				count++
			}
		}
		if count != 1 {
			t.Fatalf("database has %d records for migration %s; migrations: %#v", count, item.id, versions)
		}
	}
}
