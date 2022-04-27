package goRoutinePool

import (
	"time"
)

type workerQueue struct {
	items   []*goWorker
	expires []*goWorker
	head    int
	tail    int
	size    int
	isFull  bool
}

func newWorkerQueue(size int) workArray {
	return &workerQueue{
		items: make([]*goWorker, size),
		size:  size,
	}
}

func (wq *workerQueue) insert(g *goWorker) error {
	if wq.isFull {
		return errorQueueIsFull
	}
	wq.items[wq.tail] = g
	wq.tail++
	if wq.tail == wq.size {
		wq.tail = 0
	}
	if wq.head == wq.tail {
		wq.isFull = true
	}
	return nil
}
func (wq *workerQueue) detach() (g *goWorker) {
	if wq.isEmpty() {
		return
	}
	g = wq.items[wq.head]
	wq.head++
	if wq.head == wq.size {
		wq.head = 0
	}
	wq.isFull = false
	return
}
func (wq *workerQueue) len() int {
	if wq.head == wq.tail {
		if wq.isFull {
			return wq.size
		}
		return 0
	}
	if wq.head < wq.tail {
		return wq.tail - wq.head
	}
	return wq.size - wq.head + wq.tail
}
func (wq *workerQueue) isEmpty() bool {
	return wq.head == wq.tail && !wq.isFull
}
func (wq *workerQueue) retrieveExpire(duration time.Duration) []*goWorker {
	index := wq.binarySearch(time.Now().Add(-duration))
	if index == -1 {
		return nil
	}

	wq.expires = wq.expires[:0]

	if wq.head < wq.tail {
		wq.expires = append(wq.expires, wq.items[wq.head:index+1]...)
		for i := wq.head; i <= index; i++ {
			wq.items[i] = nil
		}
	} else {
		wq.expires = append(wq.expires, wq.items[wq.head:]...)
		for i := wq.head; i < len(wq.items); i++ {
			wq.items[i] = nil
		}
		wq.expires = append(wq.expires, wq.items[:index+1]...)
		for i := 0; i <= index; i++ {
			wq.items[i] = nil
		}
	}
	wq.head = index + 1
	if wq.head == wq.size {
		wq.head = 0
	}
	if len(wq.expires) > 0 {
		wq.isFull = false
	}
	return wq.expires
}

func (wq *workerQueue) reset() {
	if wq.isEmpty() {
		return
	}
	if wq.head < wq.tail {
		for i := wq.head; i < wq.tail; i++ {
			wq.items[i].task <- nil
			wq.items[i] = nil
		}
	} else {
		for i := wq.head; i < wq.size; i++ {
			wq.items[i].task <- nil
			wq.items[i] = nil
		}
		for i := 0; i < wq.tail; i++ {
			wq.items[i].task <- nil
			wq.items[i] = nil
		}
	}
	wq.head, wq.tail, wq.size = 0, 0, 0
	wq.isFull = false
	wq.items = wq.items[:0]
}
func (wq *workerQueue) binarySearch(expireTime time.Time) int {
	if wq.isEmpty() || expireTime.Before(wq.items[wq.head].recycleTime) || expireTime.Equal(wq.items[wq.head].recycleTime) {
		return -1
	}
	length := wq.size
	l := 0
	r := (wq.tail - 1 - wq.head + length) % length
	basal := wq.head
	for l <= r {
		mid := l + (r-l)>>1
		realMid := (mid + basal) % length
		if wq.items[realMid].recycleTime.Before(expireTime) {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return (r + basal) % length
}
