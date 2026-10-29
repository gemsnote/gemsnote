package controllers

import (
	"github.com/gemsnote/gemsnote/app/info"
	"github.com/gemsnote/gemsnote/app/service"
	//	. "github.com/gemsnote/gemsnote/app/lea"
	"github.com/revel/revel"
	"strings"
)

func isAPI2Request(c *revel.Controller) bool {
	return c.Request != nil && strings.HasPrefix(c.Request.URL.Path, "/api2/")
}

// AuthenticateAPI2 resolves the single API2 credential contract. Bearer is
// preferred; the token parameter remains accepted for older released clients.
func AuthenticateAPI2(c *revel.Controller, sessions *service.SessionService) (string, string) {
	if !isAPI2Request(c) {
		return "", ""
	}
	// Even public API2 actions must distinguish an anonymous token identity
	// from a legacy cookie. Otherwise bootstrap can advertise a login that
	// subsequent protected API2 calls will reject, causing a redirect loop.
	c.ViewArgs["Api2UserId"] = ""
	c.ViewArgs["Api2Token"] = ""
	token := ""
	authorization := strings.TrimSpace(c.Request.Header.Get("Authorization"))
	if len(authorization) > 7 && strings.EqualFold(authorization[:7], "Bearer ") {
		token = strings.TrimSpace(authorization[7:])
	}
	parameterToken := c.Params.Values.Get("token")
	if token == "" {
		token = parameterToken
	} else if parameterToken != "" && parameterToken != token {
		return "", ""
	}
	userID := sessions.FindUserId(token)
	if userID == "" {
		return "", ""
	}
	c.ViewArgs["Api2Token"] = token
	c.ViewArgs["Api2UserId"] = userID
	return userID, token
}

var userService *service.UserService
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
var sessionService *service.SessionService
var themeService *service.ThemeService

var pageSize = 1000
var defaultSortField = "UpdatedTime"

// 拦截器
// 不需要拦截的url
// Index 除了Note之外都不需要
var commonUrl = map[string]map[string]bool{"Index": map[string]bool{"Index": true,
	"Login":              true,
	"DoLogin":            true,
	"Logout":             true,
	"FindPasswword":      true,
	"DoFindPassword":     true,
	"FindPassword2":      true,
	"FindPasswordUpdate": true,
	"Suggestion":         true,
},
	"Note": map[string]bool{"ToPdf": true},
	"Blog": map[string]bool{"Index": true,
		"View":               true,
		"AboutMe":            true,
		"Cate":               true,
		"ListCateLatest":     true,
		"Search":             true,
		"GetLikeAndComments": true,
		"IncReadNum":         true,
		"ListComments":       true,
		"Single":             true,
		"Archive":            true,
		"Tags":               true,
	},
	// 用户的激活与修改邮箱都不需要登录, 通过链接地址
	"User": map[string]bool{"UpdateEmail": true,
		"ActiveEmail": true,
	},
	"Oauth":  map[string]bool{"GithubCallback": true},
	"File":   map[string]bool{"OutputImage": true, "OutputFile": true},
	"Attach": map[string]bool{"Download": true /*, "DownloadAll": true*/},
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
	var controller = strings.Title(c.Name)
	var method = strings.Title(c.MethodName)

	if isAPI2Request(c) {
		userID, _ := AuthenticateAPI2(c, sessionService)
		if !needValidate(controller, method) || userID != "" {
			return nil
		}
		c.Response.Status = 401
		return c.RenderJSON(info.Re{Ok: false, Msg: "NOTLOGIN"})
	}

	// Public file actions still require a validated optional cookie identity.
	// Revoked cookies must never authorize a download.
	if userID, ok := c.Session["UserId"].(string); ok &&
		sessionService.ValidateUserSession(c.Session.ID(), userID) {
		return nil
	}
	delete(c.Session, "UserId")
	delete(c.Session, "Username")
	if !needValidate(controller, method) {
		return nil
	}

	// 没有登录, 判断是否是ajax操作
	if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		re := info.NewRe()
		re.Msg = "NOTLOGIN"
		return c.RenderJSON(re)
	}

	return c.Redirect("/login")
}

