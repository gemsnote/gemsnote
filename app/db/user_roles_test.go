package db

import (
	"encoding/json"
	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"testing"
)

func TestRoleChangeValidation(t *testing.T) {
	actor, target := bson.NewObjectId().Hex(), bson.NewObjectId().Hex()
	for _, test := range []struct{ actor, target, role, want string }{
		{actor, actor, info.RoleUser, "cannotDemoteSelf"},
		{actor, strings.ToUpper(actor), info.RoleUser, "cannotDemoteSelf"},
		{actor, actor, info.RoleAdmin, ""},
		{actor, target, info.RoleUser, ""},
		{actor, target, info.RoleAdmin, ""},
		{actor, target, "superuser", "invalidRole"},
		{actor, "invalid", info.RoleAdmin, "userNotExists"},
	} {
		if got := validateRoleChange(test.actor, test.target, test.role); got != test.want {
			t.Fatalf("role validation got %q, want %q", got, test.want)
		}
	}
}
func TestLegacyUserDoesNotExposeRole(t *testing.T) {
	encoded, err := json.Marshal(info.User{Role: info.RoleAdmin})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), `"Role"`) {
		t.Fatal("role changed the legacy user response")
	}
}
