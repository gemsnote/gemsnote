package service

import (
	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/info"
	. "github.com/gemsnote/gemsnote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

// 找回密码
// 修改密码

type PwdService struct {
}

// 1. 找回密码, 通过email找用户,
// 用户存在, 生成code
func (this *PwdService) FindPwd(email string) (ok bool, msg string) {
	ok = false
	email = strings.ToLower(strings.TrimSpace(email))
	userId := userService.GetUserId(email)
	if userId == "" {
		msg = "用户不存在"
		return
	}

	token := tokenService.NewToken(userId, email, info.TokenPwd)
	if token == "" {
		return false, "db error"
	}

	// 发送邮件
	ok, msg = emailService.FindPwdSendEmail(token, email)
	return
}

// 重置密码时
// 修改密码
// 先验证
func (this *PwdService) UpdatePwd(token, pwd string) (bool, string) {
	ok, msg, tokenInfo := tokenService.VerifyLinkToken(token, info.TokenPwd, info.TokenInitialPassword)
	if !ok {
		return false, msg
	}
	user := userService.GetUserInfo(tokenInfo.UserId.Hex())
	if user.UserId == "" || user.Email != tokenInfo.Email {
		return false, "verificationLinkInvalid"
	}
	query := bson.M{"_id": tokenInfo.UserId, "Email": user.Email, "Pwd": user.Pwd}
	if tokenInfo.Type == info.TokenInitialPassword {
		if user.Pwd != "" || !user.Verified {
			return false, "verificationLinkInvalid"
		}
		query["Pwd"] = ""
	}

	passwd := GenPwd(pwd)
	if passwd == "" {
		return false, "GenerateHash error"
	}

	// 修改密码之
	ok = db.UpdateByQMap(db.Users, query, bson.M{"Pwd": passwd, "Verified": true})
	if !ok {
		return false, ""
	}

	// A recovered password is a credential change too: revoke every existing
	// web session and API token before consuming the recovery token.
	sessionService.ClearUserSessions(tokenInfo.UserId.Hex())

	// 删除token
	db.Delete(db.Tokens, bson.M{"_id": tokenInfo.UserId, "Token": token, "Type": tokenInfo.Type})

	return true, ""
}
