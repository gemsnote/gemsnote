package info

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"testing"
)

func TestUserAvailability(t *testing.T) {
	id := bson.NewObjectId()
	for _, test := range []struct {
		flag      string
		available bool
	}{{"A", true}, {"", true}, {"D", false}, {"unknown", false}, {"a", false}} {
		if got := (User{UserId: id, Flag: test.flag}).Available(); got != test.available {
			t.Fatalf("flag %q availability=%v", test.flag, got)
		}
	}
	if (User{Flag: FlagAvailable}).Available() {
		t.Fatal("missing user must not be available")
	}
}
func TestUserFlagPreservesLegacyJSON(t *testing.T) {
	data, err := json.Marshal(User{UserId: bson.NewObjectId(), Flag: FlagDisabled})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"Flag"`) {
		t.Fatal("legacy user response gained an unexpected flag")
	}
}
