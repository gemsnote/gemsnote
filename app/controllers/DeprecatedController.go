package controllers

import "github.com/revel/revel"

// Deprecated handles Web features deliberately removed from Gemsnote.
type Deprecated struct{ BaseController }

func (c Deprecated) Gone() revel.Result {
	c.Response.Status = 410
	return c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "featureRemoved"})
}
