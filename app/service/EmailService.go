package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/info"
	. "github.com/gemsnote/gemsnote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 发送邮件

type EmailService struct {
	tpls map[string]*template.Template
	mu   sync.Mutex
}

func NewEmailService() *EmailService {
	return &EmailService{tpls: map[string]*template.Template{}}
}

// SendEmail reports queue acceptance; SMTP delivery runs exclusively in Celery.
func (this *EmailService) SendEmail(to, subject, body string) (bool, string) {
	_, err := this.QueueEmail("", to, subject, body)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}
func (this *EmailService) RegisterSendActiveEmail(userInfo info.User, email string) bool {
	_, err := this.QueueVerification(userInfo, email, info.TokenActiveEmail)
	return err == nil
}
func (this *EmailService) UpdateEmailSendActiveEmail(userInfo info.User, email string) (bool, string) {
	_, err := this.QueueVerification(userInfo, email, info.TokenUpdateEmail)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

// Reserve before creating a token: repeated submissions must not invalidate an
// already queued verification message. The cooldown is shared across workers.
func (this *EmailService) QueueVerification(user info.User, email string, kind int) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if !IsEmail(email) {
		return "", errors.New("invalidEmail")
	}
	if kind == info.TokenUpdateEmail {
		if strings.EqualFold(user.Email, email) {
			return "", errors.New("emailUnchanged")
		}
		if userService.IsExistsUser(email) {
			return "", errors.New("该邮箱已注册")
		}
	}
	site := strings.TrimRight(configService.GetSiteUrl(), "/")
	parsed, err := url.Parse(site)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", errors.New("mailSiteURLInvalid")
	}
	c, err := mailConnect()
	if err != nil {
		return "", err
	}
	defer c.Close()
	key := "gemsnote:mail:verification:" + user.UserId.Hex()
	lease := bson.NewObjectId().Hex()
	reserved, err := redis.String(c.Do("SET", key, lease, "NX", "EX", 60))
	if err == redis.ErrNil {
		return "", errors.New("mailRateLimited")
	}
	if err != nil || reserved != "OK" {
		return "", errors.New("mailQueueUnavailable")
	}
	accepted := false
	defer func() {
		if !accepted {
			redis.NewScript(1, `if redis.call('GET',KEYS[1])==ARGV[1] then return redis.call('DEL',KEYS[1]) end return 0`).Do(c, key, lease)
		}
	}()
	prefix, path := "emailTemplateRegister", "/user/activeEmail"
	if kind == info.TokenUpdateEmail {
		prefix, path = "emailTemplateUpdateEmail", "/user/updateEmail"
	}
	subject, tpl := configService.GetGlobalStringConfig(prefix+"Subject"), configService.GetGlobalStringConfig(prefix)
	if subject == "" || tpl == "" {
		return "", errors.New("mailTemplateInvalid")
	}
	token := tokenService.NewToken(user.UserId.Hex(), email, kind)
	if token == "" {
		return "", errors.New("mailTokenFailed")
	}
	values := map[string]interface{}{"siteUrl": site, "tokenUrl": site + path + "?token=" + url.QueryEscape(token), "token": token, "tokenTimeout": strconv.Itoa(int(tokenService.GetOverHours(kind))), "newEmail": email, "user": map[string]interface{}{"userId": user.UserId.Hex(), "email": user.Email, "username": user.Username}}
	ok, _, subject, body := this.renderEmail(subject, tpl, values)
	if !ok {
		return "", errors.New("mailTemplateInvalid")
	}
	id, err := this.QueueEmail(user.UserId.Hex(), email, subject, body)
	accepted = err == nil
	return id, err
}

func (this *EmailService) FindPwdSendEmail(token, email string) (ok bool, msg string) {
	subject := configService.GetGlobalStringConfig("emailTemplateFindPasswordSubject")
	tpl := configService.GetGlobalStringConfig("emailTemplateFindPassword")

	// 发送邮件
	tokenUrl := configService.GetSiteUrl() + "/findPassword/" + token
	// {siteUrl} {tokenUrl} {token} {tokenTimeout} {user.id} {user.email} {user.username}
	token2Value := map[string]interface{}{"siteUrl": configService.GetSiteUrl(), "tokenUrl": tokenUrl,
		"token": token, "tokenTimeout": strconv.Itoa(int(tokenService.GetOverHours(info.TokenActiveEmail)))}

	ok, msg, subject, tpl = this.renderEmail(subject, tpl, token2Value)
	if !ok {
		return
	}
	// 发送邮件
	ok, msg = this.SendEmail(email, subject, tpl)
	return
}

