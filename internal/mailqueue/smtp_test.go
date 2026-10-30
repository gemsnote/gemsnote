package mailqueue

import (
	"bytes"
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/mail"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMessageUnicodeAndAlternatives(t *testing.T) {
	p := Payload{SMTP: SMTPSettings{Username: "sender@example.test"}, To: []string{"用户 <recipient@example.test>"}, Subject: "恭喜您已注册珠玑笔记，请验证邮箱。", Text: "您好，\n请验证邮箱。", Body: "<p>您好，<br>请验证邮箱。</p>"}
	raw, from, to, err := message("0123456789abcdef01234567", p)
	if err != nil || from != "sender@example.test" || len(to) != 1 || to[0] != "recipient@example.test" {
		t.Fatalf("message envelope: %v", err)
	}
	m, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	subject, err := (&mime.WordDecoder{}).DecodeHeader(m.Header.Get("Subject"))
	if err != nil || subject != p.Subject {
		t.Fatal("Unicode subject changed")
	}
	kind, params, err := mime.ParseMediaType(m.Header.Get("Content-Type"))
	if err != nil || kind != "multipart/alternative" {
		t.Fatal("missing MIME alternatives")
	}
	r := multipart.NewReader(m.Body, params["boundary"])
	for i, want := range []string{p.Text, p.Body} {
		part, err := r.NextPart()
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(part)
		if err != nil || strings.ReplaceAll(string(body), "\r\n", "\n") != want {
			t.Fatalf("part %d changed: %q %v", i, body, err)
		}
	}
	if _, err = r.NextPart(); err != io.EOF {
		t.Fatal("unexpected MIME part")
	}
}
func TestMessageRejectsHeaderInjection(t *testing.T) {
	valid := Payload{SMTP: SMTPSettings{Username: "sender@example.test"}, To: []string{"recipient@example.test"}, Subject: "Hello"}
	for _, field := range []string{"subject", "sender", "recipient"} {
		p := valid
		switch field {
		case "subject":
			p.Subject = "Hello\r\nBcc: victim@example.test"
		case "sender":
			p.SMTP.Username = "sender@example.test\r\nBcc: victim@example.test"
		case "recipient":
			p.To = []string{"recipient@example.test\r\nBcc: victim@example.test"}
		}
		if _, _, _, err := message("id", p); err == nil {
			t.Fatalf("accepted injection in %s", field)
		}
	}
}
func TestSMTPStalledGreetingHonorsContext(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, err := listener.Accept()
		if err == nil {
			defer conn.Close()
			io.Copy(io.Discard, conn)
		}
	}()
	p := Payload{SMTP: SMTPSettings{Host: "127.0.0.1", Port: listener.Addr().(*net.TCPAddr).Port, Username: "sender@example.test"}, To: []string{"recipient@example.test"}}
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	started := time.Now()
	err = Send(ctx, "id", p)
	if err == nil || err.Error() != "mailConnectionFailed" || time.Since(started) > 2*time.Second {
		t.Fatalf("SMTP did not stop promptly: %v", err)
	}
	<-done
}

func TestCustomCAMustBeReadable(t *testing.T) {
	t.Setenv("SSL_CERT_FILE", filepath.Join(t.TempDir(), "missing.pem"))
	if _, err := smtpTLSConfig("localhost"); err == nil {
		t.Fatal("unreadable custom CA silently ignored")
	}
}
