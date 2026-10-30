package service

import (
	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func TestPendingEmailRequiresCurrentValidChangeRecord(t *testing.T) {
	now := time.Now()
	user := info.User{UserId: bson.NewObjectId(), Email: "old@example.test", Flag: info.FlagAvailable}
	record := info.Token{UserId: user.UserId, Email: "new@example.test", Token: "verification", Type: info.TokenUpdateEmail, CreatedTime: now}
	if !pendingEmailValid(user, record, now) {
		t.Fatal("valid pending email missing")
	}
	for _, edit := range []func(*info.Token){
		func(r *info.Token) { r.UserId = bson.NewObjectId() },
		func(r *info.Token) { r.Type = info.TokenActiveEmail },
		func(r *info.Token) { r.Type = info.TokenPwd },
		func(r *info.Token) { r.CreatedTime = now.Add(-3 * time.Hour) },
		func(r *info.Token) { r.CreatedTime = now.Add(time.Minute) },
		func(r *info.Token) { r.Token = "" },
		func(r *info.Token) { r.Email = "" },
		func(r *info.Token) { r.Email = "OLD@example.test" },
	} {
		invalid := record
		edit(&invalid)
		if pendingEmailValid(user, invalid, now) {
			t.Fatalf("invalid record exposed: type=%d", invalid.Type)
		}
	}
	user.Flag = info.FlagDisabled
	if pendingEmailValid(user, record, now) {
		t.Fatal("disabled account exposed pending change")
	}
}
