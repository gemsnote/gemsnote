package controllers

// API2 is the versioned browser/client API facade.  It deliberately reuses
// the existing service-backed Web actions so that the first migration phase
// cannot introduce a second implementation of note permissions or conflict
// handling.  The facade accepts JSON bodies for core note operations and lives
// under /api2; the legacy /api routes are untouched.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/gemsnote/gemsnote/app/info"
	appversion "github.com/gemsnote/gemsnote/app/version"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

type Api2 struct{ BaseController }

// TokenLogin is the API2 desktop login contract. Authentication, public user
// profile and protocol version are intentionally returned in one round trip.
func (c Api2) TokenLogin() revel.Result {
	var p struct{ Email, Pwd string }
	if len(c.Params.JSON) > 0 {
		if err := c.body(&p); err != nil {
			return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
		}
	} else {
		p.Email = c.Params.Form.Get("email")
		p.Pwd = c.Params.Form.Get("pwd")
	}
	user, err := authService.Login(p.Email, p.Pwd)
	if err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: c.Message("wrongUsernameOrPassword")})
	}
	token := bson.NewObjectId().Hex()
	sessionService.SetUserId(token, user.UserId.Hex())
	return c.RenderJSON(map[string]interface{}{
		"Ok":     true,
		"Token":  token,
		"User":   info.ApiUser{UserId: user.UserId.Hex(), Username: user.Username, Email: user.Email, Verified: user.Verified, Logo: user.Logo},
		"Server": map[string]string{"Name": "gemsnote", "Version": appversion.Current, "MinVersion": ""},
	})
}

func (c Api2) Login() revel.Result {
	var p struct{ Email, Pwd, Captcha string }
	if len(c.Params.JSON) > 0 {
		if err := c.body(&p); err != nil {
			return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
		}
	} else {
		// Desktop bridge and older clients submit the login form as
		// application/x-www-form-urlencoded. Keep that transport supported
		// while JSON remains the preferred API2 format.
		p.Email = c.Params.Form.Get("email")
		p.Pwd = c.Params.Form.Get("pwd")
		p.Captcha = c.Params.Form.Get("captcha")
	}
	return (Auth{c.BaseController}).DoLogin(p.Email, p.Pwd, p.Captcha)
}

func (c Api2) Logout() revel.Result {
	sessionID := c.Session.ID()
	sessionService.Clear(sessionID)
	c.ClearSession()
	return c.RenderJSON(info.Re{Ok: true})
}

func (c Api2) Register() revel.Result {
	var p struct {
		Email string `json:"email"`
		Pwd   string `json:"pwd"`
		IU    string `json:"iu"`
	}
	if err := c.body(&p); err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
	}
	return (Auth{c.BaseController}).DoRegister(p.Email, p.Pwd, p.IU)
}

func (c Api2) RequestPasswordReset() revel.Result {
	var p struct {
		Email string `json:"email"`
	}
	if err := c.body(&p); err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
	}
	return (Auth{c.BaseController}).DoFindPassword(p.Email)
}

func (c Api2) ResetPassword() revel.Result {
	var p struct {
		Token string `json:"token"`
		Pwd   string `json:"pwd"`
	}
	if err := c.body(&p); err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
	}
	return (Auth{c.BaseController}).FindPasswordUpdate(p.Token, p.Pwd)
}

func (c Api2) body(v interface{}) error {
	// Revel parses application/json before invoking the action and stores the
	// bytes in Params.JSON. Reading Request.GetBody here would return an empty
	// reader, which caused every API2 write to report invalidJSON.
	if len(c.Params.JSON) > 0 {
		return c.Params.BindJSON(v)
	}
	reader := c.Request.GetBody()
	if reader == nil {
		return io.EOF
	}
	b, err := io.ReadAll(reader)
	if err == nil {
		return json.Unmarshal(b, v)
	}
	return err
}

func (c Api2) Bootstrap() revel.Result { return (Web{c.BaseController}).Bootstrap() }

func (c Api2) Groups() revel.Result { return (Web{c.BaseController}).Groups() }

func (c Api2) Attachments() revel.Result {
	var p struct {
		NoteID string `json:"noteId"`
	}
	if err := c.body(&p); err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
	}
	re := info.NewRe()
	re.Ok = true
	re.List = attachService.ListAttachs(p.NoteID, c.GetUserId())
	return c.RenderJSON(re)
}

func (c Api2) UploadAttachment() revel.Result {
	return (Attach{c.BaseController}).UploadAttach(c.Params.Form.Get("noteId"))
}

