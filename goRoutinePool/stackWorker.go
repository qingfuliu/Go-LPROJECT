package goRoutinePool

import (
	"time"
)

type workerStack struct {
	items   []*goWorker
	expires []*goWorker
	size    int
}

func newWorkerStack(size int) workArray {
	return &workerStack{
		size:  size,
		items: make([]*goWorker, 0, size),
	}
}

func (s *workerStack) insert(g *goWorker) error {
	s.items = append(s.items, g)
	return nil
}
func (s *workerStack) detach() (g *goWorker) {
	if len(s.items) == 0 {
		return
	}
	g = s.items[len(s.items)-1]
	s.items[len(s.items)-1] = nil
	s.items = s.items[:len(s.items)-1]
	return
}
func (s *workerStack) len() int {
	return len(s.items)
}
func (s *workerStack) isEmpty() bool {
	return len(s.items) == 0
}
func (s *workerStack) retrieveExpire(expireTime time.Duration) []*goWorker {
	length := len(s.items)
	if length == 0 {
		return nil
	}
	index := s.binarySearch(time.Now().Add(-expireTime))
	if index == -1 {
		return nil
	}
	s.expires = s.expires[:0]
	s.expires = append(s.expires, s.items[:index+1]...)
	m := copy(s.items, s.items[index+1:])
	for i := m; i < length; i++ {
		s.items[i] = nil
	}
	s.items = s.items[:m]
	return s.expires
}
func (s *workerStack) reset() {
	for i := 0; i < len(s.items); i++ {
		s.items[i].task <- nil
		s.items[i] = nil
	}
	s.items = s.items[:0]
}
func (s *workerStack) binarySearch(expireTime time.Time) int {

	if s.isEmpty() || s.items[0].recycleTime.After(expireTime) || s.items[0].recycleTime.Equal(expireTime) {
		return -1
	}
	l := 0
	r := len(s.items)
	for l < r {
		mid := l + (r-l)>>1
		if s.items[mid].recycleTime.Before(expireTime) {
			l = mid + 1
		} else {
			r = mid
		}
	}
	return r - 1
}
