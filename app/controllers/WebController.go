package controllers

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gemsnote/gemsnote/app/info"
	"github.com/gemsnote/gemsnote/app/lea"
	"github.com/gemsnote/gemsnote/app/service"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// Web serves the independent Vue application and its session-based data APIs.
type Web struct{ BaseController }

func (c Web) Page() revel.Result {
	root := revel.Config.StringDefault("frontend.dist", "frontend/dist")
	if !filepath.IsAbs(root) {
		root = filepath.Join(revel.BasePath, root)
	}
	file, err := os.Open(filepath.Join(root, "index.html"))
	if err != nil {
		c.Response.Status = 503
		return c.RenderText("前端尚未构建，请执行 npm ci --prefix frontend && npm run build --prefix frontend")
	}
	c.Response.Out.Header().Set("Cache-Control", "no-cache")
	return c.RenderFile(file, revel.Inline)
}

func (c Web) Bootstrap() revel.Result {
	userId := c.GetUserId()
	if userId == "" || !sessionService.ValidateUserSession(c.Session.ID(), userId) {
		return c.RenderJSON(map[string]interface{}{"Ok": true, "User": nil, "OpenRegister": configService.IsOpenRegister(), "NeedCaptcha": sessionService.LoginTimesIsOver(c.Session.ID())})
	}
	notebooks, users := shareService.GetShareNotebooks(userId)
	// A user may share individual notes without sharing a notebook. Expose a
	// synthetic default entry so the Vue navigation can reach those notes too.
	if notebooks == nil {
		notebooks = info.ShareNotebooksByUser{}
	}
	for _, owner := range users {
		key := owner.UserId.Hex()
		if _, ok := notebooks[key]; !ok {
			notebooks[key] = []info.ShareNotebooks{{IsDefault: true, Notebook: info.Notebook{Title: owner.Username}}}
		}
	}
	allTags := tagService.GetTags(userId)
	tags := make([]info.NoteTag, 0, len(allTags))
	for _, tag := range allTags {
		tag.Count = noteService.CountNoteByTag(userId, tag.Tag)
		if tag.Count > 0 {
			tags = append(tags, tag)
		}
	}
	starredCount, _ := noteService.ListStarNotes(userId, 1, 1, "UpdatedTime", false)
	return c.RenderJSON(map[string]interface{}{"Ok": true, "User": c.GetUserInfo(), "IsAdmin": c.GetUsername() == configService.GetAdminUsername(), "Notebooks": notebookService.GetNotebooks(userId), "SharedNotebooks": notebooks, "SharedUsers": users, "Tags": tags, "TotalNotes": noteService.CountNote(userId), "TotalStarred": starredCount, "Version": configService.GetVersion()})
}

func (c Web) RequireSession() revel.Result {
	if c.MethodName == "Page" || c.MethodName == "Bootstrap" || c.MethodName == "VerifyEmail" {
		return nil
	}
	if c.GetUserId() == "" || !sessionService.ValidateUserSession(c.Session.ID(), c.GetUserId()) {
		c.Response.Status = 401
		return c.RenderJSON(info.Re{Ok: false, Msg: "NOTLOGIN"})
	}
	if c.Request.Method != "GET" && c.Request.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		c.Response.Status = 403
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidRequest"})
	}
	return nil
}

func (c Web) Notes(notebookId, key, tag, sort string, trash bool) revel.Result {
	if notebookId != "" && !bson.IsObjectIdHex(notebookId) {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidNotebook"})
	}
	if sort != "Title" && sort != "CreatedTime" && sort != "UpdatedTime" {
		sort = defaultSortField
	}
	var notes []info.Note
	starred := c.Params.Form.Get("starred") == "true"
	if starred {
		_, notes = noteService.ListStarNotes(c.GetUserId(), c.GetPage(), 100, sort, false)
	} else if key != "" {
		_, notes = noteService.SearchNote(key, c.GetUserId(), c.GetPage(), 100, sort, false, false)
	} else if tag != "" {
		_, notes = noteService.SearchNoteByTags([]string{tag}, c.GetUserId(), c.GetPage(), 100, sort, false)
	} else {
		_, notes = noteService.ListNotes(c.GetUserId(), notebookId, trash, c.GetPage(), 100, sort, false, false)
	}
	return c.RenderJSON(notes)
}

