package mredis

import "errors"

var (
	RedisUnKnowError     = errors.New("unKnow redis error")
	RedisFormatError     = errors.New("redis format error")
	RedisKeyDoesNotExist = errors.New("key does not exist")
	RedisCondOverLoad    = errors.New("redis cond Over Load")
	RedisCondIsExpired   = errors.New("redis cond is expired")
)
