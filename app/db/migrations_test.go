package db

import (
	"strings"
	"testing"
)

type migrationTestDatabase struct {
	Database
	applied  []string
	recorded []string
}

func (d *migrationTestDatabase) AppliedMigrations() ([]string, error) {
	return append([]string{}, d.applied...), nil
}
func (d *migrationTestDatabase) RecordMigration(id string) error {
	d.applied = append(d.applied, id)
	d.recorded = append(d.recorded, id)
	return nil
}

func TestMigrationsIndependentOfReleaseVersion(t *testing.T) {
	// Earlier development builds stored these numeric IDs. They remain valid
	// even though the application release stays at 1.1.0.
	database := &migrationTestDatabase{applied: []string{"1.0.0", "1.1.0", "1.2.0", "1.3.0"}}
	if err := migrate(database); err != nil {
		t.Fatal(err)
	}
	if len(database.recorded) != 1 || database.recorded[0] != "shared-note-placements" {
		t.Fatalf("unexpected migrations: %v", database.recorded)
	}
	if err := migrate(database); err != nil {
		t.Fatal(err)
	}
	if len(database.recorded) != 1 {
		t.Fatal("migration ran twice")
	}
}
func TestUnknownDatabaseMigrationRejected(t *testing.T) {
	database := &migrationTestDatabase{applied: []string{"future-schema-change"}}
	if err := migrate(database); err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("unexpected result: %v", err)
	}
	if len(database.recorded) != 0 {
		t.Fatal("modified unsupported database")
	}
}
func TestFreshDatabaseRunsMigrationsInDeclarationOrder(t *testing.T) {
	database := &migrationTestDatabase{}
	if err := migrate(database); err != nil {
		t.Fatal(err)
	}
	if len(database.recorded) != len(migrations) {
		t.Fatalf("missing migrations: %v", database.recorded)
	}
	for i, item := range migrations {
		if database.recorded[i] != item.id {
			t.Fatalf("wrong order: %v", database.recorded)
		}
	}
}
