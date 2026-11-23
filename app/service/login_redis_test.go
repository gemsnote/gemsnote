package service

import (
	"github.com/garyburd/redigo/redis"
	"os"
	"sync"
	"testing"
	"time"
)

func TestLoginRedisIntegration(t *testing.T) {
	if os.Getenv("GEMSNOTE_LOGIN_TEST_REDIS") != "redis://127.0.0.1:16391/0" {
		t.Skip("requires disposable login Redis on port 16391")
	}
	t.Setenv("GEMSNOTE_LOGIN_REDIS_URL", os.Getenv("GEMSNOTE_LOGIN_TEST_REDIS"))
	c, err := redis.DialURL(os.Getenv("GEMSNOTE_LOGIN_TEST_REDIS"))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	identity := time.Now().String()
	ip := identity
	keys := []string{"gemsnote:login:" + loginAccountKey(identity, ip), "gemsnote:login:" + loginIPKey(ip)}
	defer c.Do("DEL", keys[0], keys[1])
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, _, err := redisLoginAttempt("failure", identity, ip); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if count, err := redis.Int(c.Do("GET", keys[0])); err != nil || count != 10 {
		t.Fatalf("atomic count=%d err=%v", count, err)
	}
	if allowed, retry, err := redisLoginAttempt("allow", identity, ip); err != nil || allowed || retry <= 0 {
		t.Fatalf("allow=%v retry=%v err=%v", allowed, retry, err)
	}
	redisLoginAttempt("success", identity, ip)
	if allowed, _, err := redisLoginAttempt("allow", identity, ip); err != nil || !allowed {
		t.Fatal("success did not clear account limit", err)
	}
	c.Do("PEXPIRE", keys[1], 1)
	time.Sleep(5 * time.Millisecond)
	if exists, _ := redis.Bool(c.Do("EXISTS", keys[1])); exists {
		t.Fatal("counter did not expire")
	}
}
