package goRoutinePool

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	pool, err := NewPool(1000, WithExpireDuration(time.Second), WithPreAllocate(true), WithMaxBlockNums(100))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v", pool)
	defer pool.Release()
	var i int64 = 0
	wait := sync.WaitGroup{}
	wait.Add(10000)

	testFunc := func() {
		for k := 0; k < 1000; k++ {
			atomic.AddInt64(&i, 1)
		}
		fmt.Println("num goroutine: ", runtime.NumGoroutine())
		wait.Done()
	}
	for j := 0; j < 10000; j++ {
		err = pool.Submit(testFunc)
		if err != nil {
			t.Fatalf("%d error:%v", i, err)
		}
	}
	wait.Wait()
	if i != 1000*10000 {
		t.Fatalf("i should be 1000,but is: %d", i)
	}

	if pool.blockingNum != 0 {
		t.Fatalf("blockingNum should be 0,but is: %d", pool.blockingNum)
	}
}

func TestNewPoolParallel(t *testing.T) {
	pool, err := NewPool(1, WithExpireDuration(time.Second), WithPreAllocate(true), WithMaxBlockNums(100))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v", pool)
	defer pool.Release()
	var i int64 = 0
	wait := sync.WaitGroup{}
	wait.Add(1000)

	testFunc := func() {
		atomic.AddInt64(&i, 1)
		wait.Done()
	}
	for j := 0; j < 1000; j++ {
		err = pool.Submit(testFunc)
		if err != nil {
			t.Fatalf("%d error:%v", i, err)
		}
	}
	wait.Wait()
	if i != 1000 {
		t.Fatalf("i should be 1000,but is: %d", i)
	}
}
