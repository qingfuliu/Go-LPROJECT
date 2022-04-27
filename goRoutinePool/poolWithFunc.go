package goRoutinePool

import (
	"sync"
	"sync/atomic"
	"time"
)

type PoolFunc struct {
	*Options
	capacity    int32
	blockingNum int32
	running     int32
	state       int32
	workers     []*goWorkerWithFunc
	workCache   sync.Pool
	goFunc      func(interface{})
	mutex       sync.Locker
	cond        sync.Cond
}

func (p *PoolFunc) binarySearch(time_ time.Time) int {
	if len(p.workers) == 0 || p.workers[0].recycleTime.After(time_) || p.workers[0].recycleTime.Equal(time_) {
		return -1
	}
	l := 0
	r := len(p.workers)
	for l < r {
		mid := l + (l-r)>>1
		if p.workers[mid].recycleTime.Before(time_) {
			l = mid + 1
		} else {
			r = mid
		}
	}
	return r - 1
}
func (p *PoolFunc) regularCleaner() {
	trickier := time.NewTicker(p.ExpireDuration)
	defer trickier.Stop()
	var expires []*goWorkerWithFunc
	for range trickier.C {
		if p.IsClosed() {
			return
		}
		index := p.binarySearch(time.Now().Add(p.ExpireDuration))
		if index == -1 {
			continue
		}
		expires = expires[:0]
		expires = append(expires, p.workers[:index+1]...)
		m := copy(p.workers, p.workers[index+1:])
		for i := m; i < len(p.workers); i++ {
			p.workers[i] = nil
		}

		p.workers = p.workers[:m]
		for i := range expires {
			expires[i].arg <- nil
			expires[i] = nil
		}

	}
}

func (p *PoolFunc) incRunning() {
	atomic.AddInt32(&p.running, 1)
}

func (p *PoolFunc) decrRunning() {
	atomic.AddInt32(&p.running, -1)
}

func (p *PoolFunc) Running() int32 {
	return atomic.LoadInt32(&p.running)
}

func (p *PoolFunc) Cap() int32 {
	return atomic.LoadInt32(&p.capacity)
}

func (p *PoolFunc) Tune(size int32) {
	if p.IsClosed() {
		return
	}
	if p.capacity == -1 || p.PreAllocate || size <= 0 || size == p.capacity {
		return
	}
	oldCapacity := p.capacity
	atomic.StoreInt32(&p.capacity, size)
	if size > oldCapacity {
		if size == oldCapacity+1 {
			p.cond.Signal()
		} else {
			p.cond.Broadcast()
		}
	}
	return
}

func (p *PoolFunc) Release() {
	if atomic.CompareAndSwapInt32(&p.state, poolStateRunning, poolStateClosed) {
		p.mutex.Lock()
		for i := range p.workers {
			p.workers[i].arg <- nil
		}
		p.workers = p.workers[:0]
		p.running = 0
		p.cond.Broadcast()
		p.mutex.Unlock()
	}
}

func (p *PoolFunc) IsClosed() bool {
	return atomic.LoadInt32(&p.state) == poolStateClosed
}

func (p *PoolFunc) retrieveWorker() (g *goWorkerWithFunc) {

	newWorkerHandle := func() {
		g = p.workCache.Get().(*goWorkerWithFunc)
		g.run()
	}

	p.mutex.Lock()
	if p.IsClosed() {
		p.mutex.Unlock()
		return
	}
	n := len(p.workers)
	if n > 0 {
		g = p.workers[0]
		p.workers = p.workers[1:]
		p.mutex.Unlock()
		return
	} else {
		capacity := p.Cap()
		if capacity == -1 || capacity > p.Running() {
			p.mutex.Unlock()
			newWorkerHandle()
			return
		} else if p.NoBlocking || p.MaxBlockNums <= p.blockingNum {
			p.mutex.Unlock()
			return
		}
	retry:
		p.blockingNum++
		p.cond.Wait()
		p.blockingNum--
		if p.Running() == 0 {
			if p.IsClosed() {
				p.mutex.Unlock()
				return
			}
			p.mutex.Unlock()
			newWorkerHandle()
			return
		}
		n = len(p.workers)
		if n > 0 {
			g = p.workers[0]
			p.workers = p.workers[1:]
			p.mutex.Unlock()
			return
		} else {
			goto retry
		}
	}
}

func (p *PoolFunc) reserveWorkers(g *goWorkerWithFunc) bool {
	if p.IsClosed() || (p.Cap() != -1 && p.Cap() <= p.Running()) {
		p.cond.Broadcast()
		return false
	}
	g.recycleTime = time.Now()
	p.mutex.Lock()
	if p.IsClosed() {
		return false
	}
	p.workers = append(p.workers, g)
	p.cond.Signal()
	p.mutex.Unlock()
	return true
}
