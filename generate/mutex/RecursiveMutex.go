package mutex

import (
	"github.com/petermattis/goid"
	"sync"
	"sync/atomic"
)

func NewRecursiveMutex()RecursiveMutex{
	return &recursiveMutex{
		gid:-1,
		countRecursive: 0,
	}
}

type recursiveMutex struct {
	mutex          sync.Mutex
	countRecursive uint64
	gid            int64
}

type RecursiveMutex interface {
	Lock()
	UnLock()
}

func (m *recursiveMutex) Lock() {
	gid := goid.Get()

	if atomic.LoadInt64(&m.gid) == gid {
		m.countRecursive++
		return
	}
	m.mutex.Lock()
	atomic.StoreInt64(&m.gid, gid)
	m.countRecursive = 1
}

func (m *recursiveMutex) UnLock() {
	gid := goid.Get()
	if gid != atomic.LoadInt64(&m.gid) {
		panic("mutex unlock fatal!")
	}
	m.countRecursive--
	if m.countRecursive!=0{
		return
	}
	atomic.StoreInt64(&m.gid,-1)
	m.mutex.Unlock()
}
