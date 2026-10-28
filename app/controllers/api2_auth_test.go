package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/gemsnote/gemsnote/app/service"
	"github.com/revel/revel"
)

func cookieAuthenticatedController(path, authorization string) *revel.Controller {
	request := httptest.NewRequest("GET", path, nil)
	request.Header.Set("Authorization", authorization)
	goRequest := &revel.GoRequest{Original: request}
	return &revel.Controller{
		Request:  &revel.Request{URL: request.URL, Header: &revel.RevelHeader{Server: &revel.GoHeader{Source: goRequest}}},
		Params:   &revel.Params{Values: request.URL.Query()},
		ViewArgs: map[string]interface{}{},
		Session:  map[string]interface{}{"UserId": "507f1f77bcf86cd799439011", "Username": "cookie-user", "Email": "cookie@example.test"},
	}
}

func TestAPI2InvalidCredentialsCannotFallBackToCookie(t *testing.T) {
	for _, scenario := range []struct{ name, path, authorization string }{
		{"missing", "/api2/bootstrap", ""},
		{"expired", "/api2/bootstrap", "Bearer expired-token"},
		{"conflicting", "/api2/bootstrap?token=different-token", "Bearer expired-token"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			controller := cookieAuthenticatedController(scenario.path, scenario.authorization)
			userID, token := AuthenticateAPI2(controller, &service.SessionService{})
			base := BaseController{controller}
			if userID != "" || token != "" || base.GetUserId() != "" || base.GetApi2Token() != "" || base.GetUsername() != "" || base.GetEmail() != "" {
				t.Fatal("invalid API2 credentials inherited an identity from the browser cookie")
			}
		})
	}
}

func TestAPI2PublicBootstrapResolvesAuthentication(t *testing.T) {
	controller := cookieAuthenticatedController("/api2/bootstrap", "")
	controller.MethodName = "Bootstrap"
	api := Api2{BaseController{controller}}
	if result := api.RequireSession(); result != nil {
		t.Fatal("anonymous bootstrap must remain public")
	}
	if api.GetUserId() != "" {
		t.Fatal("bootstrap must not claim a cookie-only user is API2 authenticated")
	}
}

func TestLegacyCookieIdentityIsUnchanged(t *testing.T) {
	controller := cookieAuthenticatedController("/web/bootstrap", "")
	AuthenticateAPI2(controller, &service.SessionService{})
	base := BaseController{controller}
	if base.GetUserId() == "" || base.GetUsername() != "cookie-user" || base.GetEmail() != "cookie@example.test" {
		t.Fatal("legacy cookie identity changed")
	}
}
