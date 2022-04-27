package main

import (
	"MFile/other"
	"context"
	"fmt"
	_ "net/http/pprof"
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
	err = timeWheel_.AddTask("lqf12", "hhh", time.Second*1, true)
	err = timeWheel_.AddTask("lqf123", "hhh", time.Second*1, true)
	err = timeWheel_.AddTask("lqf124", "hhh", time.Second*1, true)
	time.Sleep(time.Second * 35)
	timeWheel_.Stop()
	err = timeWheel_.AddTask("lqf", "hhh", time.Second*2, true)
	if err == nil {
		panic("add task should be err!!")
	}
}

func test() {
	a := ""
	for index := 0; index <= 1000; index++ {
		a += "aaaaaaaaaaaaaaaa"
	}
	fmt.Println(a)
}
func main() {
	//flag.Parse()
	//file, err := os.Create(commend.CpuProfilePath)
	//file2, err2 := os.Create(commend.HeapProfilePath)
	//if err2 != nil || err != nil {
	//	os.Exit(1)
	//}
	//defer file.Close()
	//defer file2.Close()
	//err = pprof.StartCPUProfile(file)
	//if err != nil {
	//	panic(err)
	//}
	//defer pprof.StopCPUProfile()
	//
	//err = pprof.WriteHeapProfile(file2)
	//if err != nil {
	//	panic(err)
	//}
	//for index := 0; index <= 10; index++ {
	//	test()
	//}
	//testNewTimeWheel()
	//router.Engine.Run("localhost:8081")
}
