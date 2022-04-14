package generate

import (
	"sync"
	"time"
)

var mTokenBacket SmoothBursty
var TokenSpeedPreSecond int64 = 1
var MaxStoredTokenSecond int64 = 2

func Acquire() int64 {
	return mTokenBacket.acquire(1)
}
func AcquireN(count int64) int64 {
	if count <= 0 {
		return 0
	}
	return mTokenBacket.acquire(count)
}

func Try() (bool, int64) {
	return mTokenBacket.Try()
}

func TryN(count int64) (bool, int64) {
	if count <= 0 {
		return true, 0
	}
	return mTokenBacket.TryN(count)
}

type SmoothBursty interface {
	Acquire() int64
	AcquireN(int64) int64
	acquire(count int64) int64
	reserveEarliestAvailable(count, now int64) int64
	Try() (bool, int64)
	TryN(int64) (bool, int64)
	reSync(now int64)
	ReSet(countPreSecond int64)
}

type tokenBacket struct {
	maxStoredToken int64
	storedToken    int64
	intervalMicro  int64
	nextFreeTime   int64
	mutex          sync.Mutex
}

func (t *tokenBacket) reSync(now int64) {
	//logger.MLogger.Info("duration: ", zap.Int64("before reSync t.storedToken", t.storedToken),
	//		zap.Int64("before reSync now", now),
	//	zap.Int64("before reSync t.nextFreeTime", t.nextFreeTime))
	if now > t.nextFreeTime {
		producedTokenCount := (now - t.nextFreeTime) / t.intervalMicro
		t.storedToken = min(t.maxStoredToken, t.storedToken+producedTokenCount)
		t.nextFreeTime = now
		//logger.MLogger.Info("duration: ", zap.Int64("reSync t.storedToken", t.storedToken),
		//	zap.Int64("reSync producedTokenCount", producedTokenCount))
	}
}

func (t *tokenBacket) Acquire() int64 {
	return t.acquire(1)
}
func (t *tokenBacket) AcquireN(count int64) int64 {
	return t.acquire(count)
}
func (t *tokenBacket) acquire(count int64) int64 {
	//logger.MLogger.Info("-----------------------------------------------------", zap.String("---------------------------------------------------------", ""))
	now := time.Now().UnixMicro()
	nextFreeTime := t.reserveEarliestAvailable(count, now)
	duration := nextFreeTime - now
	time.Sleep(time.Duration(duration) * time.Microsecond)
	//	logger.MLogger.Info("duration: ", zap.Int64("duration", duration))
	return duration
}
func (t *tokenBacket) reserveEarliestAvailable(count, now int64) int64 {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	//更新
	t.reSync(now)
	//返回值
	val := t.nextFreeTime
	//可以拿到多少
	accessToken := min(count, t.storedToken)
	//剩余多少
	//	logger.MLogger.Info("duration: ", zap.Int64("reserveEarliestAvailable t.storedToken", t.storedToken))
	t.storedToken = t.storedToken - accessToken
	//	logger.MLogger.Info("duration: ", zap.Int64("reserveEarliestAvailable t.storedToken", t.storedToken))
	//多少不够 透支的
	insufficientToken := count - accessToken
	//logger.MLogger.Info("duration: ", zap.Int64("insufficientToken", insufficientToken))
	//下一次更新的时间
	t.nextFreeTime = now + t.intervalMicro*insufficientToken
	return val
}
func (t *tokenBacket) Try() (bool, int64) {
	now := time.Now().UnixMicro()
	if t.storedToken == 0 && now < t.nextFreeTime+t.intervalMicro {
		return false, 0
	}
	return true, t.acquire(1)
}

func (t *tokenBacket) TryN(count int64) (bool, int64) {
	if count <= 0 {
		return true, 0
	}
	if t.storedToken < count && now < (count-t.maxStoredToken)*t.intervalMicro+t.intervalMicro {
		return false, 0
	}
	return true, t.acquire(count)
}

func (t *tokenBacket) ReSet(countPreSecond int64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	now := time.Now().UnixMicro()
	t.reSync(now)

	//扩容
	oldInterval := t.intervalMicro
	t.intervalMicro = int64(time.Microsecond) * countPreSecond
	//等比例增长
	proportion := t.intervalMicro / oldInterval
	t.maxStoredToken *= proportion
	t.storedToken *= proportion
}

func min(a, b int64) int64 {
	if a >= b {
		return b
	}
	return a
}
