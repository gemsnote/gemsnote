package service

// Celery's Redis transport uses a JSON envelope around its protocol-v2 body.
// Only the task ID enters the broker: SMTP credentials live in a short-lived
// private payload key, never in task arguments, results or application logs.
import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

const mailQueue = "gemsnote.mail"
const mailTTL = 600

type MailStatus struct {
	TaskID string `json:"TaskId"`
	State  string `json:"State"`
	Error  string `json:"Error,omitempty"`
}
type mailOwner struct {
	UserID  string
	Created int64
}
type smtpSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}
type mailPayload struct {
	SMTP    smtpSettings `json:"smtp"`
	To      []string     `json:"to"`
	Subject string       `json:"subject"`
	Body    string       `json:"body"`
}

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
func celeryMailMessage(id string, now time.Time) ([]byte, error) {
	body, err := json.Marshal([]interface{}{[]string{id}, map[string]interface{}{}, map[string]interface{}{"callbacks": nil, "errbacks": nil, "chain": nil, "chord": nil}})
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]interface{}{
		"body": base64.StdEncoding.EncodeToString(body), "content-encoding": "utf-8", "content-type": "application/json",
		"headers":    map[string]interface{}{"lang": "py", "task": "gemsnote.mail.send", "id": id, "root_id": id, "retries": 0, "expires": now.Add(mailTTL * time.Second).UTC().Format(time.RFC3339), "argsrepr": "<mail task>", "kwargsrepr": "{}"},
		"properties": map[string]interface{}{"correlation_id": id, "delivery_mode": 2, "delivery_tag": id, "body_encoding": "base64", "delivery_info": map[string]string{"exchange": "", "routing_key": mailQueue}},
	})
}
func (this *EmailService) QueueEmail(owner, to, subject, body string) (string, error) {
	port, _ := strconv.Atoi(configService.GetGlobalStringConfig("emailPort"))
	payload := mailPayload{SMTP: smtpSettings{Host: configService.GetGlobalStringConfig("emailHost"), Port: port, Username: configService.GetGlobalStringConfig("emailUsername"), Password: configService.GetGlobalStringConfig("emailPassword"), SSL: configService.GetGlobalStringConfig("emailSSL") == "1"}, Subject: subject, Body: body}
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
	c, err := mailConnect()
	if err != nil {
		return "", err
	}
	defer c.Close()
	id := bson.NewObjectId().Hex()
	now := time.Now()
	envelope, err := celeryMailMessage(id, now)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	metadata, _ := json.Marshal(mailOwner{owner, now.Unix()})
	// Atomically publish metadata, payload and task, in one transaction.
	script := redis.NewScript(0, `
 redis.call('SET','gemsnote:mail:owner:'..ARGV[2],ARGV[3],'EX',86400)
 redis.call('SET','gemsnote:mail:payload:'..ARGV[2],ARGV[4],'EX',600)
 redis.call('LPUSH','gemsnote.mail',ARGV[5])
 return 1`)
	result, err := redis.Int(script.Do(c, owner, id, metadata, data, envelope))
	if err != nil {
		return "", errors.New("mailQueueUnavailable")
	}
	if result != 1 {
		return "", errors.New("mailRateLimited")
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
