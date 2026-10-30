package service

import (
	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func TestLinkValidationDistinguishesFailureCauses(t *testing.T) {
	now := time.Now()
	user := info.User{UserId: bson.NewObjectId(), Flag: info.FlagAvailable}
	record := info.Token{UserId: user.UserId, Token: "test-proof", Type: info.TokenActiveEmail, CreatedTime: now.Add(-time.Minute)}
	for _, tc := range []struct {
		name string
		edit func(*info.Token, *info.User)
		want string
	}{
		{"valid", func(*info.Token, *info.User) {}, ""},
		{"missing or replaced", func(r *info.Token, _ *info.User) { *r = info.Token{} }, "verificationLinkInvalid"},
		{"deleted user", func(_ *info.Token, u *info.User) { *u = info.User{} }, "verificationLinkInvalid"},
		{"disabled", func(_ *info.Token, u *info.User) { u.Flag = info.FlagDisabled }, "userDisabled"},
		{"wrong purpose", func(r *info.Token, _ *info.User) { r.Type = info.TokenPwd }, "verificationLinkInvalid"},
		{"expired activation", func(r *info.Token, _ *info.User) { r.CreatedTime = now.Add(-49 * time.Hour) }, "verificationLinkExpired"},
		{"future record", func(r *info.Token, _ *info.User) { r.CreatedTime = now.Add(time.Minute) }, "verificationLinkNotYetValid"},
		{"expired setup", func(r *info.Token, _ *info.User) {
			r.Type = info.TokenInitialPassword
			r.CreatedTime = now.Add(-3 * time.Hour)
		}, "verificationLinkExpired"},
		{"activation still valid at 3 hours", func(r *info.Token, _ *info.User) { r.CreatedTime = now.Add(-3 * time.Hour) }, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, u := record, user
			tc.edit(&r, &u)
			if got := validateLinkRecord("test-proof", r, u, []int{info.TokenActiveEmail, info.TokenInitialPassword}, now); got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
	if got := validateLinkRecord("", info.Token{}, info.User{}, []int{info.TokenActiveEmail}, now); got != "verificationLinkMissing" {
		t.Fatal(got)
	}
}
