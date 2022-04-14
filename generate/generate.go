package generate

import (
	"sync"
	"time"
)

func init() {
	mTokenBacket = &tokenBacket{
		mutex:          sync.Mutex{},
		storedToken:    0,
		nextFreeTime:   time.Now().UnixMicro(),
		intervalMicro:  int64(time.Second/time.Microsecond) / TokenSpeedPreSecond,
		maxStoredToken: MaxStoredTokenSecond * TokenSpeedPreSecond,
	}
}
