package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/revel/config"
	"github.com/revel/revel"
)

func clientIPController(remoteAddr, forwarded string) *revel.Controller {
	request := httptest.NewRequest("POST", "/doLogin", nil)
	request.RemoteAddr = remoteAddr
	if forwarded != "" {
		request.Header.Set("X-Forwarded-For", forwarded)
	}
	return &revel.Controller{Request: &revel.Request{
		RemoteAddr: remoteAddr,
		Header:     &revel.RevelHeader{Server: &revel.GoHeader{Source: &revel.GoRequest{Original: request}}},
	}}
}

func TestClientIPIgnoresForwardedHeaderByDefault(t *testing.T) {
	previous := revel.Config
	revel.Config = nil
	defer func() { revel.Config = previous }()

	if got := ClientIP(clientIPController("203.0.113.5:4444", "198.51.100.9")); got != "203.0.113.5" {
		t.Fatalf("unexpected client address %q", got)
	}
}

func TestClientIPTrustsForwardedHeaderOnlyWhenEnabled(t *testing.T) {
	previous := revel.Config
	ctx := config.NewContext()
	ctx.SetOption("login.trust.forwarded", "true")
	revel.Config = ctx
	defer func() { revel.Config = previous }()

	if got := ClientIP(clientIPController("10.0.0.4:4444", "198.51.100.9, 10.0.0.1")); got != "198.51.100.9" {
		t.Fatalf("unexpected forwarded address %q", got)
	}
	if got := ClientIP(clientIPController("10.0.0.4:4444", "not-an-ip")); got != "10.0.0.4" {
		t.Fatalf("invalid forwarded address accepted: %q", got)
	}
}