// 发送评论
func (this *EmailService) SendCommentEmail(note info.Note, comment info.BlogComment, userId, content string) bool {
	subject := configService.GetGlobalStringConfig("emailTemplateCommentSubject")
	tpl := configService.GetGlobalStringConfig("emailTemplateComment")

	// title := "评论提醒"

	/*
		toUserId := note.UserId.Hex()
		// title := "评论提醒"

		// 表示回复回复的内容, 那么发送给之前回复的
		if comment.CommentId != "" {
			toUserId = comment.UserId.Hex()
		}
		toUserInfo := userService.GetUserInfo(toUserId)
		sendUserInfo := userService.GetUserInfo(userId)

		subject := note.Title + " 收到 " + sendUserInfo.Username + " 的评论";
		if comment.CommentId != "" {
			subject = "您在 " + note.Title + " 发表的评论收到 " + sendUserInfo.Username;
			if userId == note.UserId.Hex() {
				subject += "(作者)";
			}
			subject += " 的评论";
		}
	*/

	toUserId := note.UserId.Hex()
	// 表示回复回复的内容, 那么发送给之前回复的
	if comment.CommentId != "" {
		toUserId = comment.UserId.Hex()
	}
	toUserInfo := userService.GetUserInfo(toUserId) // 被评论者
	sendUserInfo := userService.GetUserInfo(userId) // 评论者

	// {siteUrl} {blogUrl}
	// {blog.id} {blog.title} {blog.url}
	// {commentUser.userId} {commentUser.username} {commentUser.email}
	// {commentedUser.userId} {commentedUser.username} {commentedUser.email}
	token2Value := map[string]interface{}{"siteUrl": configService.GetSiteUrl(), "blogUrl": configService.GetBlogUrl(),
		"blog": map[string]string{
			"id":    note.NoteId.Hex(),
			"title": note.Title,
			"url":   configService.GetBlogUrl() + "/view/" + note.NoteId.Hex(),
		},
		"commentContent": content,
		// 评论者信息
		"commentUser": map[string]interface{}{"userId": sendUserInfo.UserId.Hex(),
			"username":     sendUserInfo.Username,
			"email":        sendUserInfo.Email,
			"isBlogAuthor": userId == note.UserId.Hex(),
		},
		// 被评论者信息
		"commentedUser": map[string]interface{}{"userId": toUserId,
			"username":     toUserInfo.Username,
			"email":        toUserInfo.Email,
			"isBlogAuthor": toUserId == note.UserId.Hex(),
		},
	}

	ok := false
	ok, _, subject, tpl = this.renderEmail(subject, tpl, token2Value)
	if !ok {
		return false
	}

	// 发送邮件
	ok, _ = this.SendEmail(toUserInfo.Email, subject, tpl)
	return ok
}

// 验证模板是否正确
func (this *EmailService) ValidTpl(str string) (ok bool, msg string) {
	defer func() {
		if err := recover(); err != nil {
			ok = false
			msg = fmt.Sprint(err)
		}
	}()
	header := configService.GetGlobalStringConfig("emailTemplateHeader")
	footer := configService.GetGlobalStringConfig("emailTemplateFooter")
	str = strings.Replace(str, "{{header}}", header, -1)
	str = strings.Replace(str, "{{footer}}", footer, -1)
	_, err := template.New("tpl name").Parse(str)
	if err != nil {
		msg = fmt.Sprint(err)
		return
	}
	ok = true
	return
}

// ok, msg, subject, tpl
func (this *EmailService) getTpl(str string) (ok bool, msg string, tpl *template.Template) {
	this.mu.Lock()
	defer this.mu.Unlock()
	defer func() {
		if err := recover(); err != nil {
			ok = false
			msg = fmt.Sprint(err)
		}
	}()

	var err error
	var has bool

	if tpl, has = this.tpls[str]; !has {
		tpl, err = template.New("tpl name").Parse(str)
		if err != nil {
			msg = fmt.Sprint(err)
			return
		}
		this.tpls[str] = tpl
	}
	ok = true
	return
}

