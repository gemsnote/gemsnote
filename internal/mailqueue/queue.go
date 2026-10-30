// Package mailqueue implements the private mail payload store and Asynq worker.
package mailqueue

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

const Queue = "gemsnote.mail"
const TaskType = "gemsnote.mail.send"
const PayloadTTL = 10 * time.Minute
const StatusTTL = 24 * time.Hour

type SMTPSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}
type Payload struct {
	SMTP    SMTPSettings `json:"smtp"`
	To      []string     `json:"to"`
	Subject string       `json:"subject"`
	Body    string       `json:"body"`
	Text    string       `json:"text,omitempty"`
}
type Owner struct {
	UserID  string
	Created int64
}
type Status struct {
	TaskID string `json:"TaskId,omitempty"`
	State  string `json:"State"`
	Error  string `json:"Error,omitempty"`
}

// Errors are intentionally fixed codes: never propagate Redis URLs or SMTP replies.
func NewRedis(url string) (*redis.Client, error) {
	if url == "" {
		return nil, errors.New("mailQueueNotConfigured")
	}
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, errors.New("mailQueueNotConfigured")
	}
	opt.DialTimeout = 2 * time.Second
	opt.ReadTimeout = 2 * time.Second
	opt.WriteTimeout = 2 * time.Second
	opt.ContextTimeoutEnabled = true
	opt.MaxRetries = -1
	return redis.NewClient(opt), nil
}
func Enqueue(ctx context.Context, url, id, owner string, payload Payload) error {
	client, err := NewRedis(url)
	if err != nil {
		return err
	}
	defer client.Close()
	data, err := json.Marshal(payload)
	if err != nil {
		return errors.New("invalidEmail")
	}
	metadata, _ := json.Marshal(Owner{owner, time.Now().Unix()})
	// Store private data first, so a worker can never see a task without its payload.
	_, err = client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Set(ctx, "gemsnote:mail:owner:"+id, metadata, StatusTTL)
		p.Set(ctx, "gemsnote:mail:payload:"+id, data, PayloadTTL)
		return nil
	})
	if err != nil {
		return errors.New("mailQueueUnavailable")
	}
	broker := asynq.NewClientFromRedisClient(client)
	_, err = broker.EnqueueContext(ctx, asynq.NewTask(TaskType, []byte(id)),
		asynq.TaskID(id), asynq.Queue(Queue), asynq.MaxRetry(3), asynq.Timeout(time.Minute))
	// On an ambiguous network failure keep private keys until their TTL: Redis may
	// already have accepted the task. Never delete a possibly active payload here.
	if err != nil {
		return errors.New("mailQueueUnavailable")
	}
	return nil
}

// A persisted claim prevents SMTP replays when an interrupted task is recovered.
// SMTP cannot guarantee exactly-once delivery. Unknown outcomes require a manual resend.
var claim = redis.NewScript(`
local old = redis.call('GET',KEYS[1])
if old then return old end
redis.call('SET',KEYS[1],ARGV[1],'EX',ARGV[2])
return ''`)
var finish = redis.NewScript(`
redis.call('SET',KEYS[1],ARGV[1],'EX',ARGV[2])
redis.call('DEL',KEYS[2])
return 1`)

type Handler struct{ Redis *redis.Client }

func (h *Handler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	id := string(task.Payload())
	if _, err := hex.DecodeString(id); err != nil || len(id) != 24 {
		return asynq.SkipRetry
	}
	statusKey, payloadKey := "gemsnote:mail:status:"+id, "gemsnote:mail:payload:"+id
	sending, _ := json.Marshal(Status{State: "sending"})
	previous, err := claim.Run(ctx, h.Redis, []string{statusKey}, string(sending), int(StatusTTL.Seconds())).Text()
	if err != nil {
		return errors.New("mailQueueUnavailable")
	}
	if previous != "" {
		var status Status
		if json.Unmarshal([]byte(previous), &status) == nil && (status.State == "sent" || status.State == "failed") {
			return nil
		}
		return h.complete(id, "mailDeliveryFailed")
	}
	raw, err := h.Redis.Get(ctx, payloadKey).Bytes()
	if err == redis.Nil {
		return h.complete(id, "mailTaskExpired")
	}
	if err != nil {
		return errors.New("mailQueueUnavailable")
	}
	var payload Payload
	if json.Unmarshal(raw, &payload) != nil {
		return h.complete(id, "mailDeliveryFailed")
	}
	code := ""
	if err = Send(ctx, id, payload); err != nil {
		code = err.Error()
	}
	return h.complete(id, code)
}
func (h *Handler) complete(id, code string) error {
	status := Status{State: "sent"}
	if code != "" {
		status.State = "failed"
		status.Error = code
		log.Printf("Mail task %s failed: %s", id, code)
	}
	data, _ := json.Marshal(status)
	// Persist the result even if SMTP exhausted the task's context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := finish.Run(ctx, h.Redis, []string{"gemsnote:mail:status:" + id, "gemsnote:mail:payload:" + id}, string(data), int(StatusTTL.Seconds())).Err(); err != nil {
		return errors.New("mailQueueUnavailable")
	}
	return nil // SMTP failures are final; only infrastructure errors are retried.
}
