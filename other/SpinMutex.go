package other

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type SpinMutex uint32

const maxSpinTime = 16

func (s *SpinMutex) Lock() {
	backOff := 1
	for !atomic.CompareAndSwapUint32((*uint32)(s), 0, 1) {
		for i := 0; i <= backOff; i++ {
			runtime.Gosched()
		}
		if backOff < maxSpinTime {
			backOff <<= 1
		}
	}
}

func (s *SpinMutex) Unlock() {
	atomic.StoreUint32((*uint32)(s), 0)
}

func NewSpinLock() sync.Locker {
	return new(SpinMutex)
}
