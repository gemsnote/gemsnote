package tests

import (
	"testing"
	//	. "github.com/gemsnote/gemsnote/app/lea"
	"github.com/gemsnote/gemsnote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

// 测试登录
func TestAuth(t *testing.T) {
	requireMongoIntegration(t)
	_, err := service.AuthS.Login("admin", "abc123")
	if err != nil {
		t.Error("Admin User Auth Error")
	}
}