func (c Web) Star(noteId string, starred bool) revel.Result {
	if !bson.IsObjectIdHex(noteId) {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidNote"})
	}
	note := noteService.GetNoteById(noteId)
	if note.NoteId == "" || note.IsDeleted || note.IsTrash || note.UserId.Hex() != c.GetUserId() {
		return c.RenderJSON(info.Re{Ok: false, Msg: "noAuth"})
	}
	return c.RenderJSON(info.Re{Ok: noteService.SetStar(c.GetUserId(), noteId, starred)})
}

func (c Web) Document(noteId string) revel.Result {
	if !bson.IsObjectIdHex(noteId) {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidNote"})
	}
	note := noteService.GetNoteById(noteId)
	if note.NoteId == "" || note.IsDeleted || (note.UserId.Hex() != c.GetUserId() && !shareService.HasReadPerm(note.UserId.Hex(), c.GetUserId(), noteId)) {
		c.Response.Status = 403
		return c.RenderJSON(info.Re{Ok: false, Msg: "noAuth"})
	}
	return c.RenderJSON(map[string]interface{}{"Note": note, "Content": noteService.GetNoteContent(noteId, note.UserId.Hex()).Content, "Writable": !note.IsTrash && (note.UserId.Hex() == c.GetUserId() || shareService.HasUpdatePerm(note.UserId.Hex(), c.GetUserId(), noteId))})
}

var webSaveMutex sync.Mutex

func (c Web) Save(noteId, notebookId, ownerId, title, content, tags string, isNew, isMarkdown bool, usn int) revel.Result {
	if !bson.IsObjectIdHex(noteId) || (isNew && !bson.IsObjectIdHex(notebookId)) {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidId"})
	}
	webSaveMutex.Lock()
	defer webSaveMutex.Unlock()
	if isNew {
		if ownerId == "" {
			ownerId = c.GetUserId()
		}
		if !bson.IsObjectIdHex(ownerId) || notebookService.GetNotebook(notebookId, ownerId).NotebookId == "" || (ownerId != c.GetUserId() && !shareService.HasUpdateNotebookPerm(ownerId, c.GetUserId(), notebookId)) {
			return c.RenderJSON(info.Re{Ok: false, Msg: "invalidNotebook"})
		}
		if noteService.GetNoteById(noteId).NoteId != "" {
			return c.RenderJSON(info.Re{Ok: false, Msg: "conflict"})
		}
		note := info.Note{NoteId: bson.ObjectIdHex(noteId), NotebookId: bson.ObjectIdHex(notebookId), UserId: bson.ObjectIdHex(ownerId), Title: title, Tags: strings.Split(tags, ","), IsMarkdown: isMarkdown}
		created := noteService.AddNoteAndContentForController(note, info.NoteContent{NoteId: note.NoteId, UserId: note.UserId, Content: content}, c.GetUserId())
		if created.NoteId == "" {
			return c.RenderJSON(info.Re{Ok: false, Msg: "saveFailed"})
		}
	} else {
		note := noteService.GetNoteById(noteId)
		if note.NoteId == "" || note.IsDeleted || note.IsTrash {
			return c.RenderJSON(info.Re{Ok: false, Msg: "notExists"})
		}
		if note.Usn != usn {
			return c.RenderJSON(info.Re{Ok: false, Msg: "conflict"})
		}
		ok, msg, _ := noteService.UpdateNoteByUsn(c.GetUserId(), noteId, title, strings.Split(tags, ","), usn)
		if !ok {
			return c.RenderJSON(info.Re{Ok: false, Msg: msg})
		}
		ok, msg, _ = noteService.UpdateNoteContent(c.GetUserId(), noteId, content, "", true, -1, time.Now())
		if !ok {
			return c.RenderJSON(info.Re{Ok: false, Msg: msg})
		}
	}
	return c.Document(noteId)
}

func (c Web) Restore(noteId string) revel.Result {
	if !bson.IsObjectIdHex(noteId) {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidNote"})
	}
	note := noteService.GetNote(noteId, c.GetUserId())
	if note.NoteId == "" || !note.IsTrash {
		return c.RenderJSON(info.Re{Ok: false, Msg: "notExists"})
	}
	ok, msg, _ := noteService.UpdateNote(c.GetUserId(), noteId, bson.M{"IsTrash": false}, note.Usn)
	return c.RenderJSON(info.Re{Ok: ok, Msg: msg})
}

