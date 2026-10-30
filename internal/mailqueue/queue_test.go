package mailqueue

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func testRedis(t *testing.T) *redis.Client {
	t.Helper()
	url := os.Getenv("GEMSNOTE_MAIL_TEST_REDIS")
	if url == "" {
		t.Skip("requires disposable Redis fixture")
	}
	if url != "redis://127.0.0.1:16389/0" {
		t.Fatal("only isolated test Redis is allowed")
	}
	client, err := NewRedis(url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { client.Close() })
	return client
}
func taskID(t *testing.T) string {
	t.Helper()
	var id [12]byte
	if _, err := rand.Read(id[:]); err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(id[:])
}
func TestRedisExpiredAndRecoveredTasks(t *testing.T) {
	client := testRedis(t)
	handler := Handler{Redis: client}
	ctx := context.Background()
	for _, test := range []struct{ initial, want, code string }{{"", "failed", "mailTaskExpired"}, {"sending", "failed", "mailDeliveryFailed"}, {"sent", "sent", ""}, {"failed", "failed", "mailRecipientRejected"}} {
		t.Run(test.initial, func(t *testing.T) {
			id := taskID(t)
			key := "gemsnote:mail:status:" + id
			defer client.Del(ctx, key, "gemsnote:mail:payload:"+id)
			if test.initial != "" {
				value, _ := json.Marshal(Status{State: test.initial, Error: test.code})
				if err := client.Set(ctx, key, value, StatusTTL).Err(); err != nil {
					t.Fatal(err)
				}
			}
			if test.initial == "sending" {
				client.Set(ctx, "gemsnote:mail:payload:"+id, "private", PayloadTTL)
			}
			task := asynq.NewTask(TaskType, []byte(id))
			if err := handler.ProcessTask(ctx, task); err != nil {
				t.Fatal(err)
			}
			raw, err := client.Get(ctx, key).Bytes()
			if err != nil {
				t.Fatal(err)
			}
			var status Status
			if json.Unmarshal(raw, &status) != nil || status.State != test.want || status.Error != test.code {
				t.Fatalf("unexpected result: %s", raw)
			}
			if count, _ := client.Exists(ctx, "gemsnote:mail:payload:"+id).Result(); count != 0 {
				t.Fatal("private payload retained")
			}
			if err := handler.ProcessTask(ctx, task); err != nil {
				t.Fatal(err)
			}
			again, _ := client.Get(ctx, key).Bytes()
			if string(again) != string(raw) {
				t.Fatal("terminal task replay changed result")
			}
		})
	}
}
func TestRedisTaskContainsOnlyID(t *testing.T) {
	client := testRedis(t)
	ctx := context.Background()
	inspector := asynq.NewInspectorFromRedisClient(client)
	if err := inspector.PauseQueue(Queue); err != nil {
		t.Fatal(err)
	}
	defer inspector.UnpauseQueue(Queue)
	id := taskID(t)
	defer client.Del(ctx, "gemsnote:mail:owner:"+id, "gemsnote:mail:payload:"+id)
	payload := Payload{SMTP: SMTPSettings{Username: "secret@example.test", Password: "private-password"}, Body: "private-verification-link"}
	if err := Enqueue(ctx, os.Getenv("GEMSNOTE_MAIL_TEST_REDIS"), id, "owner", payload); err != nil {
		t.Fatal(err)
	}
	defer inspector.DeleteTask(Queue, id)
	task, err := inspector.GetTaskInfo(Queue, id)
	if err != nil {
		t.Fatal(err)
	}
	if task.Type != TaskType || string(task.Payload) != id || task.MaxRetry != 3 {
		t.Fatalf("unexpected Asynq task: %+v", task)
	}
	for key, maximum := range map[string]time.Duration{"gemsnote:mail:payload:" + id: PayloadTTL, "gemsnote:mail:owner:" + id: StatusTTL} {
		ttl, err := client.TTL(ctx, key).Result()
		if err != nil || ttl <= 0 || ttl > maximum {
			t.Fatalf("invalid TTL %s: %s %v", key, ttl, err)
		}
	}
	// Inspect actual broker hash values, not just the decoded task payload.
	values, err := client.HGetAll(ctx, "asynq:{"+Queue+"}:t:"+id).Result()
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range values {
		if strings.Contains(value, "private-") || strings.Contains(value, "secret@example.test") {
			t.Fatal("private mail data leaked into broker")
		}
	}
}
