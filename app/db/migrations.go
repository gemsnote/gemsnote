package db

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	appversion "github.com/gemsnote/gemsnote/app/version"
)

type migration struct {
	version string
	up      func(Database) error
}

// migrations must contain every database change in ascending version order.
// Version 1.0.0 establishes the migration baseline for existing Leanote and
// Gemsnote databases; it intentionally makes no changes to business data.
var migrations = []migration{
	{version: "1.0.0", up: func(Database) error { return nil }},
}

func migrate(database Database) error {
	appliedVersions, err := database.AppliedMigrations()
	if err != nil {
		return fmt.Errorf("read database version: %w", err)
	}
	applied := make(map[string]bool, len(appliedVersions))
	for _, migrationVersion := range appliedVersions {
		applied[migrationVersion] = true
		if compareVersions(migrationVersion, appversion.Current) > 0 {
			return fmt.Errorf("database version %s is newer than application version %s", migrationVersion, appversion.Current)
		}
	}

	ordered := append([]migration(nil), migrations...)
	sort.Slice(ordered, func(i, j int) bool {
		return compareVersions(ordered[i].version, ordered[j].version) < 0
	})
	for _, item := range ordered {
		if compareVersions(item.version, appversion.Current) > 0 || applied[item.version] {
			continue
		}
		if err := item.up(database); err != nil {
			return fmt.Errorf("migrate database to %s: %w", item.version, err)
		}
		if err := database.RecordMigration(item.version); err != nil {
			return fmt.Errorf("record database version %s: %w", item.version, err)
		}
	}
	return nil
}

func compareVersions(left, right string) int {
	leftParts := versionParts(left)
	rightParts := versionParts(right)
	for i := 0; i < 3; i++ {
		if leftParts[i] < rightParts[i] {
			return -1
		}
		if leftParts[i] > rightParts[i] {
			return 1
		}
	}
	return 0
}

func versionParts(value string) [3]int {
	value = strings.TrimPrefix(value, "v")
	value = strings.SplitN(value, "-", 2)[0]
	parts := strings.Split(value, ".")
	var result [3]int
	for i := 0; i < len(parts) && i < len(result); i++ {
		result[i], _ = strconv.Atoi(parts[i])
	}
	return result
}
