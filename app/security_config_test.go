package app

import "testing"

func TestDatabaseSecretFromEnvironment(t *testing.T) {
	t.Setenv("GEMSNOTE_DB_URL", "")
	t.Setenv("GEMSNOTE_DB_PASSWORD", "")
	if databaseSecretFromEnvironment() {
		t.Fatal("empty environment reported as configured")
	}
	t.Setenv("GEMSNOTE_DB_PASSWORD", "strong-password")
	if !databaseSecretFromEnvironment() {
		t.Fatal("GEMSNOTE_DB_PASSWORD not detected")
	}
	t.Setenv("GEMSNOTE_DB_PASSWORD", "")
	t.Setenv("GEMSNOTE_DB_URL", "postgresql://user:pass@db/gemsnote")
	if !databaseSecretFromEnvironment() {
		t.Fatal("GEMSNOTE_DB_URL not detected")
	}
}
