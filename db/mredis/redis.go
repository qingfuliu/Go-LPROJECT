package mredis

import (
	"MFile/logger"
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"os"
	"time"
)

var RedisDb *redis.Client
var Ctx = context.Background()

func singleNodeInit() {
	RedisDb = redis.NewClient(
		&redis.Options{
			DialTimeout:  2 * time.Minute,
			PoolTimeout:  30 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			Network:      "tcp",
			Addr:         "192.168.1.103:6379",
			Username:     "lqf",
			Password:     "",
			DB:           0,
			PoolSize:     20,
			MinIdleConns: 10,
		})
}

func sentinelInit() {
	RedisDb = redis.NewFailoverClient(
		&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{"192.168.1.103:26379"},
			DB:            0,
			Password:      "",
			DialTimeout:   2 * time.Minute,
			PoolTimeout:   30 * time.Second,
			ReadTimeout:   30 * time.Second,
			WriteTimeout:  30 * time.Second,
			PoolSize:      20,
			MinIdleConns:  10,
		},
	)

}

func clusterInit() {
	//clusterClient := redis.NewClusterClient(&redis.ClusterOptions{})
}

func init() {
	singleNodeInit()
	if err := RedisDb.Ping(context.Background()).Err(); err != nil {
		logger.MLogger.Fatal("can not connect to mredis server, err is:", zap.Error(err))
		os.Exit(1)
	}
	logger.MLogger.Info("successfully connect to mredis server")
}
