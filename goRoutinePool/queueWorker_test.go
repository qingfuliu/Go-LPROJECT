package goRoutinePool

import (
	"fmt"
	"testing"
	"time"
)

func TestWorkerQueue(t *testing.T) {
	now := time.Now()
	count := 1000
	stackWorker := newWorkerQueue(count)
	for i := 0; i < count; i++ {
		err := stackWorker.insert(&goWorker{
			recycleTime: now.Add(time.Duration(int(time.Second) * (i - 5))),
		})
		if err != nil {
			t.Fatalf("%d is err:%v", i, err)
		}
	}

	for i := 0; i < count; i++ {
		index := stackWorker.(*workerQueue).binarySearch(now.Add(time.Duration(int(time.Second) * (i - 5))))
		if index != i-1 {
			t.Fatalf(" %d is fatal:index should be %d,but is %d", i, i-1, index)
		}
	}

	for i := 0; i < count/2; i++ {
		temp := stackWorker.detach()
		if temp == nil {
			t.Fatalf("temp should not be nil")
		}
	}

	for i := 0; i < count/2; i++ {
		index := stackWorker.(*workerQueue).binarySearch(now.Add(time.Duration(int(time.Second) * (i - 5))))
		if i < count/2 && index != -1 {
			t.Fatalf(" %d is fatal:index should be %d,but is %d", i, -1, index)
		} else if i >= count/2 && index != i-1 {
			t.Fatalf(" %d is fatal:index should be %d,but is %d", i, 4, index)
		}
	}

	for i := 0; i < count/2; i++ {
		temp := stackWorker.detach()
		if temp == nil {
			t.Fatalf("temp should not be nil")
		}
	}
	for i := 0; i < count/2; i++ {
		temp := stackWorker.detach()
		if temp != nil {
			t.Fatalf("temp should  be nil")
		}
	}
}

func TestWorkerQueue2(t *testing.T) {
	now := time.Now()
	count := 1000
	stackWorker := newWorkerQueue(count)
	for i := 0; i < 1000; i++ {

		err := stackWorker.insert(&goWorker{
			recycleTime: now.Add(time.Duration(int(time.Second) * (i - 1000))),
		})
		if err != nil {
			t.Fatalf("%d is err:%v", i, err)
		}
	}
	for i := 0; i < 800; i++ {
		temp := stackWorker.detach()
		if temp == nil {
			t.Fatalf("temp should not be nil")
		}
	}

	if stackWorker.(*workerQueue).isFull {
		t.Fatalf("should not be full")
	}

	for i := 0; i < 800; i++ {
		err := stackWorker.insert(&goWorker{
			recycleTime: now.Add(time.Duration(int(time.Second) * (i))),
		})
		if err != nil {
			t.Fatalf("%d is err:%v", i, err)
		}
	}
	if !stackWorker.(*workerQueue).isFull {
		t.Fatalf("should be full")
	}
	if stackWorker.(*workerQueue).tail != 800 {
		t.Fatalf("should be full")
	}

	for i := 1; i < count; i++ {
		if i < 800 {
			index := stackWorker.(*workerQueue).binarySearch(now.Add(time.Duration(int(time.Second) * (i))))
			if index != i-1 {
				t.Fatalf(" %d is fatal:index should be %d,but is %d", i, i-1, index)
			}
		} else {
			index := stackWorker.(*workerQueue).binarySearch(now.Add(time.Duration(int(time.Second) * (i - 1000))))
			if index != i-1 && (i == 800 && index != -1) {
				fmt.Println(now.Add(time.Duration(int(time.Second) * (i - 1000))))
				fmt.Println(stackWorker.(*workerQueue).items[stackWorker.(*workerQueue).head].recycleTime)
				t.Fatalf(" %d is fatal:index should be %d,but is %d", i, i-1, index)
			}
		}
	}

	for i := 0; i < count/2; i++ {
		temp := stackWorker.detach()
		if temp == nil {
			t.Fatalf("temp should not be nil")
		}
	}

	for i := 0; i < count/2; i++ {
		temp := stackWorker.detach()
		if temp == nil {
			t.Fatalf("temp should not be nil")
		}
	}
	for i := 0; i < count/2; i++ {
		temp := stackWorker.detach()
		if temp != nil {
			t.Fatalf("temp should  be nil")
		}
	}
}
