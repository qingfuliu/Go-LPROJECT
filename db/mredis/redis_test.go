package mredis

import (
	"MFile/logger"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestValid(t *testing.T) {
	rs, err := RedisDb.LPush(context.Background(), "test list", "first", "second").Result()
	if err != nil {
		logger.MLogger.Error("redis sentinel connect faild", zap.Error(err))
	} else {
		fmt.Println("result: ", rs)
	}

	rs1, err1 := RedisDb.LRange(context.Background(), "test list", 0, 1).Result()
	if err1 != nil {
		logger.MLogger.Error("redis sentinel connect faild", zap.Error(err1))
	} else {
		fmt.Println("result: ", rs1)
	}

}

func TestSet(t *testing.T) {
	if setAdd, err := RedisDb.Set(context.Background(), "女权傻吊_", 1, time.Hour*24*7).Result(); err == nil {
		fmt.Printf("set Add result is: %v", setAdd)
	} else {
		t.Fatal("setAdd Failed")
	}

	//if setMembers, err := RedisDb.SMembers(context.Background(), "lqf").Result(); err == nil {
	//	fmt.Printf("set Add result is: %v", setMembers)
	//} else {
	//	t.Fatal("setMembers Failed")
	//}
	//
	//if SMIsMember, err := RedisDb.SMIsMember(context.Background(), "lqf", "first", "second", "sd").Result(); err == nil {
	//	fmt.Printf("set Add result is: %v", SMIsMember)
	//} else {
	//	t.Fatal("SMIsMember Failed")
	//}
	//
	//if SIsMember, err := RedisDb.SIsMember(context.Background(), "lqf", "first").Result(); err == nil {
	//	fmt.Printf("set Add result is: %v", SIsMember)
	//} else {
	//	t.Fatal("SIsMember Failed")
	//}
	//
	//if SPop, err := RedisDb.SPop(context.Background(), "lqf").Result(); err == nil {
	//	fmt.Printf("set Add result is: %v", SPop)
	//} else {
	//	t.Fatal("SPop Failed")
	//}
	//
	//if SRandMember, err := RedisDb.SRandMember(context.Background(), "lqf").Result(); err == nil {
	//	fmt.Printf("set Add result is: %v", SRandMember)
	//} else {
	//	t.Fatal("SRandMember Failed")
	//}
	//
	//if SPopN, err := RedisDb.SPopN(context.Background(), "lqf", 1).Result(); err == nil {
	//	fmt.Printf("set Add result is: %v", SPopN)
	//} else {
	//	t.Fatal("SPopN Failed")
	//}
	//
	//if SCard, err := RedisDb.SCard(context.Background(), "lqf").Result(); err == nil {
	//	fmt.Printf("set Add result is: %v", SCard)
	//} else {
	//	t.Fatal("SCard Failed")
	//}
	//
	//if setAdd, err := RedisDb.SAdd(context.Background(), "lqf2", "first1", "second1", "thread1").Result(); err == nil {
	//	fmt.Printf("set Add result is: %v", setAdd)
	//} else {
	//	t.Fatal("setAdd Failed")
	//}
	//
	//if SMove, err := RedisDb.SMove(context.Background(), "lqf", "lqf2", "second").Result(); err == nil {
	//	fmt.Printf("SMove  result is: %v", SMove)
	//} else {
	//	t.Fatal("SMove Failed")
	//}
	//fmt.Println()
	//if setAdd, err := RedisDb.SInter(context.Background(), "lqf", "lqf2").Result(); err == nil {
	//	fmt.Println("set Add result is:", setAdd)
	//} else {
	//	t.Fatal("setAdd Failed")
	//}
	//
	//if setAdd, err := RedisDb.SDiff(context.Background(), "lqf", "lqf2").Result(); err == nil {
	//	fmt.Println("SDiff SDiff result is:", setAdd)
	//} else {
	//	t.Fatal("SDiff Failed")
	//}
	//
	//if setAdd, err := RedisDb.SUnion(context.Background(), "lqf", "lqf2").Result(); err == nil {
	//	fmt.Println("SDiff SDiff result is:", setAdd)
	//} else {
	//	t.Fatal("SDiff Failed")
	//}

}

func TestZsetAdd(t *testing.T) {
	RedisDb.ZAdd(context.Background(), "zsetTest", &redis.Z{4, "lqfs4"}, &redis.Z{4, "lqfs41"}, &redis.Z{4, "lqfs42"}, &redis.Z{4, "lqfs43"}, &redis.Z{4, "lqfs44"})
	RedisDb.ZAdd(context.Background(), "zsetTest", &redis.Z{3, "lqfs3"}, &redis.Z{3, "lqfs31"}, &redis.Z{3, "lqfs32"}, &redis.Z{3, "lqfs33"}, &redis.Z{3, "lqfs34"})
	RedisDb.ZAdd(context.Background(), "zsetTest", &redis.Z{2, "lqfs2"}, &redis.Z{2, "lqfs21"}, &redis.Z{2, "lqfs22"}, &redis.Z{2, "lqfs22"}, &redis.Z{2, "lqfs24"})
	RedisDb.ZAdd(context.Background(), "zsetTest", &redis.Z{1, "lqf"}, &redis.Z{1, "lqf1"}, &redis.Z{1, "lqf11"}, &redis.Z{1, "lqf12"}, &redis.Z{1, "lqf13"}, &redis.Z{1, "lqf14"})
}

func TestZsetRange(t *testing.T) {
	rs, err := RedisDb.ZRangeWithScores(context.Background(), "zsetTest", 0, 10).Result()
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(rs)
	}
}
func TestZsetRangeByScore(t *testing.T) {
	rs, err := RedisDb.ZRangeByScore(context.Background(), "zsetTest", &redis.ZRangeBy{"1", "2", 0, 0}).Result()
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(rs)
	}
}

