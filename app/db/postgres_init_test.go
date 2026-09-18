package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindInstallationDatabaseDirHonorsEnvironment(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "schema.sql"), []byte("-- schema"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "seed.sql"), []byte("-- seed"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := findInstallationDatabaseDirFrom(dir, t.TempDir(), filepath.Join(t.TempDir(), "bin", "gemsnote"))
	if err != nil {
		t.Fatal(err)
	}
	want, _ := filepath.Abs(dir)
	if got != want {
		t.Fatalf("database directory = %q, want %q", got, want)
	}
}

func TestFindInstallationDatabaseDirRejectsIncompleteDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "schema.sql"), []byte("-- schema"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := findInstallationDatabaseDirFrom(dir, t.TempDir(), filepath.Join(t.TempDir(), "bin", "gemsnote")); err == nil {
		t.Fatal("expected missing seed.sql to be rejected")
	}
}

func TestFindInstallationDatabaseDirPrefersReleaseOverCwd(t *testing.T) {
	releaseRoot := t.TempDir()
	if err := os.Mkdir(filepath.Join(releaseRoot, "database"), 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"schema.sql", "seed.sql"} {
		if err := os.WriteFile(filepath.Join(releaseRoot, "database", name), []byte("-- "+name), 0600); err != nil {
			t.Fatal(err)
		}
	}
	cwd := t.TempDir()
	if err := os.Mkdir(filepath.Join(cwd, "database"), 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"schema.sql", "seed.sql"} {
		if err := os.WriteFile(filepath.Join(cwd, "database", name), []byte("-- cwd "+name), 0600); err != nil {
			t.Fatal(err)
		}
	}

	got, err := findInstallationDatabaseDirFrom("", cwd, filepath.Join(releaseRoot, "bin", "gemsnote"))
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(releaseRoot, "database")
	if got != want {
		t.Fatalf("database directory = %q, want bundled release path %q", got, want)
	}
}

func TestFindInstallationDatabaseDirFallsBackToCwd(t *testing.T) {
	cwd := t.TempDir()
	databaseDir := filepath.Join(cwd, "database")
	if err := os.Mkdir(databaseDir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"schema.sql", "seed.sql"} {
		if err := os.WriteFile(filepath.Join(databaseDir, name), []byte("-- "+name), 0600); err != nil {
			t.Fatal(err)
		}
	}
	got, err := findInstallationDatabaseDirFrom("", cwd, filepath.Join(t.TempDir(), "bin", "gemsnote"))
	if err != nil {
		t.Fatal(err)
	}
	if got != databaseDir {
		t.Fatalf("database directory = %q, want cwd fallback %q", got, databaseDir)
	}
}

func TestStripScriptTransactionBoundaries(t *testing.T) {
	got := stripScriptTransactionBoundaries("BEGIN;\nCREATE TABLE example (id integer);\nCOMMIT;\n")
	want := "\nCREATE TABLE example (id integer);\n"
	if got != want {
		t.Fatalf("stripped script = %q, want %q", got, want)
	}
}

func TestPostgresBootstrapLockIsProjectScoped(t *testing.T) {
	if postgresBootstrapLockSQL == "" || postgresBootstrapUnlockSQL == "" {
		t.Fatal("bootstrap advisory lock SQL must be configured")
	}
	if !strings.Contains(postgresBootstrapLockSQL, "gemsnote:postgres-bootstrap:v1") ||
		!strings.Contains(postgresBootstrapUnlockSQL, "gemsnote:postgres-bootstrap:v1") {
		t.Fatalf("bootstrap lock is not project scoped: %q / %q", postgresBootstrapLockSQL, postgresBootstrapUnlockSQL)
	}
}
