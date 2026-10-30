package service

import (
	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/info"
	. "github.com/gemsnote/gemsnote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// token
// 找回密码
// 修改密码

type TokenService struct {
}

// 生成token
func (this *TokenService) NewToken(userId string, email string, tokenType int) string {
	token := info.Token{UserId: bson.ObjectIdHex(userId), Token: NewGuidWith(email), CreatedTime: time.Now(), Email: email, Type: tokenType}

	if db.Upsert(db.Tokens, bson.M{"_id": token.UserId}, token) {
		return token.Token
	}

	return ""
}

// 删除token
func (this *TokenService) DeleteToken(userId string, tokenType int) bool {
	return db.Delete(db.Tokens, bson.M{"_id": bson.ObjectIdHex(userId), "Type": tokenType})
}

func (this *TokenService) GetOverHours(tokenType int) float64 {
	if tokenType == info.TokenPwd || tokenType == info.TokenInitialPassword {
		return info.PwdOverHours
	} else if tokenType == info.TokenUpdateEmail {
		return info.UpdateEmailOverHours
	} else {
		return info.ActiveEmailOverHours
	}
}

// 验证token, 是否存在, 过时?
func (this *TokenService) VerifyToken(token string, tokenType int) (ok bool, msg string, tokenInfo info.Token) {
	overHours := this.GetOverHours(tokenType)

	ok = false
	if token == "" {
		msg = "不存在"
		return
	}

	db.GetByQ(db.Tokens, bson.M{"Token": token, "Type": tokenType}, &tokenInfo)

	if tokenInfo.UserId == "" || tokenInfo.Type != tokenType || !db.UserAvailable(tokenInfo.UserId.Hex()) {
		msg = "不存在"
		return
	}

	// 验证是否过时
	now := time.Now()
	duration := now.Sub(tokenInfo.CreatedTime)

	if duration < 0 || duration.Hours() > overHours {
		msg = "过期"
		return
	}

	ok = true
	return
}

// VerifyLinkToken distinguishes expiry from an unknown/replaced link. Accepting
// multiple purposes in one lookup preserves the cause for setup/reset links.
func (this *TokenService) VerifyLinkToken(token string, types ...int) (bool, string, info.Token) {
	if token == "" {
		return false, "verificationLinkMissing", info.Token{}
	}
	var record info.Token
	db.GetByQ(db.Tokens, bson.M{"Token": token}, &record)
	var user info.User
	if record.UserId.Valid() {
		db.GetByQWithFields(db.Users, bson.M{"_id": record.UserId}, []string{"_id", "Flag"}, &user)
	}
	msg := validateLinkRecord(token, record, user, types, time.Now())
	return msg == "", msg, record
}

func validateLinkRecord(token string, record info.Token, user info.User, types []int, now time.Time) string {
	if token == "" {
		return "verificationLinkMissing"
	}
	allowed := false
	for _, kind := range types {
		if record.Type == kind {
			allowed = true
			break
		}
	}
	if !allowed || record.Token != token || !record.UserId.Valid() || user.UserId != record.UserId || record.CreatedTime.IsZero() {
		return "verificationLinkInvalid"
	}
	if !user.Available() {
		return "userDisabled"
	}
	age := now.Sub(record.CreatedTime)
	if age < 0 {
		return "verificationLinkNotYetValid"
	}
	if age.Hours() > (&TokenService{}).GetOverHours(record.Type) {
		return "verificationLinkExpired"
	}
	return ""
}
