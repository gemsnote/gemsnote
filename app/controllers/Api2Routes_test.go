package controllers

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"
)

// A Revel route declared with '*' lives under a separate method wildcard
// branch. Keep browser API2 actions bound to their actual HTTP method so they
// cannot silently start returning 404 when routes are added or reordered.
func TestAPI2BrowserRoutesHaveExplicitMethods(t *testing.T) {
	content, err := os.ReadFile("../../conf/routes")
	if err != nil {
		t.Fatal(err)
	}
	routes := make(map[string]bool)
	pattern := regexp.MustCompile(`^(GET|POST|\*)\s+(/api2/\S+)\s+`)
	for _, line := range strings.Split(string(content), "\n") {
		parts := pattern.FindStringSubmatch(line)
		if parts == nil {
			continue
		}
		if parts[1] == "*" {
			t.Errorf("API2 route must declare an HTTP method: %s", line)
		}
		routes[parts[1]+" "+parts[2]] = true
	}
	for _, route := range []string{
		"POST /api2/note/deleteNote",
		"POST /api2/note/deleteTrash",
		"POST /api2/note/moveNote",
		"POST /api2/note/copyNote",
		"POST /api2/share/listShareNotes",
		"POST /api2/notebook/addNotebook",
		"POST /api2/attachments",
		"POST /api2/noteContentHistory/listHistories",
		"POST /api2/avatar",
		"POST /api2/auth/login",
		"POST /api2/auth/session",
		"POST /api2/logout",
		"GET /api2/bootstrap",
		"GET /api2/groups",
		"GET /api2/user/info",
	} {
		if !routes[route] {
			t.Errorf("missing browser API2 route: %s", route)
		}
	}
	if !regexp.MustCompile(`(?m)^POST\s+/api2/auth/login\s+Api2\.TokenLogin\s*$`).Match(content) {
		t.Fatal("API2 token login must use the single-response TokenLogin action")
	}
}

func TestAPI2EmptyHistoriesAreArray(t *testing.T) {
	data, err := json.Marshal(api2HistoryArray(nil))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "[]" {
		t.Fatalf("empty API2 histories must be [], got %s", data)
	}
	content, err := os.ReadFile("../../conf/routes")
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`(?m)^POST\s+/api2/noteContentHistory/listHistories\s+Api2\.Histories\s*$`).Match(content) {
		t.Fatal("API2 history route must use the array-normalizing action")
	}
}
