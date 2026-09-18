package route

import (
	"github.com/gemsnote/gemsnote/app/db"
	"github.com/revel/revel"
	// . "github.com/gemsnote/gemsnote/app/lea"
	"net/url"
	"strings"
)

// overwite revel RouterFilter
// /api/user/Info => ApiUser.Info()
var staticPrefix = []string{"/public", "/favicon.ico", "/css", "/js", "/images", "/tinymce", "/upload", "/fonts"}

func RouterFilter(c *revel.Controller, fc []revel.Filter) {
	// 补全controller部分
	path := c.Request.URL.Path
	// Public publishing has been retired, including old controller aliases.
	lowerPath := strings.ToLower(path)
	if removedWebPath(lowerPath) {
		c.Response.Status = 410
		c.Result = c.RenderJSON(map[string]interface{}{"Ok": false, "Msg": "featureRemoved"})
		return
	}

	// Figure out the Controller/Action
	// var route *revel.RouteMatch = revel.MainRouter.Route(c.Request.Request)
	route := revel.MainRouter.Route(c.Request)

	if route == nil {
		c.Result = c.NotFound("No matching route found: " + c.Request.GetRequestURI())
		return
	}

	// Log("---------" + route.Action + "  " + path)

	// The route may want to explicitly return a 404.
	if route.Action == "404" {
		c.Result = c.NotFound("(intentionally)")
		return
	}

	//----------
	// life start
	/*
		type URL struct {
		    Scheme   string
		    Opaque   string    // encoded opaque data
		    User     *Userinfo // username and password information
		    Host     string    // host or host:port
		    Path     string
		    RawQuery string // encoded query values, without '?'
		    Fragment string // fragment for references, without '#'
		}
	*/
	if route.ControllerName != "Static" {

		// 检查数据库连接是否丢失
		db.CheckConnection()

		// /api/file/getImage -> App\file (/api/file/getImage)
		// App\auth
		// App\note
		// static\static
		//
		// Log("---------" + route.ControllerName + "  " + path)

		// api设置
		// gemsnote.com/api/user/get => ApiUser::Get
		//*       /api/login               ApiAuth.Login,  这里的设置, 其实已经转成了ApiAuth了
		// Only the original /api/ namespace uses the Leanote controller-name
		// rewrite. /api2/ has explicit routes, and rewriting it would send
		// Note.DeleteNote to the unrelated ApiNote controller (404).
		if isLegacyAPIPath(path) && !strings.HasPrefix(route.ControllerName, "App\\api") {
			route.ControllerName = "App\\api" + strings.Split(route.ControllerName, "\\")[1]
			// route.ControllerName = "App\\apifile"
		} else if strings.HasPrefix(path, "/member/") && !strings.HasSuffix(strings.ToLower(route.ControllerName), "\\web") && !strings.HasSuffix(strings.ToLower(route.ControllerName), "\\deprecated") && !strings.HasPrefix(route.ControllerName, "App\\member") {
			// member设置
			// route.ControllerName = "App\\Member" + route.ControllerName
			route.ControllerName = "App\\member" + strings.Split(route.ControllerName, "\\")[1]
		}
		// end
	}

	// Set the action.
	if err := c.SetAction(route.ControllerName, route.MethodName); err != nil {
		c.Result = c.NotFound(err.Error())
		return
	}

	// Add the route and fixed params to the Request Params.
	c.Params.Route = route.Params

	// Add the fixed parameters mapped by name.
	// TODO: Pre-calculate this mapping.
	for i, value := range route.FixedParams {
		if c.Params.Fixed == nil {
			c.Params.Fixed = make(url.Values)
		}
		if i < len(c.MethodType.Args) {
			arg := c.MethodType.Args[i]
			c.Params.Fixed.Set(arg.Name, value)
		} else {
			// revel.WARN.Println("Too many parameters to", route.Action, "trying to add", value)
			break
		}
	}

	fc[0](c, fc[1:])
}

func isLegacyAPIPath(path string) bool {
	return strings.HasPrefix(path, "/api/")
}

func removedWebPath(path string) bool {
	for _, prefix := range []string{"/blog", "/preview", "/member/blog", "/memberblog", "/admin/blog", "/adminblog"} {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}
