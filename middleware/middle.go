package middleware

import (
	"MFile/db/mredis"
	"MFile/generate/limiter"
)

var redis_limiter limiter.T0kenLimiter
func init(){
	redis_limiter=limiter.NewTokenLimiter(mredis.RedisDb,"gable_limiter",10,30)
}