// 通过subject, body和值得到内容
func (this *EmailService) renderEmail(subject, body string, values map[string]interface{}) (ok bool, msg string, o string, b string) {
	ok = false
	msg = ""
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			ok = false
			msg = fmt.Sprint(err) // 这里的err其实就是panic传入的内容，
		}
	}()

	var tpl *template.Template

	values["siteUrl"] = configService.GetSiteUrl()

	// subject
	if subject != "" {
		ok, msg, tpl = this.getTpl(subject)
		if !ok {
			return
		}
		var buffer bytes.Buffer
		err := tpl.Execute(&buffer, values)
		if err != nil {
			msg = fmt.Sprint(err)
			return
		}
		o = buffer.String()
	} else {
		o = ""
	}

	// content
	header := configService.GetGlobalStringConfig("emailTemplateHeader")
	footer := configService.GetGlobalStringConfig("emailTemplateFooter")
	body = strings.Replace(body, "{{header}}", header, -1)
	body = strings.Replace(body, "{{footer}}", footer, -1)
	values["subject"] = o
	ok, msg, tpl = this.getTpl(body)
	if !ok {
		return
	}
	var buffer2 bytes.Buffer
	err := tpl.Execute(&buffer2, values)
	if err != nil {
		msg = fmt.Sprint(err)
		return
	}
	b = buffer2.String()

	return
}

// Bulk methods synchronously report queue acceptance; Celery performs delivery.
// Never leave queue failures hidden inside a detached goroutine.
func (this *EmailService) SendEmailToUsers(users []info.User, subject, body string) (bool, string) {
	if len(users) == 0 {
		return false, "no users"
	}
	for _, user := range users {
		ok, msg, renderedSubject, renderedBody := this.renderEmail(subject, body, map[string]interface{}{"userId": user.UserId.Hex(), "username": user.Username, "email": user.Email})
		if !ok {
			return false, msg
		}
		if ok, msg = this.queueLoggedEmail(user.Email, renderedSubject, renderedBody); !ok {
			return false, msg
		}
	}
	return true, ""
}
func (this *EmailService) SendEmailToEmails(emails []string, subject, body string) (bool, string) {
	if len(emails) == 0 {
		return false, "no emails"
	}
	for _, email := range emails {
		if email == "" {
			continue
		}
		ok, msg, renderedSubject, renderedBody := this.renderEmail(subject, body, map[string]interface{}{"email": email})
		if !ok {
			return false, msg
		}
		if ok, msg = this.queueLoggedEmail(email, renderedSubject, renderedBody); !ok {
			return false, msg
		}
	}
	return true, ""
}
func (this *EmailService) queueLoggedEmail(to, subject, body string) (bool, string) {
	id, err := this.QueueEmail("", to, subject, body)
	if err != nil {
		this.AddEmailLog(to, subject, body, false, err.Error())
		return false, err.Error()
	}
	this.AddEmailLog(to, subject, body, true, "queued:"+id)
	return true, ""
}

// 添加邮件日志
func (this *EmailService) AddEmailLog(email, subject, body string, ok bool, msg string) {
	log := info.EmailLog{LogId: bson.NewObjectId(), Email: email, Subject: subject, Body: body, Ok: ok, Msg: msg, CreatedTime: time.Now()}
	db.Insert(db.EmailLogs, log)
}

// 展示邮件日志

func (this *EmailService) DeleteEmails(ids []string) bool {
	idsO := make([]bson.ObjectId, len(ids))
	for i, id := range ids {
		idsO[i] = bson.ObjectIdHex(id)
	}
	db.DeleteAll(db.EmailLogs, bson.M{"_id": bson.M{"$in": idsO}})

	return true
}
func (this *EmailService) ListEmailLogs(pageNumber, pageSize int, sortField string, isAsc bool, email string) (page info.Page, emailLogs []info.EmailLog) {
	emailLogs = []info.EmailLog{}
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)
	query := bson.M{}
	if email != "" {
		query["Email"] = bson.M{"$regex": bson.RegEx{".*?" + email + ".*", "i"}}
	}
	count := db.Count(db.EmailLogs, query)
	db.ListByQOptions(db.EmailLogs, query, &emailLogs, db.QueryOptions{
		Sort: []string{sortFieldR}, Skip: skipNum, Limit: pageSize,
	})
	page = info.NewPage(pageNumber, pageSize, count, nil)
	return
}
