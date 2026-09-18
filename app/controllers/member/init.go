package member

import (
	"github.com/gemsnote/gemsnote/app/info"
	"github.com/gemsnote/gemsnote/app/service"
	//	. "github.com/gemsnote/gemsnote/app/lea"
	"github.com/revel/revel"
	//	"strings"
)

var userService *service.UserService
var groupService *service.GroupService
var noteService *service.NoteService
var trashService *service.TrashService
var notebookService *service.NotebookService
var noteContentHistoryService *service.NoteContentHistoryService
var authService *service.AuthService
var shareService *service.ShareService
var blogService *service.BlogService
var tagService *service.TagService
var pwdService *service.PwdService
var tokenService *service.TokenService
var suggestionService *service.SuggestionService
var albumService *service.AlbumService
var noteImageService *service.NoteImageService
var fileService *service.FileService
var attachService *service.AttachService
var configService *service.ConfigService
var emailService *service.EmailService
var upgradeService *service.UpgradeService
var themeService *service.ThemeService
var sessionService *service.SessionService

// 拦截器
// 不需要拦截的url
// Index 除了Note之外都不需要
var commonUrl = map[string]map[string]bool{"Index": map[string]bool{"Index": true,
	"Login":              true,
	"DoLogin":            true,
	"Logout":             true,
	"Register":           true,
	"DoRegister":         true,
	"FindPasswword":      true,
	"DoFindPassword":     true,
	"FindPassword2":      true,
	"FindPasswordUpdate": true,
	"Suggestion":         true,
},
	"Blog": map[string]bool{"Index": true,
		"View":       true,
		"AboutMe":    true,
		"SearchBlog": true,
	},
	// 用户的激活与修改邮箱都不需要登录, 通过链接地址
	"User": map[string]bool{"UpdateEmail": true,
		"ActiveEmail": true,
	},
	"Oauth":  map[string]bool{"GithubCallback": true},
	"File":   map[string]bool{"OutputImage": true, "OutputFile": true},
	"Attach": map[string]bool{"Download": true, "DownloadAll": true},
}

func needValidate(controller, method string) bool {
	// 在里面
	if v, ok := commonUrl[controller]; ok {
		// 在commonUrl里
		if _, ok2 := v[method]; ok2 {
			return false
		}
		return true
	} else {
		// controller不在这里的, 肯定要验证
		return true
	}
}
func AuthInterceptor(c *revel.Controller) revel.Result {
	// 全部变成首字大写
	/*
		var controller = strings.Title(c.Name)
		var method = strings.Title(c.MethodName)
		// 是否需要验证?
		if !needValidate(controller, method) {
			return nil
		}
	*/

	// 验证是否已登录
	// 必须是管理员
	userId, ok := c.Session["UserId"].(string)
	if ok && sessionService.ValidateUserSession(c.Session.ID(), userId) {
		return nil // 已登录
	}
	delete(c.Session, "UserId")
	delete(c.Session, "Username")

	// 没有登录, 判断是否是ajax操作
	if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		re := info.NewRe()
		re.Msg = "NOTLOGIN"
		return c.RenderJSON(re)
	}

	return c.Redirect("/login")
}

func MutationInterceptor(c *revel.Controller) revel.Result {
	if c.MethodName == "Index" {
		return nil
	}
	if c.Request.Method != "POST" {
		c.Response.Status = 405
		return c.RenderJSON(info.Re{Ok: false, Msg: "methodNotAllowed"})
	}
	if c.Request.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		c.Response.Status = 403
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidRequest"})
	}
	return nil
}

// 最外层init.go调用
// 获取service, 单例
func InitService() {
	notebookService = service.NotebookS
	noteService = service.NoteS
	noteContentHistoryService = service.NoteContentHistoryS
	trashService = service.TrashS
	shareService = service.ShareS
	userService = service.UserS
	groupService = service.GroupS
	tagService = service.TagS
	blogService = service.BlogS
	tokenService = service.TokenS
	noteImageService = service.NoteImageS
	fileService = service.FileS
	albumService = service.AlbumS
	attachService = service.AttachS
	pwdService = service.PwdS
	suggestionService = service.SuggestionS
	authService = service.AuthS
	configService = service.ConfigS
	emailService = service.EmailS
	upgradeService = service.UpgradeS
	themeService = service.ThemeS
	sessionService = service.SessionS
}

func init() {
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &MemberIndex{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &MemberUser{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &MemberBlog{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &MemberGroup{})
	revel.InterceptFunc(MutationInterceptor, revel.BEFORE, &MemberGroup{})
	revel.OnAppStart(func() {
	})
}