func (c Api2) UploadAvatar() revel.Result {
	fileController := File{c.BaseController}
	re := fileController.uploadImage("logo", "")
	if re.Ok {
		re.Ok = userService.UpdateAvatar(c.GetUserId(), re.Id)
		if re.Ok {
			c.UpdateSession("Logo", re.Id)
		} else {
			re.Msg = "avatarUpdateFailed"
		}
	}
	if !re.Ok {
		revel.AppLog.Warn("Avatar upload failed", "userId", c.GetUserId(), "message", re.Msg)
	}
	return c.RenderJSON(re)
}

func (c Api2) DeleteAttachment() revel.Result {
	var p struct {
		AttachID string `json:"attachId"`
	}
	if err := c.body(&p); err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
	}
	return (Attach{c.BaseController}).DeleteAttach(p.AttachID)
}

func (c Api2) DownloadAttachment(attachId string) revel.Result {
	return (Attach{c.BaseController}).Download(attachId)
}

func (c Api2) AdminData() revel.Result {
	var p struct {
		Keywords string `json:"keywords"`
		Page     int    `json:"page"`
	}
	if err := c.body(&p); err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
	}
	return (Web{c.BaseController}).AdminData(p.Keywords, p.Page)
}

func (c Api2) Notes() revel.Result {
	var p struct {
		NotebookID string `json:"notebookId"`
		Key        string `json:"key"`
		Tag        string `json:"tag"`
		Sort       string `json:"sort"`
		Trash      bool   `json:"trash"`
		Starred    bool   `json:"starred"`
		Page       int    `json:"page"`
	}
	_ = c.body(&p)
	if c.Params.Form == nil {
		c.Params.Form = url.Values{}
	}
	c.Params.Form.Set("starred", boolString(p.Starred))
	if p.Page > 0 {
		c.Params.Form.Set("page", intString(p.Page))
	}
	return (Web{c.BaseController}).Notes(p.NotebookID, p.Key, p.Tag, p.Sort, p.Trash)
}

func (c Api2) Star() revel.Result {
	var p struct {
		NoteID  string `json:"noteId"`
		Starred bool   `json:"starred"`
	}
	if err := c.body(&p); err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
	}
	return (Web{c.BaseController}).Star(p.NoteID, p.Starred)
}

func (c Api2) Document() revel.Result {
	var p struct {
		NoteID string `json:"noteId"`
	}
	if err := c.body(&p); err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
	}
	return (Web{c.BaseController}).Document(p.NoteID)
}

func (c Api2) Save() revel.Result {
	var p struct {
		NoteID, NotebookID, OwnerID, Title, Content, Tags string
		IsNew                                             bool `json:"isNew"`
		IsMarkdown                                        bool `json:"isMarkdown"`
		Usn                                               int  `json:"usn"`
	}
	if err := c.body(&p); err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
	}
	return (Web{c.BaseController}).Save(p.NoteID, p.NotebookID, p.OwnerID, p.Title, p.Content, p.Tags, p.IsNew, p.IsMarkdown, p.Usn)
}

func (c Api2) Restore() revel.Result {
	var p struct {
		NoteID string `json:"noteId"`
	}
	if err := c.body(&p); err != nil {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
	}
	return (Web{c.BaseController}).Restore(p.NoteID)
}

// Histories always returns a JSON array. The legacy controller may return
// null when a note has no history; the API2 browser client expects [].
func (c Api2) Histories(noteId string) revel.Result {
	if noteId == "" && len(c.Params.JSON) > 0 {
		var payload struct {
			NoteID string `json:"noteId"`
		}
		if err := c.body(&payload); err != nil {
			return c.RenderJSON(info.Re{Ok: false, Msg: "invalidJSON"})
		}
		noteId = payload.NoteID
	}
	if !bson.IsObjectIdHex(noteId) {
		return c.RenderJSON(info.Re{Ok: false, Msg: "invalidNote"})
	}
	return c.RenderJSON(api2HistoryArray(noteContentHistoryService.ListHistories(noteId, c.GetUserId())))
}

func api2HistoryArray(histories []info.EachHistory) []info.EachHistory {
	if histories == nil {
		return []info.EachHistory{}
	}
	return histories
}

func (c Api2) RequireSession() revel.Result {
	// Login must remain public; Logout is idempotent and may be called after a
	// session has already expired. Bootstrap is also public so the login page
	// can discover registration/captcha state.
	if c.MethodName == "Bootstrap" || c.MethodName == "Login" || c.MethodName == "TokenLogin" || c.MethodName == "Logout" || c.MethodName == "Register" || c.MethodName == "RequestPasswordReset" || c.MethodName == "ResetPassword" {
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

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
func intString(v int) string {
	if v <= 0 {
		return ""
	}
	return fmt.Sprintf("%d", v)
}
