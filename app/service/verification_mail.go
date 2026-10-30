package service

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/gemsnote/gemsnote/app/info"
)

// Verification emails are application-owned so imported legacy database
// templates cannot override their brand or the request's selected language.
func verificationLanguage(preference string) string {
	selected, best := "zh-CN", 0.0
	for _, item := range strings.Split(preference, ",") {
		parts := strings.Split(item, ";")
		tag := strings.ToLower(strings.TrimSpace(parts[0]))
		base := strings.SplitN(tag, "-", 2)[0]
		if base != "zh" && base != "en" {
			continue
		}
		quality := 1.0
		for _, parameter := range parts[1:] {
			key, value, ok := strings.Cut(strings.TrimSpace(parameter), "=")
			if ok && strings.EqualFold(key, "q") {
				parsed, err := strconv.ParseFloat(value, 64)
				if err != nil || parsed < 0 || parsed > 1 {
					quality = 0
				} else {
					quality = parsed
				}
			}
		}
		if quality > best {
			best = quality
			if base == "en" {
				selected = "en-US"
			} else {
				selected = "zh-CN"
			}
		}
	}
	return selected
}

type verificationMail struct {
	Language, Brand, Slogan, SiteURL, TokenURL             string
	Greeting, Instruction, Button, Fallback, Expiry, Setup string
}

var verificationMailHTML = template.Must(template.New("verification-mail").Parse(`<!doctype html>
<html lang="{{.Language}}"><head><meta charset="utf-8"></head>
<body style="margin:0;padding:24px 12px;background:#f5f6f3;color:#243b36;font-family:Arial,'Microsoft YaHei',sans-serif;line-height:1.7">
<div style="max-width:600px;margin:0 auto;padding:28px;background:#fff;border:1px solid #dde3d9;border-radius:10px">
<a href="{{.SiteURL}}" style="color:#193b35;font-size:24px;font-weight:bold;text-decoration:none">{{.Brand}}</a>
<hr style="border:0;border-top:1px solid #dde3d9;margin:20px 0">
<p>{{.Greeting}}<br>{{.Instruction}}</p>
<p style="margin:24px 0"><a href="{{.TokenURL}}" style="display:inline-block;padding:12px 28px;border-radius:6px;background:#193b35;color:#fff;font-weight:bold;text-decoration:none">{{.Button}}</a></p>
<p style="font-size:12px;color:#63766e">{{.Fallback}}<br><a href="{{.TokenURL}}" style="font-size:12px;color:#246f58;overflow-wrap:anywhere;word-break:break-all">{{.TokenURL}}</a></p>
<p>{{.Expiry}}{{if .Setup}}<br>{{.Setup}}{{end}}</p>
<hr style="border:0;border-top:1px solid #dde3d9;margin:28px 0 14px">
<p style="margin:0;font-size:12px;color:#63766e"><a href="{{.SiteURL}}" style="color:#193b35;text-decoration:none">{{.Brand}}</a><br>{{.Slogan}}</p>
</div></body></html>`))

func renderVerificationMail(user info.User, email, siteURL, tokenURL, preference string, kind int) (subject, html, plain string, err error) {
	initial := kind == info.TokenActiveEmail && user.Pwd == ""
	hours := int(info.ActiveEmailOverHours)
	if kind == info.TokenUpdateEmail {
		hours = int(info.UpdateEmailOverHours)
	}
	name := user.UsernameRaw
	if name == "" {
		name = user.Username
	}
	m := verificationMail{Language: verificationLanguage(preference), SiteURL: siteURL, TokenURL: tokenURL}
	if m.Language == "en-US" {
		m.Brand = "Gemsnote"
		m.Slogan = "Small notes become lasting insights."
		m.Greeting = fmt.Sprintf("Hello %s (%s),", name, email)
		m.Instruction = "Please verify your email address using the link below."
		m.Button = "Verify email"
		m.Fallback = "If the button doesn’t work, click the link below or copy it into your browser:"
		m.Expiry = fmt.Sprintf("This verification link is valid for %d hours.", hours)
		subject = "Verify your Gemsnote email address."
		if initial {
			subject = "Welcome to Gemsnote — please verify your email."
			m.Setup = "After verification, you’ll be taken to the password setup page to create your sign-in password."
		}
	} else {
		m.Brand = "珠玑笔记"
		m.Slogan = "日积字句，终得珠玑。"
		m.Greeting = fmt.Sprintf("%s（%s），您好：", name, email)
		m.Instruction = "请点击链接验证邮箱。"
		m.Button = "验证邮箱"
		m.Fallback = "如果按钮无法使用，请点击或复制以下链接到浏览器："
		m.Expiry = fmt.Sprintf("此验证链接 %d 小时内有效。", hours)
		subject = "请验证您的珠玑笔记邮箱。"
		if initial {
			subject = "恭喜您已注册珠玑笔记，请验证邮箱。"
			m.Setup = "验证成功后，您将进入密码设置页面，自行设置登录密码。"
		}
	}
	var out bytes.Buffer
	if err = verificationMailHTML.Execute(&out, m); err != nil {
		return
	}
	html = out.String()
	lines := []string{m.Brand + " (" + siteURL + ")", "----", m.Greeting, m.Instruction, "", m.Button + ": " + tokenURL, "", m.Expiry}
	if m.Setup != "" {
		lines = append(lines, m.Setup)
	}
	lines = append(lines, "", "----", m.Brand, m.Slogan)
	plain = strings.Join(lines, "\n")
	return
}
