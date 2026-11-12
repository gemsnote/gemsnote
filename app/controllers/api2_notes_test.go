package controllers

import (
	"fmt"
	"testing"

	"github.com/revel/revel"
)

func TestAPI2NotesUsesJSONPage(t *testing.T) {
	for _, page := range []int{1, 2, 3} {
		t.Run(fmt.Sprint(page), func(t *testing.T) {
			c := cookieAuthenticatedController("/api2/notes?page=99", "")
			c.Params.JSON = []byte(fmt.Sprintf(`{"page":%d,"sort":"UpdatedTime"}`, page))
			c.Response = &revel.Response{Status: 200}
			api := Api2{BaseController{c}}
			api.Notes()
			if got := api.GetPage(); got != page {
				t.Fatalf("JSON page %d became %d; facade must update the parameters used by Bind", page, got)
			}
		})
	}
}
