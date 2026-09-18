package db

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPostgresInstallationSeed(t *testing.T) {
	seedPath := filepath.Join("..", "..", "database", "seed.sql")
	data, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatal(err)
	}
	seed := string(data)

	for _, required := range []string{
		"INSERT INTO public.users",
		"'admin@gemsnote.com'",
		"'4533cd5ac65d7b5ccb02f396ff4b6615'",
		"INSERT INTO public.configs",
		"INSERT INTO public.notebooks",
		"INSERT INTO public.notes",
		"'love gemsnote'",
		"ON CONFLICT DO NOTHING",
	} {
		if !strings.Contains(seed, required) {
			t.Errorf("installation seed is missing %q", required)
		}
	}
	for _, removed := range []string{"'demo@gemsnote.com'", "'demoUserId'", "'demoUsername'", "'demoPassword'"} {
		if strings.Contains(seed, removed) {
			t.Errorf("installation seed contains removed demo data %s", removed)
		}
	}

	for _, excludedTable := range []string{"sessions", "tokens", "email_logs", "suggestions", "reports"} {
		if strings.Contains(seed, "INSERT INTO public."+excludedTable) {
			t.Errorf("installation seed contains runtime-only table %s", excludedTable)
		}
	}
}

func TestMongoInstallationGroupName(t *testing.T) {
	path := filepath.Join("..", "..", "mongodb_backup", "gemsnote_install_data", "groups.bson")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("love gemsnote")) {
		t.Fatal("MongoDB installation seed is missing the Gemsnote group name")
	}
}
