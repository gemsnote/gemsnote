package api

import (
	"github.com/gemsnote/gemsnote/app/info"
	"github.com/gemsnote/gemsnote/app/service"
	"github.com/revel/revel"
	//	"strconv"
)

// 用户登录后生成一个token, 将这个token保存到session中
// 以后每次的请求必须带这个token, 并从session中获取userId

// 用户登录/注销/找回密码

type ApiAuth struct {
	ApiBaseContrller
}

// 登录
// [ok]
// 成功返回 {Ok: true, Item: token }
// 失败返回 {Ok: false, Msg: ""}
func (c ApiAuth) Login(email, pwd string) revel.Result {
	if email == "" && len(c.Params.JSON) > 0 {
		var payload struct {
			Email string `json:"email"`
			Pwd   string `json:"pwd"`
		}
		if err := c.Params.BindJSON(&payload); err != nil {
			return c.RenderJSON(info.ApiRe{Ok: false, Msg: "invalidJSON"})
		}
		email, pwd = payload.Email, payload.Pwd
	}
	var msg = ""

	userInfo, err := authService.Login(email, pwd)
	if err == nil {
		token, tokenErr := service.NewAPIToken()
		if tokenErr != nil {
			revel.AppLog.Error("API token generation failed", "error", tokenErr)
			c.Response.Status = 500
			return c.RenderJSON(info.ApiRe{Ok: false, Msg: "tokenGenerationFailed"})
		}
		sessionService.SetUserId(token, userInfo.UserId.Hex())
		return c.RenderJSON(info.AuthOk{Ok: true, Token: token, UserId: userInfo.UserId, Email: userInfo.Email, Username: userInfo.Username})
	} else {
		// 登录错误, 则错误次数++
		msg = "wrongUsernameOrPassword"
	}
	return c.RenderJSON(info.ApiRe{Ok: false, Msg: c.Message(msg)})
}

// 注销
// [Ok]
func (c ApiAuth) Logout() revel.Result {
	token := c.getToken()
	sessionService.Clear(token)
	re := info.NewApiRe()
	re.Ok = true
	return c.RenderJSON(re)
}
