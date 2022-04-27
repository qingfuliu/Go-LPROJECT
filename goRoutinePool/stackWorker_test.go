package goRoutinePool

import (
	"testing"
	"time"
)

func TestWorkerStack(t *testing.T) {
	now := time.Now()
	stackWorker := newWorkerStack(10)
	for i := 0; i < 10; i++ {
		err := stackWorker.insert(&goWorker{
			recycleTime: now.Add(time.Duration(int(time.Second) * (i - 5))),
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 10; i++ {
		index := stackWorker.(*workerStack).binarySearch(now.Add(time.Duration(int(time.Second) * (i - 5))))
		if index != i-1 {
			t.Fatalf(" %d is fatal:index should be %d,but is %d", i, i-1, index)
		}
	}

	for i := 0; i < 5; i++ {
		temp := stackWorker.detach()
		if temp == nil {
			t.Fatalf("temp should not be nil")
		}
	}

	for i := 0; i < 10; i++ {
		index := stackWorker.(*workerStack).binarySearch(now.Add(time.Duration(int(time.Second) * (i - 5))))
		if i < 5 && index != i-1 {
			t.Fatalf(" %d is fatal:index should be %d,but is %d", i, i-1, index)
		} else if i >= 5 && index != 4 {
			t.Fatalf(" %d is fatal:index should be %d,but is %d", i, 4, index)
		}
	}

	for i := 0; i < 5; i++ {
		temp := stackWorker.detach()
		if temp == nil {
			t.Fatalf("temp should not be nil")
		}
	}

	for i := 0; i < 5; i++ {
		temp := stackWorker.detach()
		if temp != nil {
			t.Fatalf("temp should  be nil")
		}
	}
}

func BenchmarkWorkerStack(t *testing.B) {
	for i := 0; i < t.N; i++ {
		count := 1000
		now := time.Now()
		stackWorker := newWorkerStack(10)
		for i := 0; i < count; i++ {
			err := stackWorker.insert(&goWorker{
				recycleTime: now.Add(time.Duration(int(time.Second) * (i - 5))),
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		for i := 0; i < count; i++ {
			index := stackWorker.(*workerStack).binarySearch(now.Add(time.Duration(int(time.Second) * (i - 5))))
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

		for i := 0; i < count; i++ {
			index := stackWorker.(*workerStack).binarySearch(now.Add(time.Duration(int(time.Second) * (i - 5))))
			if i < count/2 && index != i-1 {
				t.Fatalf(" %d is fatal:index should be %d,but is %d", i, i-1, index)
			} else if i >= count/2 && index != count/2-1 {
				t.Fatalf(" %d is fatal:index should be %d,but is %d", i, count/2-1, index)
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
}