func TestZCard(t *testing.T) {
	rs, err := RedisDb.ZCard(context.Background(), "zsetTest").Result()
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(rs)
	}
}

func TestZRangeByAges(t *testing.T) {

}

func TestZRank(t *testing.T) {
	rs, err := RedisDb.ZRank(context.Background(), "zsetTest", "lqf1").Result()
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(rs)
	}
}

func TestMulti(t *testing.T) {

	transcation := RedisDb.TxPipeline()
	transcation.ZCard(context.Background(), "zsetTest")
	transcation.ZAdd(context.Background(), "zsetTest", &redis.Z{10, "key"})
	rs, err := transcation.Exec(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range rs {
		fmt.Println(value.String())
	}
}

func TestBfReserve(t *testing.T) {
	cmd1, err1 := BfReserve(RedisDb, context.Background(), "lqf2221", 0.0001, 10000).Result()
	if err1 != nil {
		t.Fatal(err1)
	}
	fmt.Println("cmd1:", cmd1)
	cmd, err := BfMExists(RedisDb, context.Background(), "lqfff", 55, "666", "777", "lqfwd").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cmd)
}
func TestWatch(t *testing.T) {
	key := "lqfhhh"
	err := RedisDb.Watch(context.Background(), func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(context.Background(), func(pipeliner redis.Pipeliner) error {
			_, err := pipeliner.Set(context.Background(), "lqfhhh", "lqfwatch", 0).Result()
			return err
		})
		return err
	}, key)
	if err != nil {
		t.Fatal(err)
	}
}
func TestLikes(t *testing.T) {
	targetName := "女权杀掉"
	userName := "all"
	tx := RedisDb.TxPipeline()
	tx.ZIncrBy(Ctx, "like_count_zset", 1, targetName)
	tx.HIncrBy(Ctx, "like_count_hash", targetName, 1)
	BfAdd(tx, Ctx, targetName+"_", userName)
	_, err := tx.Exec(Ctx)
	if err != nil {
		t.Fatal(err)
	}
}
