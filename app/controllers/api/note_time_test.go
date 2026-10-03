package api

import (
	"testing"
	"time"
)

func TestAPI2ModificationTime(t *testing.T) {
	now := time.Date(2026, 10, 3, 1, 0, 0, 0, time.UTC)
	old := now.Add(-time.Hour)
	legacyNormalized := time.Date(old.Year(), old.Month(), old.Day(), old.Hour(), old.Minute(), old.Second(), 0, time.Local)
	for _, tc := range []struct {
		path           string
		supplied, want time.Time
	}{
		{"/api2/note/updateNote", time.Time{}, now},
		{"/api2/client/note/update", time.Time{}, now},
		{"/api2/note/addNote", time.Time{}, now},
		{"/api2/note/updateNote", old, legacyNormalized},
		{"/api/note/updateNote", time.Time{}, time.Time{}},
		{"/api/note/updateNote", old, old},
	} {
		if got := api2ModificationTime(tc.path, tc.supplied, now); !got.Equal(tc.want) {
			t.Errorf("%s: %v != %v", tc.path, got, tc.want)
		}
	}
}
