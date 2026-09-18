package api

import (
	appversion "github.com/gemsnote/gemsnote/app/version"
	"github.com/revel/revel"
)

// ApiSystem exposes public server metadata used by newer clients to decide
// whether the endpoint is Gemsnote-compatible or an older Leanote server.
type ApiSystem struct {
	ApiBaseContrller
}

func (c ApiSystem) Version() revel.Result {
	return c.RenderJSON(map[string]string{
		"server":      "gemsnote",
		"version":     appversion.Current,
		"min_version": "",
	})
}
