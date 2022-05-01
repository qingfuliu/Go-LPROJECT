package mredis

import (
	"MFile/db/mredis"
	"MFile/generate"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestRedisCond_Acquire(t *testing.T) {
	redisCond := NewRedisCond(mredis.RedisDb, strconv.FormatInt(generate.GetSnakeId(), 10), time.Second*10, "TestRedisCond_Acquire4")
	if ok, err := redisCond.Acquire(); !ok {
		t.Fatal(err)
	} else {
		fmt.Println("success acquire cond")
	}
	if ok, err := redisCond.Release(); !ok || err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("success release cond")
	}

	if ok, err := redisCond.Release(); ok {
		t.Fatal(err)
	} else {
		fmt.Println("success refresh cond")
	}
}
