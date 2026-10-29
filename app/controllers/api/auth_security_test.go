package api

import (
	"github.com/gemsnote/gemsnote/app/controllers"
	"github.com/revel/revel"
	"testing"
)

func TestAPIIdentityCannotFallBackToAnotherCookie(t *testing.T) {
	c := &revel.Controller{Session: map[string]interface{}{"UserId": "cookie-owner", "_userId": "token-owner", "_token": "old-token"}, ViewArgs: map[string]interface{}{}}
	base := ApiBaseContrller{controllers.BaseController{Controller: c}}
	if got := base.getUserId(); got != "token-owner" {
		t.Fatalf("legacy token identity = %q", got)
	}
	c.ViewArgs["Api2UserId"] = ""
	c.ViewArgs["Api2Token"] = ""
	if base.getUserId() != "" || base.getToken() != "" {
		t.Fatal("anonymous API2 inherited stale legacy credentials")
	}
	c.ViewArgs["Api2UserId"] = "api2-owner"
	c.ViewArgs["Api2Token"] = "api2-token"
	if base.getUserId() != "api2-owner" || base.getToken() != "api2-token" {
		t.Fatal("API2 credential not authoritative")
	}
}

func TestCookieAPIReadsExcludeMutations(t *testing.T) {
	for _, action := range []string{"UpdateNote", "AddNote", "DeleteTrash", "RestoreHistory"} {
		if legacyCookieReadAction("ApiNote", action) {
			t.Fatalf("mutation %s allowed as cookie GET", action)
		}
	}
	if !legacyCookieReadAction("ApiNote", "GetNoteContent") {
		t.Fatal("legacy read disabled")
	}
}
