package mredis

import (
	"MFile/generate"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLock_Lock(t *testing.T) {
	key := generate.RandStringN(16)
	fmt.Println(key)
	redisLock := NewRedisLock(context.Background(), RedisDb, key)
	ok, err := redisLock.Lock()
	if !ok {
		fmt.Println("lock redis error err is :", err)
	} else {
		fmt.Println("lock redis successful")
	}

	ok, err = redisLock.UnLock()
	if !ok {
		fmt.Println("UnLock redis error err is :", err)
	} else {
		fmt.Println("UnLock redis successful")
	}
}

func TestLock_Lock2(t *testing.T) {
	var key1 string
	var w sync.WaitGroup
	w.Add(1)
	go func() {
		key1 = generate.RandStringN(16)
		lock := NewRedisLock(context.Background(), RedisDb, key1)
		lock.Lock()
		w.Done()
	}()
	w.Wait()

	lock := NewRedisLock(context.Background(), RedisDb, key1)

	ok, err := lock.Lock()
	if !ok {
		fmt.Println("lock redis error err is :", err)
	} else {
		fmt.Println("lock redis successful")
	}

	ok, err = lock.UnLock()
	if !ok {
		fmt.Println("UnLock redis error err is :", err)
	} else {
		fmt.Println("UnLock redis successful")
	}

}

func TestRedisTime(t *testing.T){
	rs,err:=RedisDb.Set(context.Background(),"time",time.Now().Unix(),0).Result()
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(rs)
}
