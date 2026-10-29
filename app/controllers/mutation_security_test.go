package controllers

import (
	"github.com/revel/revel"
	"testing"
)

func TestLegacyMutationAliasesCannotUseGETOrCrossSiteForms(t *testing.T) {
	for _, target := range [][2]string{{"Web", "AdminSettings"}, {"Web", "Save"}, {"Tag", "DeleteTag"}, {"Album", "DeleteAlbum"}, {"File", "CopyImage"}} {
		c := cookieAuthenticatedController("/"+target[0]+"/"+target[1], "")
		c.Name, c.MethodName = target[0], target[1]
		c.Response = &revel.Response{}
		c.Request.Method = "GET"
		if WebMutationInterceptor(c) == nil || c.Response.Status != 405 {
			t.Fatalf("unsafe GET allowed: %v", target)
		}
		c.Request.Method = "POST"
		if WebMutationInterceptor(c) == nil || c.Response.Status != 403 {
			t.Fatalf("cross-site form allowed: %v", target)
		}
		c.Request.Header.Server.(*revel.GoHeader).Source.(*revel.GoRequest).Original.Header.Set("X-Requested-With", "XMLHttpRequest")
		if WebMutationInterceptor(c) != nil {
			t.Fatalf("same-origin POST blocked: %v", target)
		}
	}
}

func TestPublicAPI2DownloadClearsCookieIdentity(t *testing.T) {
	c := cookieAuthenticatedController("/api2/attach/download", "")
	c.Name, c.MethodName = "Attach", "Download"
	if AuthInterceptor(c) != nil {
		t.Fatal("public file endpoint should perform its own authorization")
	}
	if (BaseController{c}).GetUserId() != "" {
		t.Fatal("public download inherited cookie identity")
	}
}
