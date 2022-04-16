package limiter

import (
	"MFile/db/mredis"
	"MFile/logger"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	xrate "golang.org/x/time/rate"
	"sync"
	"sync/atomic"
	"time"
)

const (
	//KEY[1] storedTokenKey
	//KEY[2] freeTimeStampKey
	//ARGV[1] now
	//ARGV[2] requestN
	//ARGV[3] capacity
	//ARGV[4] rate
	requestCommend = `
local nowTime=tonumber(ARGV[1])
local requestn=tonumber(ARGV[2])
local capacity=tonumber(ARGV[3])
local rate=tonumber(ARGV[4])

local nextFreeTIme=tonumber(redis.call("get",KEYS[2]))
if nextFreeTIme==nil then
    nextFreeTIme=0
end

local storedTokens=tonumber(redis.call("get",KEYS[1]))
if storedTokens==nil then
    storedTokens=capacity
end

if nowTime<nextFreeTIme then
return "FAILD"
else
    storedTokens=math.max(capacity,storedTokens+(nowTime-nextFreeTIme)*rate)
    local access=math.max(storedTokens,requestn)

    local diff=requestn-access
    storedTokens=storedTokens-access
    nextFreeTIme=nowTime+math.ceil(diff/rate)
    redis.call("Set",KEYS[1],storedTokens)
    redis.call("Set",KEYS[2],nextFreeTIme)
return "OK"
end
`
	storedTokenKeyFormat="{%s}.StoredToken"
	freeTimeStampKeyFormat="{%s}.nextFreeTIme"
)

func NewTokenLimiter(r mredis.CmdAble,key string,rate,capacity int)T0kenLimiter{
	return &tokenBucket{
		redis: r,
		capacity: capacity,
		rate: rate,
		storedTokenKey: fmt.Sprintf(storedTokenKeyFormat,key),
		freeTimeStampKey: fmt.Sprintf(freeTimeStampKeyFormat,key),
		redisAlive:true,
		monitorStarted:0,
		tempLimiter:	xrate.NewLimiter(xrate.Limit(rate),capacity),
	}
}

type tokenBucket struct {
	capacity int
	rate     int
	//redis
	redis            mredis.CmdAble
	storedTokenKey   string
	freeTimeStampKey string
	redisAlive       bool
	//temp
	monitorStarted int32
	tempLimiter    *xrate.Limiter
	//protected
	mutex sync.Mutex
}

type T0kenLimiter interface {
	AllowN(n int) bool
	Allow() bool
	reserve(now time.Time, n int) bool
	startMonitor()
	waitForRedis()
}

func (t *tokenBucket) AllowN(n int) bool {
	return t.reserve(time.Now(), n)
}

func (t *tokenBucket) Allow() bool {
	return t.AllowN(1)
}

func (t *tokenBucket) reserve(now time.Time, n int) bool {
	if atomic.LoadInt32(&t.monitorStarted) == 1 {
		return t.tempLimiter.AllowN(now, 1)
	}

	rs, err := t.redis.Eval(context.Background(), requestCommend, []string{t.storedTokenKey, t.freeTimeStampKey}, now.Unix(), n, t.capacity, t.rate).Result()

	if err == redis.Nil {
		return false
	}
	if err != nil {
		logger.MLogger.Error("token limiter failed!", zap.Error(err))
		t.startMonitor()
		return t.tempLimiter.AllowN(now, n)
	}

	code, ok := rs.(string)
	if !ok {
		logger.MLogger.Error("token limiter failed!", zap.Error(err))
		t.startMonitor()
		return t.tempLimiter.AllowN(now, n)
	}

	return code == "OK"
}

func (t *tokenBucket) startMonitor() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if !t.redisAlive {
		return
	}
	t.redisAlive = false
	atomic.StoreInt32(&t.monitorStarted, 1)

	go t.waitForRedis()
}
func (t *tokenBucket) waitForRedis() {
	tiker := time.NewTicker(time.Millisecond * 100)
	defer func() {
		tiker.Stop()
		atomic.StoreInt32(&t.monitorStarted, 0)
		t.mutex.Lock()
		t.redisAlive = true
		t.mutex.Unlock()
	}()
	for {
		if rs, err := t.redis.Ping(context.Background()).Result(); err == nil && rs == "PONG" {
			break
		}
	}
}
