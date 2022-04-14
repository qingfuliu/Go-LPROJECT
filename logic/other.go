package logic

import (
	"MFile/db/mredis"
	"fmt"
	"math/rand"
	"time"
)

func Likes(targetName, userName string) error {
	//查看文章点赞时间是否过期
	if rs, err := mredis.RedisDb.TTL(mredis.Ctx, targetName+"_").Result(); err != nil || rs == -2 {
		return fmt.Errorf("文章点赞时间已经过啦")
	}
	//查看自己是否已经点赞
	if rs, err := mredis.BfExists(mredis.RedisDb, mredis.Ctx, targetName+"_liked", userName).Result(); err != nil || rs {
		return fmt.Errorf("你已经赞过")
	}
	//set nx 分布式锁
	if rs, err := mredis.RedisDb.SetNX(mredis.Ctx, targetName+"_"+userName, 1, time.Duration(300+rand.Intn(200))).Result(); err != nil || !rs {
		return fmt.Errorf("点赞太频繁")
	}
	//执行事务
	//1 点赞人数加1 zset hash
	//2.username 加入到布隆过滤器
	tx := mredis.RedisDb.TxPipeline()
	tx.ZIncrBy(mredis.Ctx, "like_count_zset", 1, targetName)
	tx.HIncrBy(mredis.Ctx, "like_count_hash", targetName, 1)
	mredis.BfAdd(tx, mredis.Ctx, targetName+"_liked", userName)
	_, err := tx.Exec(mredis.Ctx)
	return err
}
