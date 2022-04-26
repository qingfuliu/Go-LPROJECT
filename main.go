package main

import (
	"MFile/other"
	"MFile/router"
	"context"
	"fmt"
	"time"
)

func testNewTimeWheel() {
	timeWheel_ := other.NewTimeWheel(time.Second, func(key, value interface{}) {
		fmt.Println(key, value)
	}, 10, context.Background())
	timeWheel_.Run()
	err := timeWheel_.AddTask("lqf1", "hhh", time.Second*1, true)
	if err != nil {
		panic(err)
	}

	err = timeWheel_.AddTask("lqf5", "hhh", time.Second*15, true)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 35)
	timeWheel_.Stop()
	err = timeWheel_.AddTask("lqf", "hhh", time.Second*2, true)
	if err == nil {
		panic("add task should be err!!")
	}
}
func main() {
	testNewTimeWheel()
	router.Engine.Run("localhost:8081")
}
