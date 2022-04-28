package goRoutinePool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewPoolFuncBinarySearch(t *testing.T) {
	pool, err := NewPoolFunc(0, nil)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	for i := 0; i < 100; i++ {
		pool.workers = append(pool.workers, &goWorkerWithFunc{
			recycleTime: now.Add(time.Duration(int(time.Second) * i)),
		})
	}
	fmt.Println("len pool.workers:  ", len(pool.workers))
	for i := 0; i < 100; i++ {
		index := pool.binarySearch(now.Add(time.Duration(int(time.Second) * i)))
		if index != i-1 {
			t.Fatalf("%d is fatal,index should be %d,but is %d", i, i-1, index)
		}
	}
}

func TestNewPoolFunc(t *testing.T) {
	//runtime.BlockProfile()
	var i int32
	count := 1000
	interval := 1000
	var wait sync.WaitGroup
	wait.Add(count)
	pool, err := NewPoolFunc(0, func(arg interface{}) {
		wait.Done()
		p := arg.(int)
		atomic.AddInt32(&i, int32(p))
	})
	if err != nil {
		t.Fatal(err)
	}
	for p := 0; p < count; p++ {
		err = pool.Submit(interval)
		if err != nil {
			t.Fatal(err)
		}
	}
	wait.Wait()
	if int(i) != count*interval {
		t.Fatalf("i should be %d,but is %d", count*interval, i)
	}
	pool.Release()
	if pool.blockingNum != 0 || pool.Running() != 0 || len(pool.workers) != 0 {
		t.Fatal("error release", pool.blockingNum, pool.Running(), len(pool.workers))
	}
}
