package service

// Queue task arguments contain only IDs; credentials stay in short-lived private keys.
import (
	"context"
	"encoding/json"
	"errors"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gemsnote/gemsnote/internal/mailqueue"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

const mailTTL = int64(mailqueue.PayloadTTL / time.Second)

type MailStatus = mailqueue.Status
type mailOwner = mailqueue.Owner
type smtpSettings = mailqueue.SMTPSettings
type mailPayload = mailqueue.Payload

func mailRedisURL() string {
	if value := os.Getenv("GEMSNOTE_MAIL_REDIS_URL"); value != "" {
		return value
	}
	if revel.Config != nil {
		value, _ := revel.Config.String("mail.redis_url")
		return value
	}
	return ""
}
func mailConnect() (redis.Conn, error) {
	url := mailRedisURL()
	if url == "" {
		return nil, errors.New("mailQueueNotConfigured")
	}
	c, err := redis.DialURL(url, redis.DialConnectTimeout(2*time.Second), redis.DialReadTimeout(2*time.Second), redis.DialWriteTimeout(2*time.Second))
	if err != nil {
		return nil, errors.New("mailQueueUnavailable")
	}
	return c, nil
}
func (this *EmailService) QueueEmail(owner, to, subject, body string) (string, error) {
	return this.queueEmail(owner, to, subject, body, "")
}
func (this *EmailService) queueEmail(owner, to, subject, body, plain string) (string, error) {
	port, _ := strconv.Atoi(configService.GetGlobalStringConfig("emailPort"))
	payload := mailPayload{SMTP: smtpSettings{Host: configService.GetGlobalStringConfig("emailHost"), Port: port, Username: configService.GetGlobalStringConfig("emailUsername"), Password: configService.GetGlobalStringConfig("emailPassword"), SSL: configService.GetGlobalStringConfig("emailSSL") == "1"}, Subject: subject, Body: body, Text: plain}
	if payload.SMTP.Host == "" || port < 1 || port > 65535 || payload.SMTP.Username == "" || payload.SMTP.Password == "" {
		return "", errors.New("mailNotConfigured")
	}
	if _, err := mail.ParseAddress(payload.SMTP.Username); err != nil {
		return "", errors.New("mailNotConfigured")
	}
	if strings.ContainsAny(subject, "\r\n") {
		return "", errors.New("invalidEmail")
	}
	for _, recipient := range strings.Split(to, ";") {
		address, err := mail.ParseAddress(strings.TrimSpace(recipient))
		if err != nil {
			return "", errors.New("invalidEmail")
		}
		payload.To = append(payload.To, address.Address)
	}
	return enqueueMail(owner, payload)
}
func enqueueMail(owner string, payload mailPayload) (string, error) {
	id := bson.NewObjectId().Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mailqueue.Enqueue(ctx, mailRedisURL(), id, owner, payload); err != nil {
		return "", err
	}
	return id, nil
}
func GetMailStatus(owner, id string) (MailStatus, error) {
	result := MailStatus{TaskID: id, State: "queued"}
	if owner == "" || !bson.IsObjectIdHex(id) {
		return result, errors.New("mailTaskNotFound")
	}
	c, err := mailConnect()
	if err != nil {
		return result, err
	}
	defer c.Close()
	data, err := redis.Bytes(c.Do("GET", "gemsnote:mail:owner:"+id))
	if err == redis.ErrNil {
		return result, errors.New("mailTaskNotFound")
	}
	if err != nil {
		return result, errors.New("mailQueueUnavailable")
	}
	var metadata mailOwner
	if json.Unmarshal(data, &metadata) != nil || metadata.UserID != owner {
		return result, errors.New("mailTaskNotFound")
	}
	data, err = redis.Bytes(c.Do("GET", "gemsnote:mail:status:"+id))
	if err == nil {
		if json.Unmarshal(data, &result) != nil {
			return result, errors.New("mailQueueUnavailable")
		}
		result.TaskID = id
	} else if err != redis.ErrNil {
		return result, errors.New("mailQueueUnavailable")
	}
	if (result.State == "queued" || result.State == "sending") && time.Now().Unix()-metadata.Created > mailTTL {
		result.State = "failed"
		result.Error = "mailTaskExpired"
	}
	return result, nil
}
