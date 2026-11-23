package controllers

import (
	"net/url"
	"strings"
	"testing"

	"github.com/revel/config"
	"github.com/revel/revel"
)

func imageTestConfig(t *testing.T, secret string) {
	t.Helper()
	previous := revel.Config
	ctx := config.NewContext()
	ctx.SetOption("app.secret", secret)
	revel.Config = ctx
	t.Cleanup(func() { revel.Config = previous })
}

func imageRequest(params url.Values, fileId string) *revel.Controller {
	return &revel.Controller{Params: &revel.Params{Values: params}}
}

func TestSignedImagePathAuthorisesItsOwnViewer(t *testing.T) {
	imageTestConfig(t, strings.Repeat("s", 32))
	const user, file = "507f1f77bcf86cd799439011", "0123456789abcdef01234567"

	signed := SignedImagePath(file, user)
	if !strings.HasPrefix(signed, "/api2/file/getImage?fileId="+file+"&uid="+user+"&token=") {
		t.Fatalf("unexpected signed path %q", signed)
	}
	parsed, err := url.Parse(signed)
	if err != nil {
		t.Fatal(err)
	}
	if got := ImageCapabilityUser(imageRequest(parsed.Query(), file), file); got != user {
		t.Fatalf("capability resolved to %q, want %q", got, user)
	}
	if got := ImageCapabilityUser(imageRequest(parsed.Query(), "0123456789abcdef01234568"), "0123456789abcdef01234568"); got != "" {
		t.Fatal("capability accepted for another file")
	}
	tampered := parsed.Query()
	tampered.Set("uid", "507f1f77bcf86cd799439012")
	if got := ImageCapabilityUser(imageRequest(tampered, file), file); got != "" {
		t.Fatal("tampered viewer accepted")
	}
	if got := ImageCapabilityUser(imageRequest(url.Values{}, file), file); got != "" {
		t.Fatal("missing capability accepted")
	}
	if SignedImagePath("not-an-id", user) != "" || SignedImagePath(file, "") != "" {
		t.Fatal("invalid ids must not produce capabilities")
	}
}

func TestSignedImagePathRequiresSecret(t *testing.T) {
	imageTestConfig(t, "")
	if got := SignedImagePath("0123456789abcdef01234567", "507f1f77bcf86cd799439011"); got != "" {
		t.Fatalf("capability issued without a secret: %q", got)
	}
}
