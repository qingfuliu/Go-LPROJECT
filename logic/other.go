package logic

import (
	"MFile/db/mredis"
	"MFile/logger"
	"context"
	"fmt"
	"go.uber.org/zap"
)

func Likes(targetName, userName string) error {
	//查看文章点赞时间是否过期
	if rs, err := mredis.RedisDb.TTL(mredis.Ctx, targetName).Result(); err != nil || rs == -2 {
		return fmt.Errorf("文章点赞时间已经过啦")
	}
	bloomKey := targetName + "_liked"
	//查看自己是否已经点赞
	if rs, err := mredis.BfExists(mredis.RedisDb, mredis.Ctx, bloomKey, userName).Result(); err != nil || rs {
		return fmt.Errorf("你已经赞过")
	}
	//redis 分布式锁
	lockKey := targetName + "_" + userName
	redisLock := mredis.NewRedisLock(context.Background(), mredis.RedisDb, lockKey)
	if ok, err := redisLock.Lock(); !ok {
		return err
	}
	defer func(redisLock mredis.RedisLock) {
		_, err := redisLock.UnLock()
		if err != nil {
			logger.MLogger.Error("Service timeout!!", zap.String("targetName", targetName))
		}
	}(redisLock)
	//执行事务
	//1 点赞人数加1 zset hash
	//2.username 加入到布隆过滤器

	reserveKey := "like_count_"
	tx := mredis.RedisDb.TxPipeline()
	tx.ZIncrBy(mredis.Ctx, reserveKey+"zset", 1, targetName)
	tx.HIncrBy(mredis.Ctx, reserveKey+"hash", targetName, 1)
	mredis.BfAdd(tx, mredis.Ctx, bloomKey, userName)
	_, err := tx.Exec(mredis.Ctx)
	return err
}