func (c Web) Groups() revel.Result {
	return c.RenderJSON(service.GroupS.GetGroupsAndUsers(c.GetUserId()))
}

func (c Web) ShareMembers(noteId string) revel.Result {
	if !bson.IsObjectIdHex(noteId) {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidNote"})
	}
	return c.RenderJSON(map[string]interface{}{"Users": shareService.ListNoteShareUserInfo(noteId, c.GetUserId())})
}

func (c Web) EmailChange(email, pwd string) revel.Result {
	if !lea.IsEmail(email) {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidEmail"})
	}
	user, err := authService.Login(c.GetEmail(), pwd)
	if err != nil || user.UserId.Hex() != c.GetUserId() {
		return c.RenderJSON(info.Re{Ok: false, Msg: "wrongPassword"})
	}
	ok, msg := emailService.UpdateEmailSendActiveEmail(user, email)
	return c.RenderJSON(info.Re{Ok: ok, Msg: msg})
}

func (c Web) VerifyEmail(token string, change bool) revel.Result {
	var ok bool
	var msg, email string
	if change {
		ok, msg, email = userService.UpdateEmail(token)
	} else {
		ok, msg, email = userService.ActiveEmail(token)
	}
	return c.RenderJSON(map[string]interface{}{"Ok": ok, "Msg": msg, "Email": email})
}

func (c Web) isAdmin() bool { return c.GetUsername() == configService.GetAdminUsername() }

var webSettingKeys = []string{
	"siteUrl", "openRegister", "emailHost", "emailPort", "emailUsername", "emailPassword", "emailSSL",
	"uploadImageSize", "uploadAvatarSize", "uploadAttachSize", "exportPdfBinPath",
}

func (c Web) AdminData(keywords string, page int) revel.Result {
	if !c.isAdmin() {
		c.Response.Status = 403
		return c.RenderJSON(info.Re{Ok: false, Msg: "noAuth"})
	}
	if page < 1 {
		page = 1
	}
	_, users := userService.ListUsers(page, 20, "CreatedTime", false, keywords)
	settings := map[string]string{}
	for _, key := range webSettingKeys {
		if key != "emailPassword" {
			settings[key] = configService.GetGlobalStringConfig(key)
		}
	}
	return c.RenderJSON(map[string]interface{}{"Users": users, "Settings": settings})
}

func (c Web) AdminSettings() revel.Result {
	if !c.isAdmin() {
		c.Response.Status = 403
		return c.RenderJSON(info.Re{Ok: false, Msg: "noAuth"})
	}
	for _, key := range webSettingKeys {
		if c.Has(key) {
			value := c.Params.Values.Get(key)
			if key == "emailPassword" && value == "" {
				continue
			}
			if !configService.UpdateGlobalStringConfig(c.GetUserId(), key, value) {
				return c.RenderJSON(info.Re{Ok: false, Msg: "updateFailed"})
			}
		}
	}
	return c.RenderJSON(info.Re{Ok: true})
}

func (c Web) AdminResetPwd(userId, pwd string) revel.Result {
	if !c.isAdmin() {
		c.Response.Status = 403
		return c.RenderJSON(info.Re{Ok: false, Msg: "noAuth"})
	}
	if ok, msg := lea.Vd("password", pwd); !ok {
		return c.RenderJSON(info.Re{Ok: false, Msg: msg})
	}
	ok, msg := userService.ResetPwd(c.GetUserId(), userId, pwd)
	return c.RenderJSON(info.Re{Ok: ok, Msg: msg})
}

func (c Web) AdminRegister(email, pwd string) revel.Result {
	if !c.isAdmin() {
		c.Response.Status = 403
		return c.RenderJSON(info.Re{Ok: false, Msg: "noAuth"})
	}
	if ok, msg := lea.Vd("email", email); !ok {
		return c.RenderJSON(info.Re{Ok: false, Msg: msg})
	}
	if ok, msg := lea.Vd("password", pwd); !ok {
		return c.RenderJSON(info.Re{Ok: false, Msg: msg})
	}
	ok, msg := authService.Register(email, pwd, "")
	return c.RenderJSON(info.Re{Ok: ok, Msg: msg})
}

func init() { revel.InterceptMethod(Web.RequireSession, revel.BEFORE) }
