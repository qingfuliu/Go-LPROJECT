package mredis

import (
	"MFile/db/mredis"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	acquire = `
local limiterKey = KEYS[1].."_limiter"

local limiter = tonumber(redis.call("get", limiterKey))


if (limiter == nil) then
    limiter = 100
    redis.call("set", limiterKey, limiter, "nx")
end

redis.call("zRemRangeByScore", KEYS[1], "-inf", ARGV[1])

local rs = tonumber(redis.call("zAdd",KEYS[1],"CH",ARGV[2], KEYS[2]))

if (rs ~= 1) then
    return "FAILED"
end

local rank=tonumber(redis.call("zRank",KEYS[1],KEYS[2]))

if (rank>=limiter)then
    redis.call("zRem",KEYS[1],KEYS[2])
    return "FAILED"
end

return "OK"
`
	release = `
local rs = tonumber(redis.call("zRem", KEYS[1], KEYS[2]))
if (rs == nil) or (rs~=1) then
    return "FAILED"
end
return "OK"
`
	refresh = `
local limiterKey = KEYS[1].."_limiter"

local limiter = tonumber(redis.call("get", limiterKey))


if (limiter == nil) then
    limiter = 100
    redis.call("set", limiterKey, limiter, "nx")
end

redis.call("zRemRangeByScore", KEYS[1], "-inf", ARGV[1])

local rs=tonumber(redis.call("zRank",KEYS[1],KEYS[2]))

if (rs==nil) then
    return "FAILED"
end

rs = tonumber(redis.call("zAdd",KEYS[1],"CH","XX",ARGV[2], KEYS[2]))

if (rs ~= 1) then
    return "FAILED"
end

return "OK"
`
)

/*
--ARGV 		  1			  2             3
--			zet key     unique id   nowTimeStamp             limiterKey=ARGV[2]+"_limiter"
--KEYS         1
--          expireDuration
*/

type RedisCond interface {
	Acquire() (bool, error)
	Release() (bool, error)
	Refresh() (bool, error)
}

type redisCond struct {
	redis          mredis.CmdAble
	id             string
	expireDuration time.Duration
	zSetKey        string
}

func NewRedisCond(redis mredis.CmdAble, id string, expireDuration time.Duration, zSetKey string) RedisCond {
	return &redisCond{
		redis:          redis,
		id:             id,
		expireDuration: expireDuration,
		zSetKey:        zSetKey,
	}
}

func (rCond *redisCond) Acquire() (bool, error) {

	/*
		--ARGV 		  1			  2             3
		--			zet key     unique id   nowTimeStamp             limiterKey=ARGV[2]+"_limiter"
		--KEYS         1
		--          expireDuration
	*/
	now := time.Now()
	rs, err := rCond.redis.Eval(context.Background(), acquire, []string{rCond.zSetKey, rCond.id},
		[]interface{}{now.Add(-rCond.expireDuration).Unix(), now.Unix()}).Result()
	if err != nil {
		return false, err
	} else if err == redis.Nil {
		return false, mredis.RedisKeyDoesNotExist
	} else if rs == nil {
		return false, mredis.RedisUnKnowError
	}
	if rsStr, ok := rs.(string); !ok {
		return false, mredis.RedisFormatError
	} else if rsStr == "OK" {
		return true, nil
	} else {
		return false, mredis.RedisCondOverLoad
	}
}

func (rCond *redisCond) Release() (bool, error) {
	rs, err := rCond.redis.Eval(context.Background(), release, []string{rCond.zSetKey, rCond.id}, []interface{}{}).Result()
	if err != nil {
		return false, err
	} else if err == redis.Nil {
		return false, err
	} else if rs == nil {
		return false, mredis.RedisUnKnowError
	}
	if rsStr, ok := rs.(string); !ok {
		return false, mredis.RedisFormatError
	} else if rsStr == "OK" {
		return true, nil
	}
	return false, mredis.RedisCondIsExpired
}

func (rCond *redisCond) Refresh() (bool, error) {
	now := time.Now()
	rs, err := rCond.redis.Eval(context.Background(), refresh, []string{rCond.zSetKey, rCond.id},
		[]interface{}{now.Add(-rCond.expireDuration).Unix(), now.Unix()}).Result()
	if err != nil {
		return false, err
	} else if err == redis.Nil {
		return false, mredis.RedisKeyDoesNotExist
	} else if rs == nil {
		return false, mredis.RedisUnKnowError
	}
	if rsStr, ok := rs.(string); !ok {
		return false, mredis.RedisFormatError
	} else if rsStr == "OK" {
		return true, nil
	} else {
		return false, mredis.RedisCondIsExpired
	}
}
