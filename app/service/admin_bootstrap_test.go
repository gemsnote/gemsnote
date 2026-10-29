package service

import (
	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestLegacyAdministratorConfigOwner(t *testing.T) {
	first, second := bson.NewObjectId(), bson.NewObjectId()
	for _, test := range []struct {
		name    string
		configs []info.Config
		want    string
	}{
		{"missing", nil, ""},
		{"renamed owner", []info.Config{{UserId: first, Key: "siteUrl"}, {UserId: first, Key: "emailHost"}}, first.Hex()},
		{"conflicting owners", []info.Config{{UserId: first}, {UserId: second}}, ""},
		{"missing owner", []info.Config{{UserId: first}, {}}, ""},
		{"invalid owner", []info.Config{{UserId: bson.ObjectId("invalid")}}, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := legacyConfigOwner(test.configs); got != test.want {
				t.Fatalf("owner=%q, want %q", got, test.want)
			}
		})
	}
}
