package service

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2/bson"
)

// Uses an isolated Redis, a real Asynq worker and a local TLS SMTP fixture.
// Never run against a production queue.
func TestAsynqDeliveryIntegration(t *testing.T) {
	url := os.Getenv("GEMSNOTE_MAIL_TEST_REDIS")
	if url == "" {
		t.Skip("set GEMSNOTE_MAIL_TEST_REDIS for disposable mail integration")
	}
	if url != "redis://127.0.0.1:16389/0" {
		t.Fatal("only the isolated local test Redis is allowed")
	}
	t.Setenv("GEMSNOTE_MAIL_REDIS_URL", url)
	for _, test := range []struct{ name, to, password, state, code string }{
		{"sent", "recipient@example.test", "test-secret", "sent", ""},
		{"starttls", "recipient@example.test", "test-secret", "sent", ""},
		{"rejected", "reject@example.test", "test-secret", "failed", "mailRecipientRejected"},
		{"auth", "recipient@example.test", "wrong-secret", "failed", "mailAuthenticationFailed"},
	} {
		t.Run(test.name, func(t *testing.T) {
			port := 10258
			implicitTLS := test.name != "starttls"
			if !implicitTLS {
				port = 10259
			}
			id, err := enqueueMail("owner", mailPayload{SMTP: smtpSettings{Host: "localhost", Port: port, Username: "sender@example.test", Password: test.password, SSL: implicitTLS}, To: []string{test.to}, Subject: "测试 verification", Body: "<p>verification</p>"})
			if err != nil {
				t.Fatal(err)
			}
			if _, err = GetMailStatus("different-owner", id); err == nil {
				t.Fatal("another user read task state")
			}
			deadline := time.Now().Add(15 * time.Second)
			var status MailStatus
			for time.Now().Before(deadline) {
				status, err = GetMailStatus("owner", id)
				if err != nil {
					t.Fatal(err)
				}
				if status.State == "sent" || status.State == "failed" {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			if status.State != test.state || status.Error != test.code {
				t.Fatalf("unexpected task result %+v", status)
			}
			c, err := mailConnect()
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()
			// Worker deletes credentials after finishing; allow a short scheduling gap.
			time.Sleep(50 * time.Millisecond)
			exists, _ := redis.Bool(c.Do("EXISTS", "gemsnote:mail:payload:"+id))
			if exists {
				t.Fatal("SMTP payload retained after delivery")
			}
		})
	}
	c, err := mailConnect()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	id := bson.NewObjectId().Hex()
	data, _ := json.Marshal(mailOwner{"owner", time.Now().Add(-11 * time.Minute).Unix()})
	if _, err = c.Do("SETEX", "gemsnote:mail:owner:"+id, 60, data); err != nil {
		t.Fatal(err)
	}
	status, err := GetMailStatus("owner", id)
	if err != nil || status.Error != "mailTaskExpired" {
		t.Fatalf("worker outage not reported: %+v %v", status, err)
	}
}
