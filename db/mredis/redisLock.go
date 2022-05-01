package mredis

import (
	"MFile/generate"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//KEYS[1] key
//ARGV[1] id
//ARGV[2] expire
var unKnowErr = fmt.Errorf("unKnow error")

const (
	maxIdLength = 16
	lockCommend = `if redis.call("GET",KEYS[1])==ARGV[1] then
	redis.call("expire",KEYS[1],ARGV[2])
	return "ok"
else
	return redis.call("SET",KEYS[1],ARGV[1],"NX","EX",ARGV[2])
end`
	unLockCommend = `if redis.call("GET",KEYS[1])==ARGV[1] then
	return redis.call("DEL",KEYS[1])
else
	return 0
end
`
)

type RedisLock interface {
	Lock() (bool, error)
	Unlock() (bool, error)
	SetExpire(duration time.Duration)
	GetExpire() time.Duration
	GetId() string
	GetKey() string
	SetKey(key string)
}

type lock struct {
	redisClient CmdAble
	ctx         context.Context
	expire      time.Duration
	id          string
	key         string
}

func (rl *lock) Lock() (bool, error) {
	rs, err := rl.redisClient.Eval(rl.ctx, lockCommend, []string{rl.key}, rl.id, rl.expire).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else if rs == nil {
		return false, err
	}

	str, ok := rs.(string)
	if ok && str == "OK" {
		return true, nil
	}
	return false, unKnowErr
}
func (rl *lock) Unlock() (bool, error) {
	rs, err := rl.redisClient.Eval(rl.ctx, unLockCommend, []string{rl.key}, rl.id).Result()

	if err != nil || err == redis.Nil {
		return false, err
	}

	reply, ok := rs.(int64)

	if ok && reply == 1 {
		return true, nil
	}

	return false, unKnowErr
}
func (rl *lock) SetExpire(duration time.Duration) {
	rl.expire = duration
}
func (rl *lock) GetExpire() time.Duration {
	return rl.expire
}
func (rl *lock) GetId() string {
	return rl.id
}
func (rl *lock) GetKey() string {
	return rl.key
}
func (rl *lock) SetKey(key string) {
	rl.key = key
}

func NewRedisLock(ctx context.Context, redis CmdAble, key string) RedisLock {

	return &lock{
		ctx:         ctx,
		expire:      30,
		key:         key,
		redisClient: redis,
		id:          generate.RandStringN(maxIdLength),
	}
}
