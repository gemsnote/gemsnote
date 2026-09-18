package tests

import (
	"testing"
	//	. "github.com/gemsnote/gemsnote/app/lea"
	"github.com/gemsnote/gemsnote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

// 测试登录
func TestSendMail(t *testing.T) {
	requireMongoIntegration(t)
	ok, err := service.EmailS.SendEmail("life@gemsnote.com", "你好", "你好吗")
	t.Log(ok)
	t.Log(err)
}