var webMutationMethods = map[string]map[string]bool{
	"Tag":      {"UpdateTag": true, "DeleteTag": true},
	"Album":    {"AddAlbum": true, "UpdateAlbum": true, "DeleteAlbum": true},
	"Web":      {"Save": true, "Star": true, "Restore": true, "EmailChange": true, "AdminSettings": true, "AdminResetPwd": true, "AdminRegister": true},
	"Notebook": {"DeleteNotebook": true, "AddNotebook": true, "UpdateNotebookTitle": true, "MoveNotebook": true, "DragNotebooks": true, "SetNotebook2Blog": true},
	"Note":     {"UpdateNoteOrContent": true, "DeleteNote": true, "DeleteTrash": true, "MoveNote": true, "CopyNote": true, "CopySharedNote": true, "SetNote2Blog": true},
	"Share":    {"AddShareNote": true, "AddShareNotebook": true, "UpdateShareNotePerm": true, "UpdateShareNotebookPerm": true, "DeleteShareNote": true, "DeleteShareNotebook": true, "DeleteShareNoteBySharedUser": true, "DeleteShareNotebookBySharedUser": true, "DeleteUserShareNoteAndNotebook": true, "AddShareNoteGroup": true, "DeleteShareNoteGroup": true, "UpdateShareNoteGroupPerm": true, "AddShareNotebookGroup": true, "DeleteShareNotebookGroup": true, "UpdateShareNotebookGroupPerm": true},
	"User":     {"UpdateUsername": true, "UpdatePwd": true, "UpdateTheme": true, "ReSendActiveEmail": true, "UpdateColumnWidth": true, "UpdateLeftIsMin": true},
	"File":     {"UploadBlogLogo": true, "PasteImage": true, "UploadAvatar": true, "UploadImageLeaui": true, "UpdateImageTitle": true, "DeleteImage": true, "CopyImage": true, "CopyHttpImage": true},
	"Attach":   {"UploadAttach": true, "DeleteAttach": true},
}

// Legacy Web actions are retained for API compatibility, but browser state
// changes must be same-origin POST requests. /api controllers are separate and
// continue to use their token contract.
func WebMutationInterceptor(c *revel.Controller) revel.Result {
	methods := webMutationMethods[strings.Title(c.Name)]
	if methods == nil || !methods[strings.Title(c.MethodName)] {
		return nil
	}
	if c.Request.Method != "POST" {
		c.Response.Status = 405
		return c.RenderJSON(info.Re{Ok: false, Msg: "methodNotAllowed"})
	}
	if !isAPI2Request(c) && c.Request.Header.Get("X-Requested-With") != "XMLHttpRequest" {
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
	sessionService = service.SessionS
	themeService = service.ThemeS
}

// 初始化博客模板
// 博客模板不由revel的
func initBlogTemplate() {
}

func init() {
	// interceptor
	// revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Index{}) // Index.Note自己校验
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Notebook{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Note{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Share{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &User{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Album{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &File{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Attach{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Tag{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &NoteContentHistory{})
	revel.InterceptFunc(WebMutationInterceptor, revel.BEFORE, &Notebook{})
	revel.InterceptFunc(WebMutationInterceptor, revel.BEFORE, &Note{})
	revel.InterceptFunc(WebMutationInterceptor, revel.BEFORE, &Share{})
	revel.InterceptFunc(WebMutationInterceptor, revel.BEFORE, &User{})
	revel.InterceptFunc(WebMutationInterceptor, revel.BEFORE, &File{})
	revel.InterceptFunc(WebMutationInterceptor, revel.BEFORE, &Attach{})
	revel.InterceptFunc(WebMutationInterceptor, revel.BEFORE, &Tag{})
	revel.InterceptFunc(WebMutationInterceptor, revel.BEFORE, &Album{})
	revel.InterceptFunc(WebMutationInterceptor, revel.BEFORE, &Web{})
	revel.InterceptMethod(Api2.RequireSession, revel.BEFORE)

}
