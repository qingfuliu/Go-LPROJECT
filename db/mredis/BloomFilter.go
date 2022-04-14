package mredis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type mOperation interface {
	Process(ctx context.Context, cmd redis.Cmder) error
}

func BfReserve(c mOperation, ctx context.Context, key string, errorRate float64, initialSize int) *redis.BoolCmd {
	argsSlice := make([]interface{}, 4)
	argsSlice[0] = "bf.reserve"
	argsSlice[1] = key
	argsSlice[2] = errorRate
	argsSlice[3] = initialSize
	boolCmd := redis.NewBoolCmd(ctx, argsSlice...)
	_ = c.Process(ctx, boolCmd)
	return boolCmd
}

func BfAdd(c mOperation, ctx context.Context, key string, val interface{}) *redis.IntCmd {
	argSlice := make([]interface{}, 3)
	argSlice[0] = "bf.add"
	argSlice[1] = key
	argSlice[2] = val
	intCmd := redis.NewIntCmd(ctx, argSlice...)
	_ = c.Process(ctx, intCmd)
	return intCmd
}

func BfExists(c mOperation, ctx context.Context, key, val interface{}) *redis.BoolCmd {
	argSlice := make([]interface{}, 3)
	argSlice[0] = "bf.exists"
	argSlice[1] = key
	argSlice[2] = val
	boolCmd := redis.NewBoolCmd(ctx, argSlice...)
	_ = c.Process(ctx, boolCmd)
	return boolCmd
}

func BfMAdd(c mOperation, ctx context.Context, args ...interface{}) *redis.IntSliceCmd {
	argSlice := make([]interface{}, 1, len(args)+1)
	argSlice[0] = "bf.madd"
	argSlice = redisAppendArgs(argSlice, args)
	intSliceCmd := redis.NewIntSliceCmd(ctx, argSlice...)
	_ = c.Process(ctx, intSliceCmd)
	return intSliceCmd
}

func BfMExists(c mOperation, ctx context.Context, args ...interface{}) *redis.BoolSliceCmd {
	argSlice := make([]interface{}, 1, len(args)+1)
	argSlice[0] = "bf.MExists"
	argSlice = redisAppendArgs(argSlice, args)
	boolSliceCmd := redis.NewBoolSliceCmd(ctx, argSlice...)
	_ = c.Process(ctx, boolSliceCmd)
	return boolSliceCmd
}

func redisAppendArgs(dst []interface{}, arg ...interface{}) []interface{} {
	if len(arg) == 1 {
		return redisAppendArg(dst, arg[0])
	}
	return append(dst, arg...)
}

func redisAppendArg(dst []interface{}, arg interface{}) []interface{} {
	switch arg := arg.(type) {
	case []string:
		for _, value := range arg {
			dst = append(dst, value)
		}
		return dst
	case []interface{}:
		for _, value := range arg {
			dst = append(dst, value)
		}
		return dst

	case map[string]interface{}:
		for key, value := range arg {
			dst = append(dst, key, value)
		}
		return dst
	case map[string]string:
		for key, value := range arg {
			dst = append(dst, key, value)
		}
		return dst
	default:
		return append(dst, arg)
	}
}

//func BfReserve(c *redis.Client, ctx context.Context) *redis.BoolCmd {
//
//}
