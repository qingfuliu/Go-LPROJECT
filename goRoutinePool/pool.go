package goRoutinePool

import (
	"MFile/other"
	"hash/crc32"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	*Options
	cond        *sync.Cond
	mutex       sync.Locker
	state       int32
	capacity    int32
	running     int32
	blockingNum int32
	workerCache sync.Pool
	works       workArray
}

func (p *Pool) incrRunning() {
	atomic.AddInt32(&p.running, 1)
}

func (p *Pool) decrRunning() {
	atomic.AddInt32(&p.running, -1)
}

func (p *Pool) regularCleaner() {
	trickier := time.NewTicker(p.ExpireDuration)
	defer trickier.Stop()
	for range trickier.C {
		if p.IsClosed() {
			return
		}
		p.mutex.Lock()
		expires := p.works.retrieveExpire(p.ExpireDuration)
		p.mutex.Unlock()
		for index := range expires {
			expires[index].task <- nil
			expires[index] = nil
		}
		if p.Running() == 0 {
			p.cond.Broadcast()
		}
	}
}

//
// NewPool
//  @Description: return a new pool
//  @param capacity
//  @param options
//  @return pool
//
func NewPool(capacity int32, options ...Option) (pool *Pool, err error) {
	opts := loadOption(options...)
	if opts.PreAllocate {
		if capacity == -1 {
			return nil, errorInvalidCapacity
		} else if capacity == 0 {
			capacity = DefaultAntsPoolSize
		}
	} else {
		capacity = -1
	}

	if opts.ExpireDuration < 0 {
		return nil, errorInvalidExpire
	} else if opts.ExpireDuration == 0 {
		opts.ExpireDuration = DefaultCleanIntervalTime
	}

	mutex := other.NewSpinLock()
	pool = &Pool{
		Options:  opts,
		mutex:    mutex,
		capacity: capacity,
	}
	pool.cond = sync.NewCond(mutex)
	pool.workerCache = sync.Pool{
		New: func() interface{} {
			return &goWorker{
				pool: pool,
				task: make(chan func()),
			}
		},
	}
	if pool.PreAllocate {
		pool.works = newWorkArrayWithType(queueWorkArray, int(capacity))
	} else {
		pool.works = newWorkArrayWithType(stackWorkArray, int(crc32.Size))
	}
	go pool.regularCleaner()
	return
}

//
// Submit
//  @Description: submit to run
//  @receiver p
//  @param f
//
func (p *Pool) Submit(f func()) error {
	if p.IsClosed() {
		return errorPoolClosed
	}
	if w := p.retrieveWorker(); w == nil {
		return errorPoolOverLoad
	} else {
		w.task <- f
	}
	return nil
}

func (p *Pool) AsyncSubmit() error {
	return nil
}

//
// Release
//  @Description: close the pool
//  @receiver p
//
func (p *Pool) Release() {

	if atomic.CompareAndSwapInt32(&p.state, poolStateRunning, poolStateClosed) {
		p.mutex.Lock()
		p.running = 0
		p.works.reset()
		p.mutex.Unlock()
		p.cond.Broadcast()
	}
}

//
// Reboot
//  @Description: restart the pool,If it's closed
//  @receiver p
//
func (p *Pool) Reboot() {
	if atomic.CompareAndSwapInt32(&p.state, poolStateClosed, poolStateRunning) {
		go p.regularCleaner()
	}
}

//
// Running
//  @Description: atomic return the nums of the running goroutine
//  @receiver p
//  @return int32
//
func (p *Pool) Running() int32 {
	return atomic.LoadInt32(&p.running)
}

//
// BlockingNum
//  @Description:atomic return blockingNum
//  @receiver p
//
func (p *Pool) BlockingNum() int32 {
	return atomic.LoadInt32(&p.blockingNum)
}

//
// Cap
//  @Description: atomic return the capacity
//  @receiver p
//  @return int32
//
func (p *Pool) Cap() int32 {
	return atomic.LoadInt32(&p.capacity)
}

//
// Free
//  @Description: atomic return the Free num
//  @receiver p
//  @return int32
//
func (p *Pool) Free() int32 {
	capacity := p.Cap()
	if capacity == -1 {
		return -1
	}
	return capacity - p.Running()
}

//
// Tune
//  @Description:resize the worker array with the special size
//  @receiver p
//  @param size
//
func (p *Pool) Tune(size int32) {
	p.mutex.Lock()
	oldSize := p.capacity
	if p.capacity == -1 || size == p.capacity || size <= 0 || p.PreAllocate {
		return
	}
	p.capacity = size
	if oldSize < size {
		if size == oldSize+1 {
			p.cond.Signal()
		} else {
			p.cond.Broadcast()
		}
	}
	p.mutex.Unlock()
}

func (p *Pool) IsClosed() bool {
	return atomic.LoadInt32(&p.state) == poolStateClosed
}

//
// retrieveWorker
//  @Description: return a available worker from pool,If there's none,make decisions according to noBlocking and MaxBlockingNum
//  @receiver p
//  @return *goWorker
//
func (p *Pool) retrieveWorker() (g *goWorker) {
	newWorkerHandle := func() {
		g = p.workerCache.Get().(*goWorker)
		g.run()
	}
	p.mutex.Lock()

	if p.IsClosed() {
		p.mutex.Unlock()
		return
	}

	if g = p.works.detach(); g != nil {
		p.mutex.Unlock()
		return
	} else if capacity := p.Cap(); capacity == -1 || p.Running() < capacity {
		p.mutex.Unlock()
		newWorkerHandle()
		return
	} else {
		if p.NoBlocking || p.blockingNum == p.MaxBlockNums {
			p.mutex.Unlock()
			return
		}
	retry:
		p.blockingNum++
		p.cond.Wait()
		p.blockingNum--
		if p.Running() == 0 {
			p.mutex.Unlock()
			if !p.IsClosed() {
				newWorkerHandle()
			}
			return
		}

		if g = p.works.detach(); g == nil {
			if p.Cap() < p.Running() {
				p.mutex.Unlock()
				newWorkerHandle()
				return
			}
			goto retry
		}
		p.mutex.Unlock()
	}
	return
}

//
// reserveWorkers
//  @Description: return the worker to the pool,and notify other blocked goroutine to worker continue
//  @receiver p
//  @param w
//  @return bool
//
func (p *Pool) reserveWorkers(w *goWorker) bool {
	if p.IsClosed() || p.Cap() > p.Running() {
		return false
	}
	w.recycleTime = time.Now()
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.IsClosed() {
		return false
	}
	if err := p.works.insert(w); err != nil {
		return false
	}
	p.cond.Signal()
	return true
}
