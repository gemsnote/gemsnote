package service

import (
	"github.com/gemsnote/gemsnote/app/info"
	"strings"
	"testing"
)

func TestVerificationLanguage(t *testing.T) {
	for input, want := range map[string]string{"": "zh-CN", "zh-CN": "zh-CN", "en-US": "en-US", "zh-TW,en;q=0.5": "zh-CN", "en;q=0.1,zh;q=0.9": "zh-CN", "zh;q=0,en;q=1": "en-US", "fr,en;q=0.8": "en-US", "en;q=invalid": "zh-CN"} {
		if got := verificationLanguage(input); got != want {
			t.Errorf("%q => %s, want %s", input, got, want)
		}
	}
}

func TestVerificationMailBrandLanguageAndSafety(t *testing.T) {
	user := info.User{Username: "alice", UsernameRaw: "Alice <script>alert(1)</script>"}
	for _, lang := range []string{"zh-CN", "en-US"} {
		for _, kind := range []int{info.TokenActiveEmail, info.TokenUpdateEmail} {
			subject, body, plain, err := renderVerificationMail(user, "alice@example.test", "https://notes.example.test", "https://notes.example.test/user/activeEmail?token=test-proof", lang, kind)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(body, "<script>") || !strings.Contains(body, "&lt;script&gt;") {
				t.Fatal("username not HTML escaped")
			}
			for _, bad := range []string{"leanote", "pearlnote", "your own cloud note"} {
				if strings.Contains(strings.ToLower(subject+body+plain), bad) {
					t.Fatalf("old brand %s", bad)
				}
			}
			if strings.Count(body, `href="https://notes.example.test/user/activeEmail?token=test-proof"`) != 2 {
				t.Fatal("button and fallback must share verification link")
			}
			if strings.Count(body, `href="https://notes.example.test"`) != 2 {
				t.Fatal("header and footer must link to configured site")
			}
			if !strings.Contains(body, `font-size:12px`) || !strings.Contains(body, `background:#193b35`) {
				t.Fatal("missing fallback typography or button")
			}
			if lang == "zh-CN" {
				if !strings.Contains(body, "您好：<br>请点击链接验证邮箱。") || !strings.Contains(plain, "您好：\n请点击链接验证邮箱。") {
					t.Fatal("Chinese greeting must have line break")
				}
				if !strings.Contains(body, "日积字句，终得珠玑。") {
					t.Fatal("missing Chinese slogan")
				}
				if kind == info.TokenActiveEmail && subject != "恭喜您已注册珠玑笔记，请验证邮箱。" {
					t.Fatal(subject)
				}
			} else {
				if !strings.Contains(body, "),<br>Please verify") || !strings.Contains(plain, "),\nPlease verify") {
					t.Fatal("English greeting must have line break")
				}
				if !strings.Contains(body, "Small notes become lasting insights.") {
					t.Fatal("missing English slogan")
				}
				if kind == info.TokenActiveEmail && subject != "Welcome to Gemsnote — please verify your email." {
					t.Fatal(subject)
				}
			}
			if kind == info.TokenActiveEmail && !strings.Contains(body, "48") {
				t.Fatal("registration links expire in 48 hours")
			}
			if kind == info.TokenUpdateEmail && (strings.Contains(body, "48") || strings.Contains(body, "password setup") || strings.Contains(body, "密码设置")) {
				t.Fatal("email changes must not use registration expiration or password setup instructions")
			}
		}
	}
	user.Pwd = "existing-password-hash"
	_, body, _, err := renderVerificationMail(user, "alice@example.test", "https://notes.example.test", "https://notes.example.test/verify", "zh-CN", info.TokenActiveEmail)
	if err != nil || strings.Contains(body, "密码设置") {
		t.Fatal("existing users must not be instructed to set an initial password")
	}
}
