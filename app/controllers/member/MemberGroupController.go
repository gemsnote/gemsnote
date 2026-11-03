package member

import (
	"github.com/gemsnote/gemsnote/app/info"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

// 分组管理
type MemberGroup struct {
	MemberBaseController
}

// 首页, 显示所有分组和用户
func (c MemberGroup) Index() revel.Result {
	c.SetUserInfo()
	c.SetLocale()
	c.ViewArgs["title"] = c.Message("My Group")
	c.ViewArgs["groups"] = groupService.GetGroupsAndUsers(c.GetUserId())
	return c.RenderTemplate("member/group/index.html")
}

// 添加分组
func (c MemberGroup) AddGroup(title string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Item = groupService.AddGroup(c.GetUserId(), title)
	return c.RenderJSON(re)
}

func (c MemberGroup) UpdateGroupTitle(groupId, title string) revel.Result {
	re := info.NewRe()
	re.Ok = groupService.UpdateGroupTitle(c.GetUserId(), groupId, title)
	return c.RenderJSON(re)
}

func (c MemberGroup) DeleteGroup(groupId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = groupService.DeleteGroup(c.GetUserId(), groupId)
	return c.renderGroupResult(re)
}

// 添加用户
func (c MemberGroup) AddUser(groupId, email string) revel.Result {
	re := info.NewRe()
	// Keep the wire field name for existing clients. API2 treats it as an
	// identity, never as an ObjectId (a username can itself be 24 hex digits).
	var userInfo info.User
	if c.Request != nil && strings.HasPrefix(c.Request.URL.Path, "/api2/") {
		if !bson.IsObjectIdHex(groupId) || !bson.IsObjectIdHex(c.GetUserId()) {
			re.Msg = "notMyGroup"
			return c.renderGroupResult(re)
		}
		identity := strings.TrimSpace(email)
		if identity == "" {
			re.Msg = "userNotExists"
			return c.renderGroupResult(re)
		}
		userInfo = userService.GetUserInfoByName(identity)
	} else {
		userInfo = userService.GetUserInfoByAny(email)
	}
	if userInfo.UserId == "" {
		re.Msg = "userNotExists"
	} else {
		re.Ok, re.Msg = groupService.AddUser(c.GetUserId(), groupId, userInfo.UserId.Hex())
		re.Item = userInfo
	}
	return c.renderGroupResult(re)
}

func (c MemberGroup) DeleteUser(groupId, userId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = groupService.DeleteUser(c.GetUserId(), groupId, userId)
	return c.renderGroupResult(re)
}

// API2 clients choose their own locale. Preserve legacy server-side messages
// only on the historical routes.
func (c MemberGroup) renderGroupResult(re info.Re) revel.Result {
	if c.Request != nil && strings.HasPrefix(c.Request.URL.Path, "/api2/") {
		return c.RenderJSON(re)
	}
	return c.RenderRe(re)
}
