package db

import (
	"testing"

	appversion "github.com/gemsnote/gemsnote/app/version"
)

func TestMigrationListEndsAtCurrentVersion(t *testing.T) {
	if len(migrations) == 0 {
		t.Fatal("migration list is empty")
	}
	if got := migrations[len(migrations)-1].version; got != appversion.Current {
		t.Fatalf("latest migration is %s, application version is %s", got, appversion.Current)
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		left, right string
		want        int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.2.0", "1.10.0", -1},
		{"v2.0.0", "1.9.9", 1},
	}
	for _, test := range tests {
		if got := compareVersions(test.left, test.right); got != test.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", test.left, test.right, got, test.want)
		}
	}
}
