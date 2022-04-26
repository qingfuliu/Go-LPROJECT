package other

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewTimeWheel(t *testing.T) {
	timeWheel_ := NewTimeWheel(time.Second, func(key, value interface{}) {
		fmt.Println(key, value)
	}, 10, context.Background())
	timeWheel_.Run()
	err := timeWheel_.AddTask("lqf", "hhh", time.Second*2, true)
	if err != nil {
		t.Fatal(err)
	} else {
		time.Sleep(time.Second * 2)
		if val, ok := timeWheel_.(*timeWheel).tasksMap.Load("lqf"); !ok {
			t.Fatal("map add error!!!")
		} else {
			fmt.Println("value of the lqf: ", val)
		}
	}
	time.Sleep(time.Second * 10)
	timeWheel_.Stop()
	err = timeWheel_.AddTask("lqf", "hhh", time.Second*2, true)
	if err == nil {
		t.Fatal("add task should be err!!")
	}
}
