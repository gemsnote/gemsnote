package service

import (
	"os"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/revel/revel"
)

// Both counters are changed atomically; expiration is set only on the first
// failure so further guesses cannot prolong another user's lockout forever.
var loginAttemptScript = redis.NewScript(2, `
if ARGV[1] == 'success' then redis.call('DEL', KEYS[1]); return 0 end
if ARGV[1] == 'failure' then
  for i=1,2 do
    local count=redis.call('INCR', KEYS[i])
    if count == 1 then redis.call('PEXPIRE', KEYS[i], ARGV[2]) end
  end
  return 0
end
local retry=0
for i=1,2 do
  local count=tonumber(redis.call('GET', KEYS[i]) or '0')
  if count >= tonumber(ARGV[i+2]) then
    local ttl=redis.call('PTTL', KEYS[i])
    if ttl < 0 then redis.call('PEXPIRE', KEYS[i], ARGV[2]); ttl=tonumber(ARGV[2]) end
    retry=math.max(retry,ttl)
  end
end
return retry
`)

var loginRedisWarning struct {
	sync.Mutex
	last time.Time
}

func redisLoginAttempt(action, identity, clientIP string) (bool, time.Duration, error) {
	url := os.Getenv("GEMSNOTE_LOGIN_REDIS_URL")
	if url == "" && revel.Config != nil {
		url = revel.Config.StringDefault("login.redis_url", "")
	}
	if url == "" {
		url = mailRedisURL()
	}
	c, err := redis.DialURL(url, redis.DialConnectTimeout(time.Second), redis.DialReadTimeout(time.Second), redis.DialWriteTimeout(time.Second))
	if err == nil {
		defer c.Close()
		var retry int64
		retry, err = redis.Int64(loginAttemptScript.Do(c,
			"gemsnote:login:"+loginAccountKey(identity, clientIP), "gemsnote:login:"+loginIPKey(clientIP),
			action, int64(loginFailureWindow()/time.Millisecond), loginFailureLimit("user:"), loginFailureLimit("ip:")))
		if err == nil {
			return retry == 0, time.Duration(retry) * time.Millisecond, nil
		}
	}
	loginRedisWarning.Lock()
	if time.Since(loginRedisWarning.last) > time.Minute {
		if revel.AppLog != nil {
			revel.AppLog.Warn("Login rate limit Redis unavailable; using local counters")
		}
		loginRedisWarning.last = time.Now()
	}
	loginRedisWarning.Unlock()
	return false, 0, err
}
